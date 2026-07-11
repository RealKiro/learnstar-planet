<?php

declare(strict_types=1);

namespace App\Http\Requests\Admin;

use Illuminate\Foundation\Http\FormRequest;

class BatchCreateTeachersRequest extends FormRequest
{
    public function authorize(): bool
    {
        return $this->user()?->role === 'school_admin';
    }

    public function rules(): array
    {
        return [
            'teachers' => ['required', 'array', 'min:1', 'max:200'],
            'teachers.*.name' => ['required', 'string', 'max:50'],
            'teachers.*.password' => ['required', 'string', 'min:6', 'max:64'],
            'teachers.*.phone' => ['nullable', 'string', 'max:20'],
            'teachers.*.email' => ['nullable', 'email', 'max:100'],
        ];
    }

    public function messages(): array
    {
        return [
            'teachers.*.name.required' => '教师姓名不能为空',
            'teachers.*.password.required' => '教师密码不能为空',
            'teachers.*.password.min' => '密码至少6位',
            'teachers.*.email.email' => '邮箱格式不正确',
        ];
    }
}
