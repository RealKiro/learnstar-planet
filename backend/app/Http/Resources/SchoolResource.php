<?php

declare(strict_types=1);

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class SchoolResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'code' => $this->code,
            'address' => $this->address,
            'logo' => $this->logo,
            'contact_phone' => $this->contact_phone,
            'contact_email' => $this->contact_email,
            'status' => $this->status,
            'settings' => $this->settings,
        ];
    }
}
