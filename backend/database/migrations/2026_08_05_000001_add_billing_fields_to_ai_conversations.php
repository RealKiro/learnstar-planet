<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class () extends Migration {
    public function up(): void
    {
        Schema::table('ai_conversations', function (Blueprint $table) {
            if (!Schema::hasColumn('ai_conversations', 'prompt_tokens')) {
                $table->integer('prompt_tokens')->default(0)->after('tokens_used');
            }
            if (!Schema::hasColumn('ai_conversations', 'completion_tokens')) {
                $table->integer('completion_tokens')->default(0)->after('prompt_tokens');
            }
            if (!Schema::hasColumn('ai_conversations', 'cost')) {
                $table->decimal('cost', 12, 6)->default(0)->after('completion_tokens');
            }
            if (!Schema::hasColumn('ai_conversations', 'currency')) {
                $table->string('currency', 8)->default('CNY')->after('cost');
            }
        });
    }

    public function down(): void
    {
        Schema::table('ai_conversations', function (Blueprint $table) {
            $columns = ['currency', 'cost', 'completion_tokens', 'prompt_tokens'];
            foreach ($columns as $column) {
                if (Schema::hasColumn('ai_conversations', $column)) {
                    $table->dropColumn($column);
                }
            }
        });
    }
};
