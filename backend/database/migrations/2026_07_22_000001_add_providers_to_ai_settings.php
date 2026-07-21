<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class () extends Migration {
    public function up(): void
    {
        Schema::table('ai_settings', function (Blueprint $table) {
            if (!Schema::hasColumn('ai_settings', 'providers')) {
                $table->json('providers')->nullable()->after('tokens_limit');
            }
        });
    }

    public function down(): void
    {
        Schema::table('ai_settings', function (Blueprint $table) {
            $table->dropColumn('providers');
        });
    }
};
