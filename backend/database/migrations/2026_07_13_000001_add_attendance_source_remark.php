<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class () extends Migration {
    public function up(): void
    {
        Schema::table('attendances', function (Blueprint $t) {
            if (!Schema::hasColumn('attendances', 'source')) {
                $t->string('source', 20)->default('auto')->after('remark');
            }
            if (!Schema::hasColumn('attendances', 'remark')) {
                $t->string('remark', 500)->nullable()->after('status');
            }
            if (!Schema::hasColumn('attendances', 'leave_record_id')) {
                $t->unsignedBigInteger('leave_record_id')->nullable()->after('remark');
            }
        });
    }

    public function down(): void
    {
        Schema::table('attendances', function (Blueprint $t) {
            foreach (['source', 'remark', 'leave_record_id'] as $c) {
                if (Schema::hasColumn('attendances', $c)) {
                    $t->dropColumn($c);
                }
            }
        });
    }
};
