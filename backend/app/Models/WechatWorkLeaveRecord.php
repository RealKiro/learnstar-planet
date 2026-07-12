<?php
declare(strict_types=1);
namespace App\Models;
use Illuminate\Database\Eloquent\Model;
class WechatWorkLeaveRecord extends Model
{
    protected $table = 'wechat_work_leave_records';
    protected $fillable = ['school_id','class_id','student_id','parent_wework_userid','student_name_from_wework','sp_no','leave_start_date','leave_end_date','leave_type','reason','approve_status','approved_at','raw_data','synced_at'];
    protected $casts = ['leave_start_date'=>'date','leave_end_date'=>'date','approved_at'=>'datetime','raw_data'=>'array','synced_at'=>'datetime'];
    public static function findForStudentOnDate(int $studentId, string $date): ?self { return static::where('student_id',$studentId)->where('approve_status','approved')->whereDate('leave_start_date','<=',$date)->whereDate('leave_end_date','>=',$date)->latest('approved_at')->first(); }
    /** @return \Illuminate\Database\Eloquent\Collection<int, static> */
    public static function getUnsyncedForDate(string $date) { return static::where('approve_status','approved')->whereNull('synced_at')->whereDate('leave_start_date','<=',$date)->whereDate('leave_end_date','>=',$date)->get(); }
}
