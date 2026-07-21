<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class () extends Migration {
    public function up(): void
    {
        Schema::table('shop_items', function (Blueprint $table) {
            if (!Schema::hasColumn('shop_items', 'school_id')) {
                $table->foreignId('school_id')->nullable()->constrained()->cascadeOnDelete()->after('class_id');
            }
            if (!Schema::hasColumn('shop_items', 'currency_type')) {
                $table->string('currency_type', 20)->default('score')->after('cost_score');
            }
        });
    }

    public function down(): void
    {
        Schema::table('shop_items', function (Blueprint $table) {
            $table->dropColumn(['school_id', 'currency_type']);
        });
    }
};
