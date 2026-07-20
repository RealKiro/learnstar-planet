<?php

declare(strict_types=1);

namespace App\Http\Requests\Admin;

use Illuminate\Foundation\Http\FormRequest;

class CreateClassRoomRequest extends FormRequest
{
    public function authorize(): bool
    {
        return $this->user()?->role === 'school_admin';
    }

    public function rules(): array
    {
        return [
            'name' => ['required', 'string', 'max:100'],
            'grade' => ['nullable', 'string', 'max:50'],
            'year' => ['nullable', 'string', 'max:20'],
        ];
    }
}

