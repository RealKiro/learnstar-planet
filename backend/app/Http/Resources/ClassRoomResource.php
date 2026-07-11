<?php

declare(strict_types=1);

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class ClassRoomResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'grade' => $this->grade,
            'year' => $this->year,
            'teacher_id' => $this->teacher_id,
            'teacher_name' => $this->whenLoaded('teacher', fn () => $this->teacher?->name),
            'max_students' => $this->max_students,
            'student_count' => $this->whenCounted('students'),
            'status' => $this->status,
            'created_at' => $this->created_at?->toDateTimeString(),
        ];
    }
}
