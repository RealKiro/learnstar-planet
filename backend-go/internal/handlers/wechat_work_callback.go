// 企业微信回调（公开路由）：GET 验签回显 echostr + POST 接收审批事件。
//
// 逐条移植自 Laravel App\Http\Controllers\Api\WechatWorkWebhookController：
//
//	GET  verify  —— query `msg_signature/timestamp/nonce/echostr` 任一为空、配置缺失、签名错、
//	                解密失败一律返回**空字符串**（Laravel 记 error 日志）；成功则把解密出的
//	                echostr 明文**直接作为响应体**（不是 JSON）。
//	POST receive —— 空 body / 非 XML / 无 `<Encrypt>` / 解密失败 / 解密后非 XML 一律返回
//	                `{"errcode":0,"errmsg":"ok"}`（HTTP 200）；`Event` 或 `ChangeType` 命中
//	                `sys_approval_change`、`<ApprovalInfo><SpNo>` 非空、query `school_id > 0` 时
//	                才调用请假同步服务；任何异常都吞掉（记日志）并仍返回 ok。
//
// 加密件只用标准库（crypto/aes、crypto/cipher、crypto/sha1、encoding/base64、encoding/xml，零新依赖）：
// key = base64(encoding_aes_key + "=")，IV = key[:16]，AES-256-CBC，**无填充校验**
// （对应 PHP `OPENSSL_RAW_DATA | OPENSSL_ZERO_PADDING`），解密后先按末字节 1..32 去掉 PKCS7 填充
// （Laravel `strip()`），再跳过前 16 字节随机串，最后优先匹配 `<echostr>…</echostr>`，
// 否则按 4 字节**大端**长度截取正文。
package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// WechatWorkCallbackVerify GET /api/v1/wechat-work/callback（公开）。
func (h *Handlers) WechatWorkCallbackVerify(c *gin.Context) {
	msgSignature := c.Query("msg_signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echoStr := c.Query("echostr")
	if msgSignature == "" || timestamp == "" || nonce == "" || echoStr == "" {
		c.String(http.StatusOK, "")
		return
	}

	result, err := h.decryptEcho(msgSignature, timestamp, nonce, echoStr)
	if err != nil {
		log.Printf("企微验证失败：%v", err)
		c.String(http.StatusOK, "")
		return
	}
	c.String(http.StatusOK, result)
}

// WechatWorkCallbackReceive POST /api/v1/wechat-work/callback（公开）。
func (h *Handlers) WechatWorkCallbackReceive(c *gin.Context) {
	body, readErr := io.ReadAll(c.Request.Body)
	if readErr != nil || len(body) == 0 {
		wechatWorkCallbackOK(c)
		return
	}

	envelope, ok := parseWechatWorkXML(body)
	if !ok || envelope.Encrypt == "" {
		wechatWorkCallbackOK(c)
		return
	}

	plain, ok := h.decryptMsg(
		c.Query("msg_signature"), c.Query("timestamp"), c.Query("nonce"), envelope.Encrypt,
	)
	if !ok {
		wechatWorkCallbackOK(c)
		return
	}

	message, ok := parseWechatWorkXML([]byte(plain))
	if !ok {
		wechatWorkCallbackOK(c)
		return
	}

	if isWechatWorkApprovalChange(message) && message.ApprovalInfo != nil {
		if spNo := message.ApprovalInfo.SpNo; spNo != "" {
			schoolID := services.PhpIntCast(c.Query("school_id"))
			if schoolID > 0 {
				spStatus := 1
				if message.ApprovalInfo.SpStatus != nil {
					spStatus = services.PhpIntCast(*message.ApprovalInfo.SpStatus)
				}
				if err := h.wechatAttendance.HandleWebhookCallback(uint(schoolID), spNo, spStatus); err != nil {
					log.Printf("企微回调异常：%v", err)
				}
			}
		}
	}

	wechatWorkCallbackOK(c)
}

// wechatWorkCallbackOK 统一的回调响应（同 Laravel `response()->json(['errcode' => 0,'errmsg' => 'ok'])`）。
func wechatWorkCallbackOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"errcode": 0, "errmsg": "ok"})
}

