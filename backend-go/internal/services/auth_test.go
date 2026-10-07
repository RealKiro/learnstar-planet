package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedAdminUser(t *testing.T, db *gorm.DB, schoolID uint, password string) models.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user := models.User{
		SchoolID: schoolID, Role: "school_admin", Username: "admin",
		PasswordHash: string(hash), Name: "管理员", Status: "active",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return user
}

func TestLoginSuccess(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	user := seedAdminUser(t, db, school.ID, "secret123")

	jwtMgr := auth.New("test-secret", 72)
	svc := services.NewAuthService(db, jwtMgr)

	token, got, err := svc.Login("admin", "secret123")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if token == "" {
		t.Fatalf("token is empty")
	}
	if got.ID != user.ID || got.LastLoginAt == nil {
		t.Fatalf("got = %+v", got)
	}

	claims, err := jwtMgr.Parse(token)
	if err != nil {
		t.Fatalf("Parse token: %v", err)
	}
	if claims.UserID != user.ID || claims.Role != "school_admin" || claims.SchoolID != school.ID {
		t.Fatalf("claims = %+v", claims)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	seedAdminUser(t, db, school.ID, "secret123")

	jwtMgr := auth.New("test-secret", 72)
	svc := services.NewAuthService(db, jwtMgr)

	_, _, err := svc.Login("admin", "wrong-password")
	ae, ok := services.AsAppError(err)
	if !ok || ae.Status != 422 {
		t.Fatalf("Login err = %v, want 422 AppError", err)
	}
}

func TestChangePassword(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	user := seedAdminUser(t, db, school.ID, "secret123")

	svc := services.NewAuthService(db, auth.New("test-secret", 72))

	if err := svc.ChangePassword(&user, "secret123", "newpass456"); err != nil {
		t.Fatalf("ChangePassword: %v", err)
	}

	var reloaded models.User
	db.First(&reloaded, user.ID)
	if bcrypt.CompareHashAndPassword([]byte(reloaded.PasswordHash), []byte("newpass456")) != nil {
		t.Fatalf("password hash not updated")
	}
	if !reloaded.PasswordChanged {
		t.Fatalf("PasswordChanged = false, want true")
	}
}

func TestChangePasswordWrongOld(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	user := seedAdminUser(t, db, school.ID, "secret123")

	svc := services.NewAuthService(db, auth.New("test-secret", 72))
	err := svc.ChangePassword(&user, "wrong-old", "newpass456")
	ae, ok := services.AsAppError(err)
	if !ok || ae.Status != 422 {
		t.Fatalf("ChangePassword err = %v, want 422 AppError", err)
	}
}
