<?php

// Stub file for overtrue/pinyin (optional dependency — used only when installed)
// This allows PHPStan to resolve Overtrue\Pinyin\Pinyin without requiring the package.

namespace Overtrue\Pinyin;

class Pinyin
{
    /**
     * Convert Chinese string to pinyin permalink (without tone marks)
     *
     * @param string $string
     * @param string $delimiter
     * @return string
     */
    public static function permalink(string $string, string $delimiter = '-'): string
    {
        return '';
    }
}