// wechatWorkCallbackEnvelope 回调/事件报文（`<xml>…</xml>`）的最小结构。
//
// 不带 XMLName：PHP 的 `$xml->Encrypt` 取的是**根元素的直接子元素**，无根名的结构体语义一致
// （根元素名不同也能取到；嵌套在别的元素里的同名元素两边都取不到）。
// `SpStatus` 用指针以区分「元素缺失」（PHP `?? 1` → 1）与「元素为空」（`(int) ”` → 0）。
type wechatWorkCallbackEnvelope struct {
	Encrypt      string                  `xml:"Encrypt"`
	Event        string                  `xml:"Event"`
	ChangeType   string                  `xml:"ChangeType"`
	ApprovalInfo *wechatWorkApprovalInfo `xml:"ApprovalInfo"`
}

// wechatWorkApprovalInfo 审批事件载荷（`<ApprovalInfo>`）。
type wechatWorkApprovalInfo struct {
	SpNo       string  `xml:"SpNo"`
	SpStatus   *string `xml:"SpStatus"`
	ChangeType string  `xml:"ChangeType"`
}

// parseWechatWorkXML 用 encoding/xml 解析报文；非 XML 返回 ok=false（对应 simplexml_load_string 失败）。
// `LIBXML_NOCDATA` 无需额外处理：encoding/xml 会把 `<![CDATA[…]]>` 直接并入文本。
func parseWechatWorkXML(body []byte) (wechatWorkCallbackEnvelope, bool) {
	var envelope wechatWorkCallbackEnvelope
	if err := xml.Unmarshal(body, &envelope); err != nil {
		return wechatWorkCallbackEnvelope{}, false
	}
	return envelope, true
}

// isWechatWorkApprovalChange 判定是否审批状态变更事件：
// 顶层 `Event` 或 `ChangeType`（Laravel 只判这两个）；
// 有意补一条兜底 —— 真实企微报文里 `ChangeType` 也可能出现在 `ApprovalInfo` 内层（Laravel 会漏判）。
func isWechatWorkApprovalChange(message wechatWorkCallbackEnvelope) bool {
	const change = "sys_approval_change"
	if message.Event == change || message.ChangeType == change {
		return true
	}
	return message.ApprovalInfo != nil && message.ApprovalInfo.ChangeType == change
}

// decryptEcho 同 Laravel WechatWorkWebhookController::decryptEcho。
func (h *Handlers) decryptEcho(msgSignature, timestamp, nonce, echoStr string) (string, error) {
	token, aesKey := h.wechatWorkConfig()
	if token == "" || aesKey == "" {
		return "", errors.New("未配置")
	}
	if wechatWorkSignature(token, timestamp, nonce, echoStr) != msgSignature {
		return "", errors.New("签名错误")
	}
	if wechatWorkAESKey(aesKey) == nil {
		return "", errors.New("key decode fail")
	}
	plain, ok := wechatWorkDecrypt(aesKey, echoStr)
	if !ok {
		// 同 Laravel：openssl_decrypt 失败时 `strip(false)` → ''，最终**无异常**地返回空串。
		return "", nil
	}
	return wechatWorkEcho(plain), nil
}

// decryptMsg 同 Laravel WechatWorkWebhookController::decryptMsg（失败返回 ok=false）。
func (h *Handlers) decryptMsg(msgSignature, timestamp, nonce, encrypted string) (string, bool) {
	token, aesKey := h.wechatWorkConfig()
	if token == "" || aesKey == "" {
		return "", false
	}
	if wechatWorkSignature(token, timestamp, nonce, encrypted) != msgSignature {
		return "", false
	}
	if wechatWorkAESKey(aesKey) == nil {
		return "", false
	}
	plain, ok := wechatWorkDecrypt(aesKey, encrypted)
	if !ok {
		return "", false
	}
	return wechatWorkMessage(plain), true
}

// wechatWorkConfig 读取企微回调凭据（服务构造时从环境变量载入，键名同 Laravel config/wechat-work.php）。
func (h *Handlers) wechatWorkConfig() (string, string) {
	return h.wechatWork.Token, h.wechatWork.EncodingAESKey
}

