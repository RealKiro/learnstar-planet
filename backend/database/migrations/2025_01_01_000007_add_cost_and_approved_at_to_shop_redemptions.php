<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::table('shop_redemptions', function (Blueprint $table) {
            $table->integer('cost')->default(0)->after('class_id');
            $table->timestamp('approved_at')->nullable()->after('approved_by');
        });
    }

    public function down(): void
    {
        Schema::table('shop_redemptions', function (Blueprint $table) {
            $table->dropColumn(['cost', 'approved_at']);
        });
    }
};
