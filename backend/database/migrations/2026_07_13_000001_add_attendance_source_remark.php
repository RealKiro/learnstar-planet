<?php declare(strict_types=1); use Illuminate\Database\Migrations\Migration; use Illuminate\Database\Schema\Blueprint; use Illuminate\Support\Facades\Schema; return new class extends Migration {
    public function up(): void { Schema::table('attendances', function (Blueprint $t) { if(!Schema::hasColumn('attendances','source')) $t->string('source',20)->default('auto')->after('remark'); if(!Schema::hasColumn('attendances','remark')) $t->string('remark',500)->nullable()->after('status'); if(!Schema::hasColumn('attendances','leave_record_id')) $t->unsignedBigInteger('leave_record_id')->nullable()->after('remark'); }); }
    public function down(): void { Schema::table('attendances', function (Blueprint $t) { foreach(['source','remark','leave_record_id'] as $c) if(Schema::hasColumn('attendances',$c)) $t->dropColumn($c); }); }
}; EOF

cat > "$BP/backend/database/migrations/2026_07_13_000002_create_wechat_work_leave_records.php" << 'EOF'
<?php declare(strict_types=1); use Illuminate\Database\Migrations\Migration; use Illuminate\Database\Schema\Blueprint; use Illuminate\Support\Facades\Schema; return new class extends Migration {
    public function up(): void { Schema::create('wechat_work_leave_records', function (Blueprint $t) { $t->id(); $t->unsignedBigInteger('school_id'); $t->unsignedBigInteger('class_id')->nullable(); $t->unsignedBigInteger('student_id')->nullable(); $t->string('parent_wework_userid',128); $t->string('student_name_from_wework',100)->nullable(); $t->string('sp_no',128); $t->date('leave_start_date'); $t->date('leave_end_date'); $t->string('leave_type',50)->nullable(); $t->text('reason')->nullable(); $t->string('approve_status',20)->default('approved'); $t->dateTime('approved_at')->nullable(); $t->json('raw_data')->nullable(); $t->dateTime('synced_at')->nullable(); $t->timestamps(); $t->unique('sp_no'); $t->index(['school_id','leave_start_date']); $t->index(['student_id','leave_start_date']); }); }
    public function down(): void { Schema::dropIfExists('wechat_work_leave_records'); }
}; EOF

cat > "$BP/backend/database/migrations/2026_07_13_000003_create_class_room_teachers.php" << 'EOF'
<?php declare(strict_types=1); use Illuminate\Database\Migrations\Migration; use Illuminate\Database\Schema\Blueprint; use Illuminate\Support\Facades\Schema; return new class extends Migration {
    public function up(): void { Schema::create('class_room_teachers', function (Blueprint $t) { $t->id(); $t->foreignId('class_room_id')->constrained('class_rooms')->cascadeOnDelete(); $t->foreignId('user_id')->constrained('users')->cascadeOnDelete(); $t->string('role',20)->default('co_teacher'); $t->timestamps(); $t->unique(['class_room_id','user_id']); }); }
    public function down(): void { Schema::dropIfExists('class_room_teachers'); }
}; EOF

# Config
cat > "$BP/backend/config/wechat-work.php" << 'EOF'
<?php return ['corp_id'=>env('WECOM_CORP_ID',''),'agent_id'=>(int)env('WECOM_AGENT_ID',0),'secret'=>env('WECOM_SECRET',''),'token'=>env('WECOM_TOKEN',''),'encoding_aes_key'=>env('WECOM_ENCODING_AES_KEY',''),'leave_template_id'=>env('WECOM_LEAVE_TEMPLATE_ID',''),];
