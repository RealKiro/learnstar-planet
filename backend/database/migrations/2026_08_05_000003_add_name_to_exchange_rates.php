<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class () extends Migration {
    public function up(): void
    {
        Schema::table('exchange_rates', function (Blueprint $table) {
            if (!Schema::hasColumn('exchange_rates', 'name')) {
                $table->string('name', 100)->nullable()->after('school_id');
            }
        });
    }

    public function down(): void
    {
        Schema::table('exchange_rates', function (Blueprint $table) {
            if (Schema::hasColumn('exchange_rates', 'name')) {
                $table->dropColumn('name');
            }
        });
    }
};
