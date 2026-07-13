<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

return new class () extends Migration {
    public function up(): void
    {
        DB::table('class_rooms')->orderBy('id')->each(function (object $row): void {
            $name = \App\Models\ClassRoom::normalizeClassName($row->name);
            if ($name !== $row->name) {
                DB::table('class_rooms')->where('id', $row->id)->update(['name' => $name]);
            }
        });
    }

    public function down(): void
    {
        // no rollback — normalization is one-way
    }
};
