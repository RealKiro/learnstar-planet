<?php

declare(strict_types=1);

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class UserResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'username' => $this->username,
            'name' => $this->name,
            'nickname' => $this->nickname,
            'role' => $this->role,
            'school_id' => $this->school_id,
            'avatar_path' => $this->avatar_path,
            'phone' => $this->phone,
            'email' => $this->email,
            'status' => $this->status,
            'bindings' => $this->whenLoaded('bindings', fn () => $this->bindings->pluck('platform')),
            'class_names' => $this->whenLoaded('classRooms', fn () => $this->classRooms->pluck('name')),
        ];
    }
}