// wechatWorkSignature 同 Laravel `sha1(implode(”, sort([token, timestamp, nonce, encrypt])))`：
// PHP `sort(..., SORT_STRING)` 是字节序比较，与 Go 的字符串排序一致。
func wechatWorkSignature(token, timestamp, nonce, encrypt string) string {
	parts := []string{token, timestamp, nonce, encrypt}
	sort.Strings(parts)
	sum := sha1.Sum([]byte(strings.Join(parts, "")))
	return hex.EncodeToString(sum[:])
}

// wechatWorkAESKey 同 Laravel `base64_decode($aesKey . '=', true)`：
// 企微 EncodingAESKey 为 43 字符，补一个 '=' 后是合法 base64 → 32 字节。
// PHP 对长度 %4 == 2/3 的输入会自动补 '=' 解码，对 %4 == 1 判为非法（返回 false）；
// Go 的标准库要求显式填充，故这里按 PHP 口径预处理（长度 %4 == 1 直接判非法）。
func wechatWorkAESKey(aesKey string) []byte {
	padded := aesKey + "="
	if rest := len(padded) % 4; rest != 0 {
		if rest == 1 {
			return nil
		}
		padded += strings.Repeat("=", 4-rest)
	}
	key, err := base64.StdEncoding.DecodeString(padded)
	if err != nil {
		return nil
	}
	return key
}

// wechatWorkDecrypt AES-256-CBC 解密并去掉填充（见文件头说明）；ok=false 表示密文/密钥不可解。
func wechatWorkDecrypt(aesKey, ciphertextBase64 string) (string, bool) {
	key := wechatWorkAESKey(aesKey)
	if len(key) != aes.BlockSize*2 {
		return "", false
	}

	// 有意差异：PHP 用**非严格** base64_decode（忽略非法字符），Go 标准库严格报错 —— 只在报文被
	// 篡改/截断时有差别，两边最终都判定为「解密失败」。
	raw, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil || len(raw) == 0 || len(raw)%aes.BlockSize != 0 {
		return "", false
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", false
	}
	plain := make([]byte, len(raw))
	cipher.NewCBCDecrypter(block, key[:aes.BlockSize]).CryptBlocks(plain, raw)
	return wechatWorkStrip(plain), true
}

// wechatWorkStrip 同 Laravel `strip()`：末字节为 1..32 时按 PKCS7 去掉该长度（企微按 32 字节块填充），
// 否则原样返回（例如零填充场景下末字节为 0）。
func wechatWorkStrip(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	pad := int(data[len(data)-1])
	if pad < 1 || pad > 32 {
		return string(data)
	}
	if pad >= len(data) {
		return ""
	}
	return string(data[:len(data)-pad])
}

// wechatWorkEcho 取 echostr：先按 `<echostr>(.*?)</echostr>` 匹配（Go 正则的 `.` 同样不匹配换行，
// 与 Laravel 的 `/…/` 无 `s` 修饰符一致），未命中再走 4 字节大端长度截取。
func wechatWorkEcho(plain string) string {
	if len(plain) < 16 {
		return ""
	}
	body := plain[16:]
	if match := wechatWorkEchoStrPattern.FindStringSubmatch(body); match != nil {
		return match[1]
	}
	return wechatWorkLengthPrefixed(body)
}

// wechatWorkEchoStrPattern 同 Laravel 解密后正文的 `<echostr>` 提取正则。
var wechatWorkEchoStrPattern = regexp.MustCompile(`<echostr>(.*?)</echostr>`)

// wechatWorkMessage 同 Laravel `decryptMsg` 的返回：跳过前 16 字节随机串 → 4 字节大端长度 → 截取。
func wechatWorkMessage(plain string) string {
	if len(plain) < 16 {
		return ""
	}
	return wechatWorkLengthPrefixed(plain[16:])
}

// wechatWorkLengthPrefixed 读 4 字节大端长度并截取正文（长度越界时按 PHP substr 的截断语义处理）。
func wechatWorkLengthPrefixed(body string) string {
	if len(body) < 4 {
		return ""
	}
	length := int(binary.BigEndian.Uint32([]byte(body[:4])))
	body = body[4:]
	if length <= 0 {
		return ""
	}
	if length > len(body) {
		return body
	}
	return body[:length]
}
