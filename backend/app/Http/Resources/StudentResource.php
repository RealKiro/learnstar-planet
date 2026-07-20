<?php

declare(strict_types=1);

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class StudentResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'gender' => $this->gender,
            'student_no' => $this->student_no,
            'class_id' => $this->class_id,
            'parent_id' => $this->parent_id,
            'total_score' => (int) $this->total_score,
            'status' => $this->status,
            'class_room' => $this->whenLoaded('classRoom', fn () => [
                'id' => $this->classRoom->id,
                'name' => $this->classRoom->name,
                'grade' => $this->classRoom->grade,
            ]),
            'created_at' => $this->created_at?->toDateTimeString(),
        ];
    }
}

