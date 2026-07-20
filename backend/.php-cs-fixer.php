<?php

# PHP CS Fixer configuration
# @see https://github.com/FriendsOfPHP/PHP-CS-Fixer

$finder = (new PhpCsFixer\Finder())
    ->in(__DIR__)
    ->name('*.php')
    ->notName('*.blade.php')
    ->exclude('vendor')
    ->exclude('storage')
    ->exclude('bootstrap/cache')
    ->exclude('tests/phpstan-stubs');

return (new PhpCsFixer\Config())
    ->setRiskyAllowed(true)
    ->setRules([
        // PSR-12 基础规则
        '@PSR12' => true,

        // 数组语法
        'array_syntax' => ['syntax' => 'short'],

        // 类型声明
        'declare_strict_types' => true,
        'visibility_required' => ['elements' => ['method', 'property']],

        // 空白和缩进
        'no_trailing_whitespace' => true,
        'no_trailing_whitespace_in_comment' => true,
        'blank_line_after_namespace' => true,
        'blank_line_after_opening_tag' => true,
        'blank_line_before_statement' => [
            'statements' => ['declare', 'return', 'throw', 'try'],
        ],
        'no_extra_blank_lines' => [
            'tokens' => [
                'curly_brace_block',
                'extra',
                'parenthesis_brace_block',
                'square_brace_block',
                'throw',
                'use',
            ],
        ],

        // 括号和空格
        'no_spaces_inside_parenthesis' => true,
        'space_after_semicolon' => true,

        // 方法和函数
        'method_argument_space' => [
            'on_multiline' => 'ensure_fully_multiline',
            'keep_multiple_spaces_after_comma' => false,
        ],
        'no_empty_statement' => true,

        // 命名约定
        'constant_case' => ['case' => 'lower'],
        'lowercase_keywords' => true,
        'lowercase_static_reference' => true,

        // 导入和命名空间
        'ordered_imports' => [
            'imports_order' => ['class', 'function', 'const'],
            'sort_algorithm' => 'alpha',
        ],
        'single_line_after_imports' => true,
        'no_unused_imports' => true,

        // 控制结构
        'control_structure_braces' => true,
        'control_structure_continuation_position' => [
            'position' => 'same_line',
        ],
        'trailing_comma_in_multiline' => true,

        // 字符串和类型转换
        'single_quote' => true,
        'cast_spaces' => true,
        'concat_space' => [
            'spacing' => 'one',
        ],

        // 类、接口、Trait
        'class_attributes_separation' => [
            'elements' => [
                'method' => 'one',
                'property' => 'one',
            ],
        ],
        'no_blank_lines_after_class_opening' => true,
        'single_class_element_per_statement' => false,

        // 文件末尾
        'single_blank_line_at_eof' => true,
        'no_closing_tag' => true,
    ])
    ->setFinder($finder);
