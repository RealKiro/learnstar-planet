<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * 宠物体系统一改造：
     * 1. pets 表新增 species 列（物种 ID，对齐前端 petData.ts 的 10 系列 120 物种体系）
     * 2. 旧 type 数据映射回填到 species（旧体系为 type 字典，含版权命名）
     * 3. 新建 pet_collections 收藏表（图鉴收集 + 切换进度保留，为玩法阶段预留）
     */
    public function up(): void
    {
        Schema::table('pets', function (Blueprint $table) {
            $table->string('species')->nullable()->after('type')->index();
            // 记录最近一次切换宠物时间（备用字段）
            $table->timestamp('last_switched_at')->nullable()->after('last_fed_at');
        });

        // 旧 type → 新 species 映射（未命中的回退到 'zhulong' 烛龙）
        $map = [
            'pikachu' => 'pikachu',
            'eevee' => 'eevee',
            'charmander' => 'charmander',
            'charizard' => 'charmander',
            'squirtle' => 'squirtle',
            'blastoise' => 'squirtle',
            'bulbasaur' => 'bulbasaur',
            'venusaur' => 'bulbasaur',
            'dragonite' => 'wyvern',
            'mewtwo' => 'riolu',
            'gengar' => 'quantum_beast',
            'snorlax' => 'panda',
            'panda' => 'panda',
            'golden_monkey' => 'golden_monkey',
            'crested_ibis' => 'crested_ibis',
            'red_crowned_crane' => 'red_crowned_crane',
            'chinese_alligator' => 'chinese_alligator',
            'snow_leopard' => 'south_china_tiger',
            'siberian_tiger' => 'south_china_tiger',
            'qilin' => 'qilin',
            'pixiu' => 'qilin',
            'fenghuang' => 'fenghuang',
            'zhuque' => 'fenghuang',
            'kunpeng' => 'kunpeng',
            'qinglong' => 'yinglong',
            'zhulong' => 'zhulong',
            'cloud_fox' => 'nine_tail_fox',
            'hope_dragon' => 'zhulong',
            'light_horse' => 'unicorn',
        ];

        // 兼容旧版整班切换写入的带前缀 type（如 'pokemon_pikachu'、'myth_qilin'）
        $seriesPrefixes = ['myth', 'pokemon', 'national', 'mecha', 'magic', 'prehistoric', 'constellation', 'folklore', 'festival', 'qixia', 'cosmic', 'cute', 'treasure', 'mythic'];

        $cases = [];
        $bindings = [];
        foreach ($map as $old => $new) {
            $cases[] = 'WHEN ? THEN ?';
            $bindings[] = $old;
            $bindings[] = $new;
            foreach ($seriesPrefixes as $prefix) {
                $cases[] = 'WHEN ? THEN ?';
                $bindings[] = "{$prefix}_{$old}";
                $bindings[] = $new;
            }
        }
        if ($cases) {
            $sql = 'UPDATE pets SET species = CASE type ' . implode(' ', $cases)
                . ' ELSE \'zhulong\' END WHERE species IS NULL';
            DB::statement($sql, $bindings);
        }

        Schema::create('pet_collections', function (Blueprint $table) {
            $table->id();
            $table->foreignId('student_id')->constrained()->cascadeOnDelete();
            $table->string('species')->index();
            $table->integer('level')->default(1);
            $table->integer('experience')->default(0);
            $table->integer('mood')->default(80);
            $table->boolean('is_active')->default(false);
            $table->timestamps();

            // 同一学生同一物种只保留一条收藏记录
            $table->unique(['student_id', 'species']);
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('pet_collections');
        Schema::table('pets', function (Blueprint $table) {
            $table->dropColumn('species');
            $table->dropColumn('last_switched_at');
        });
    }
};
