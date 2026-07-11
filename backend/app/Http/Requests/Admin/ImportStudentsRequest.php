<?php

declare(strict_types=1);

namespace App\Http\Requests\Admin;

use Illuminate\Foundation\Http\FormRequest;

class ImportStudentsRequest extends FormRequest
{
    public function authorize(): bool
    {
        return $this->user()?->role === 'school_admin';
    }

    public function rules(): array
    {
        return [
            'students' => ['required', 'array', 'min:1', 'max:500'],
            'students.*.name' => ['required', 'string', 'max:50'],
            'students.*.class_name' => ['required', 'string', 'max:100'],
            'students.*.gender' => ['required', 'string', 'in:男,女,未知'],
            'students.*.student_no' => ['nullable', 'string', 'max:50'],
            'students.*.phone' => ['nullable', 'string', 'max:20'],
        ];
    }

    public function messages(): array
    {
        return [
            'students.*.name.required' => '学生姓名不能为空',
            'students.*.class_name.required' => '班级不能为空',
            'students.*.gender.required' => '性别不能为空',
        ];
    }
}
