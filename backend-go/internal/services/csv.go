// CSV 输出辅助：Laravel 端各类导出走 maatwebsite/excel（.xlsx），
// Go 端遵循「零新依赖」约定改用标准库 encoding/csv。
//
// 统一口径：
//   - 行尾 "\n"（csv.Writer 默认，非 CRLF）；
//   - 开头写入 UTF-8 BOM，让 Excel 正确识别中文（同 Laravel 批量导入模板的做法）；
//   - 单元格内的换行按 CSV 规则加引号（Excel 显示为单元格内换行）。
package services

import (
	"bytes"
	"encoding/csv"
)

// csvBOM UTF-8 BOM。
const csvBOM = "\ufeff"

// csvText 用 encoding/csv 生成带 BOM 的 CSV 字节流。
// 空行（len(row) == 0）输出为仅含换行的空记录，与 xlsx 网格里的「空行」对应。
func csvText(rows [][]string) []byte {
	var buf bytes.Buffer
	buf.WriteString(csvBOM)

	writer := csv.NewWriter(&buf)
	for _, row := range rows {
		// bytes.Buffer 的写入不会失败；csv.Writer 内部错误（无）忽略。
		_ = writer.Write(row)
	}
	writer.Flush()

	return buf.Bytes()
}
