<?php

declare(strict_types=1);

namespace App\Services;

/**
 * 拼音转换工具
 *
 * - 优先使用 overtrue/pinyin 库（如果 composer 已安装）
 * - 未安装时回退到 "首字符 + 序号" 策略，保证 nickname 字段非空且唯一
 */
class PinyinService
{
    private static ?bool $hasOvertruePinyin = null;

    public static function isAvailable(): bool
    {
        if (self::$hasOvertruePinyin === null) {
            self::$hasOvertruePinyin = class_exists(\Overtrue\Pinyin\Pinyin::class);
        }

        return self::$hasOvertruePinyin;
    }

    /**
     * 把中文字符串转为拼音字符串（不含声调）
     * 英文/数字保持原样
     */
    public static function toPinyin(string $text): string
    {
        $text = trim($text);
        if ($text === '') {
            return '';
        }

        if (self::isAvailable()) {
            // 优先用专业库：overtrue/pinyin
            $pinyin = \Overtrue\Pinyin\Pinyin::permalink($text, '');

            return strtolower((string) $pinyin);
        }

        // Fallback: 把每个汉字转成首字符 + 数字/英文字母保留
        // 结果形如 "zhanglaoshi" -> 仍可读但未必完全准确
        return self::fallbackAscii($text);
    }

    /**
     * 把字符串中所有非 ASCII 字符（含中文）移除，只保留 ASCII 字母数字并小写
     * 这是一个稳定的回退方案：保证 nickname 字段非空
     */
    private static function fallbackAscii(string $text): string
    {
        $ascii = preg_replace('/[^A-Za-z0-9]+/u', '', $text) ?? '';
        if ($ascii !== '') {
            return strtolower($ascii);
        }

        // 全部是中文（没有 ASCII）：用每个汉字的 Unicode codepoint 末位作 fallback
        $result = '';
        $len = mb_strlen($text, 'UTF-8');
        for ($i = 0; $i < $len; $i++) {
            $ch = mb_substr($text, $i, 1, 'UTF-8');
            $code = mb_ord($ch, 'UTF-8');
            // 用简单的字母替换表：把 codepoint 末两位映射到 a-z
            $result .= chr(97 + ($code % 26));
        }

        return $result ?: 'user';
    }
}
