<?php
declare(strict_types=1);
namespace App\Http\Controllers\Api;
use App\Http\Controllers\Controller;
use App\Services\WechatWorkAttendanceService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
class WechatWorkWebhookController extends Controller
{
    public function __construct(private readonly WechatWorkAttendanceService $s) {}
    public function verify(Request $request): string {
        $ms=$request->query('msg_signature','');$ts=$request->query('timestamp','');$nn=$request->query('nonce','');$es=$request->query('echostr','');
        if(empty($ms)||empty($ts)||empty($nn)||empty($es))return '';
        try{return $this->decryptEcho($ms,$ts,$nn,$es);}catch(\Throwable $e){Log::error('企微验证失败',['e'=>$e->getMessage()]);return'';}
    }
    private function decryptEcho(string $ms,string $ts,string $nn,string $es):string{
        $tk=config('wechat-work.token','');$ak=config('wechat-work.encoding_aes_key','');
        if(empty($tk)||empty($ak))throw new \RuntimeException('未配置');
        if($this->sig($tk,$ts,$nn,$es)!==$ms)throw new \RuntimeException('签名错误');
        $akey=base64_decode($ak.'=',true);if($akey===false)throw new \RuntimeException('key decode fail');
        $dec=openssl_decrypt(base64_decode($es)?:'','AES-256-CBC',$akey,OPENSSL_RAW_DATA|OPENSSL_ZERO_PADDING,substr($akey,0,16));
        $dec=$this->strip($dec);$c=substr($dec,16);
        if(preg_match('/<echostr>(.*?)<\/echostr>/',$c,$m))return $m[1];
        $l=unpack('N',substr($c,0,4));return substr($c,4,$l[1]??0);
    }
    public function receive(Request $request): JsonResponse {
        $ms=$request->query('msg_signature','');$ts=$request->query('timestamp','');$nn=$request->query('nonce','');
        $body=$request->getContent();if($body===null||$body==='')return response()->json(['errcode'=>0,'errmsg'=>'ok']);
        try{
            $xml=simplexml_load_string($body,\SimpleXMLElement::class,LIBXML_NOCDATA);
            if($xml===false)return response()->json(['errcode'=>0,'errmsg'=>'ok']);
            $enc=(string)($xml->Encrypt??'');if($enc==='')return response()->json(['errcode'=>0,'errmsg'=>'ok']);
            $dec=$this->decryptMsg($ms,$ts,$nn,$enc);if($dec===null)return response()->json(['errcode'=>0,'errmsg'=>'ok']);
            $ex=simplexml_load_string($dec,\SimpleXMLElement::class,LIBXML_NOCDATA);if($ex===false)return response()->json(['errcode'=>0,'errmsg'=>'ok']);
            $ev=(string)($ex->Event??'');$ct=(string)($ex->ChangeType??'');
            if($ev==='sys_approval_change'||$ct==='sys_approval_change'){$ai=$ex->ApprovalInfo??null;$sn=(string)($ai?->SpNo??'');$ss=(int)($ai?->SpStatus??1);if($sn!==''){$sid=(int)$request->query('school_id','0');if($sid>0)$this->s->handleWebhookCallback($sid,['sp_no'=>$sn,'sp_status'=>$ss]);}}
        }catch(\Throwable $e){Log::error('企微回调异常',['e'=>$e->getMessage()]);}
        return response()->json(['errcode'=>0,'errmsg'=>'ok']);
    }
    private function decryptMsg(string $ms,string $ts,string $nn,string $enc):?string{
        $tk=config('wechat-work.token','');$ak=config('wechat-work.encoding_aes_key','');
        if(empty($tk)||empty($ak))return null;
        if($this->sig($tk,$ts,$nn,$enc)!==$ms)return null;
        $akey=base64_decode($ak.'=',true);if($akey===false)return null;
        $dec=openssl_decrypt(base64_decode($enc)?:'','AES-256-CBC',$akey,OPENSSL_RAW_DATA|OPENSSL_ZERO_PADDING,substr($akey,0,16));
        $dec=$this->strip($dec);$c=substr($dec,16);$l=unpack('N',substr($c,0,4));return substr($c,4,$l[1]??0);
    }
    private function sig(string $tk,string $ts,string $nn,string $enc):string{$p=[$tk,$ts,$nn,$enc];sort($p,SORT_STRING);return sha1(implode('',$p));}
    private function strip(string $d):string{$l=strlen($d);if($l===0)return'';$p=ord($d[$l-1]);if($p<1||$p>32)return$d;return substr($d,0,$l-$p);}
}
