<?php

declare(strict_types=1);

namespace App\Services;

use App\Events\NoticePublished;
use App\Models\Notice;
use App\Models\User;
use Illuminate\Contracts\Pagination\LengthAwarePaginator;

/**
 * 班级通知（公告）读写。
 *
 * 只负责数据操作与事件广播；HTTP 校验与响应包装留在控制器。
 * 所有查询都经 TeacherClassScope 限定班级作用域。
 */
class NoticeService
{
    public function __construct(
        private readonly TeacherClassScope $scope,
    ) {
    }

    /** 通知列表（按教师可管理班级过滤，每页 20 条） */
    public function paginate(User $teacher): LengthAwarePaginator
    {
        return Notice::whereIn('class_id', $this->scope->ids($teacher))
            ->orderBy('created_at', 'desc')
            ->paginate(20);
    }

    /**
     * 新建通知，落在教师的第一个班级，默认未发布。
     *
     * @param  array{title: string, content: string, type: string|null}  $attributes
     */
    public function create(User $teacher, array $attributes): Notice
    {
        return Notice::create([
            'class_id' => $this->scope->ids($teacher)->first(),
            'title' => $attributes['title'],
            'content' => $attributes['content'],
            'type' => $attributes['type'],
            'published_by' => $teacher->id,
            'is_published' => false,
        ]);
    }

    /** 按作用域查找通知（越权即 404） */
    public function findInScope(User $teacher, int $id): Notice
    {
        return Notice::whereIn('class_id', $this->scope->ids($teacher))->findOrFail($id);
    }

    /**
     * @param  array<string, mixed>  $attributes
     */
    public function update(Notice $notice, array $attributes): Notice
    {
        $notice->update($attributes);

        return $notice;
    }

    /** 发布通知并广播 NoticePublished 事件到班级大屏 */
    public function publish(Notice $notice): Notice
    {
        $notice->update(['is_published' => true]);

        event(new NoticePublished(
            $notice->class_id,
            $notice->id,
            $notice->title,
            $notice->type,
        ));

        return $notice;
    }

    /** 撤回（取消发布）通知 */
    public function unpublish(Notice $notice): Notice
    {
        $notice->update(['is_published' => false]);

        return $notice;
    }

    public function delete(Notice $notice): void
    {
        $notice->delete();
    }
}
