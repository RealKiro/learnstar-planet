<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class () extends Migration {
    public function up(): void
    {
        Schema::table('users', function (Blueprint $table) {
            if (!Schema::hasColumn('users', 'subject')) {
                $table->string('subject', 50)->nullable()->after('nickname')->comment('科目');
            }
            if (!Schema::hasColumn('users', 'grade_team')) {
                $table->string('grade_team', 50)->nullable()->after('subject')->comment('所属年级团队');
            }
        });
    }

    public function down(): void
    {
        Schema::table('users', function (Blueprint $table) {
            $table->dropColumn(['subject', 'grade_team']);
        });
    }
};
