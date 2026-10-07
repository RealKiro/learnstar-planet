package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
)

// 通知在创建时落在教师首个管辖班级，默认未发布。
func TestNoticeCreateDefaultsToFirstClassUnpublished(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	class, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	n, err := services.NewNoticeService(db, services.NewScope(db)).Create(&teacher, "春游通知", "本周五春游", "event")
	if err != nil {
		t.Fatalf("create notice: %v", err)
	}
	if n.ClassID != class.ID {
		t.Fatalf("ClassID = %d, want %d", n.ClassID, class.ID)
	}
	if n.IsPublished {
		t.Fatalf("new notice should be unpublished")
	}
	if n.Type != "event" {
		t.Fatalf("Type = %q, want event", n.Type)
	}
}

// 未传类型时默认为 info。
func TestNoticeCreateDefaultsTypeInfo(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	n, err := services.NewNoticeService(db, services.NewScope(db)).Create(&teacher, "标题", "内容", "")
	if err != nil {
		t.Fatalf("create notice: %v", err)
	}
	if n.Type != "info" {
		t.Fatalf("Type = %q, want info", n.Type)
	}
}

// 教师在无可管辖班级时，创建通知应报 422。
func TestNoticeCreateNoClass(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")

	_, err := services.NewNoticeService(db, services.NewScope(db)).Create(&teacher, "标题", "内容", "info")
	if err == nil {
		t.Fatalf("expected error for teacher with no class")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 422 {
		t.Fatalf("expected 422, got %v", err)
	}
}

// 教师只能看到本班通知，越权查找返回 404。
func TestNoticeScopeGuard(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	t1 := seedTeacher(t, db, school.ID, "teacher1")
	t2 := seedTeacher(t, db, school.ID, "teacher2")
	_, _ = seedTeacherClass(t, db, school.ID, t1.ID, "一班")
	_, _ = seedTeacherClass(t, db, school.ID, t2.ID, "二班")

	svc := services.NewNoticeService(db, services.NewScope(db))
	n, err := svc.Create(&t1, "保密", "仅一班", "info")
	if err != nil {
		t.Fatalf("create notice: %v", err)
	}

	// t2 查不到 t1 班的通知。
	if _, err := svc.FindInScope(&t2, n.ID); err == nil {
		t.Fatalf("expected not-found for cross-class notice")
	}

	// t1 可以看到自己的通知。
	found, err := svc.FindInScope(&t1, n.ID)
	if err != nil {
		t.Fatalf("find own notice: %v", err)
	}
	if found.ID != n.ID {
		t.Fatalf("found wrong notice")
	}
}

// 发布 / 撤回后字段正确。
func TestNoticePublishUnpublish(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewNoticeService(db, services.NewScope(db))
	n, err := svc.Create(&teacher, "标题", "内容", "info")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	published, err := svc.Publish(n)
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	if !published.IsPublished || published.PublishedAt == nil {
		t.Fatalf("publish did not set fields")
	}

	unpublished, err := svc.Unpublish(n)
	if err != nil {
		t.Fatalf("unpublish: %v", err)
	}
	if unpublished.IsPublished {
		t.Fatalf("unpublish failed")
	}
}

// 更新可编辑字段。
func TestNoticeUpdate(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewNoticeService(db, services.NewScope(db))
	n, err := svc.Create(&teacher, "旧标题", "旧内容", "info")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	updated, err := svc.Update(n, map[string]any{"title": "新标题", "content": "新内容"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Title != "新标题" || updated.Content != "新内容" {
		t.Fatalf("update not persisted, got %+v", updated)
	}
}

// 删除通知。
func TestNoticeDelete(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewNoticeService(db, services.NewScope(db))
	n, err := svc.Create(&teacher, "标题", "内容", "info")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := svc.Delete(n); err != nil {
		t.Fatalf("delete: %v", err)
	}

	if _, err := svc.FindInScope(&teacher, n.ID); err == nil {
		t.Fatalf("expected not-found after delete")
	}
}

// List 只含本班通知。
func TestNoticeListScoped(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	t1 := seedTeacher(t, db, school.ID, "teacher1")
	seedTeacherClass(t, db, school.ID, t1.ID, "一班")

	svc := services.NewNoticeService(db, services.NewScope(db))
	if _, err := svc.Create(&t1, "A", "a", "info"); err != nil {
		t.Fatalf("create A: %v", err)
	}
	if _, err := svc.Create(&t1, "B", "b", "info"); err != nil {
		t.Fatalf("create B: %v", err)
	}

	list, err := svc.List(&t1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("list len = %d, want 2", len(list))
	}
	// 倒序：B 在前。
	if list[0].Title != "B" {
		t.Fatalf("list not newest-first, got %+v", list[0])
	}

}
