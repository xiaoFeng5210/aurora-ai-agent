package points

import (
	"aurora-agent/database"
	"aurora-agent/database/model"
	"aurora-agent/handler/vo"
	"aurora-agent/utils/enum"
	"errors"
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"
)

func createPointsTestUser(t *testing.T, suffix string) model.User {
	t.Helper()
	conn, err := database.DBConnect()
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("points-%s-%d", suffix, time.Now().UnixNano())
	user := model.User{Username: name, Email: name + "@example.com", Password: "test"}
	if err := conn.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Where("user_id = ?", user.Id).Delete(&model.PointsRecord{})
		conn.Where("user_id = ?", user.Id).Delete(&model.PointsBalance{})
		conn.Unscoped().Delete(&user)
	})
	return user
}

func TestLegacyBalanceAndCredits(t *testing.T) {
	user := createPointsTestUser(t, "credit")
	if got, err := Query(user.Id); err != nil || got != 0 {
		t.Fatalf("legacy balance = %d, %v", got, err)
	}
	if _, err := database.GetPointsBalanceByUserID(user.Id); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("read must not initialize: %v", err)
	}
	for _, amount := range []int{0, -1} {
		if _, err := Add(user.Id, amount, "invalid", enum.TriggerMode_Register); !errors.Is(err, gorm.ErrInvalidData) {
			t.Fatalf("invalid credit: %v", err)
		}
	}
	for _, want := range []int{15, 30} {
		got, err := Add(user.Id, 15, "test credit", enum.TriggerMode_Register)
		if err != nil || got.BalanceAfter != want {
			t.Fatalf("credit = %+v, %v; want %d", got, err, want)
		}
	}
}

func TestDeductPoints(t *testing.T) {
	user := createPointsTestUser(t, "deduct")

	if _, err := Deduct(user.Id, 1, "empty wallet", enum.TriggerMode_Admin); !errors.Is(err, vo.ErrInsufficientPoints) {
		t.Fatalf("empty deduct: %v", err)
	}
	if _, err := Add(user.Id, 20, "seed", enum.TriggerMode_Admin); err != nil {
		t.Fatal(err)
	}
	got, err := Deduct(user.Id, 7, "admin deduct", enum.TriggerMode_Admin)
	if err != nil || got.BalanceAfter != 13 {
		t.Fatalf("deduct = %+v, %v", got, err)
	}
	if _, err := Deduct(user.Id, 14, "too much", enum.TriggerMode_Admin); !errors.Is(err, vo.ErrInsufficientPoints) {
		t.Fatalf("overdraft: %v", err)
	}
	if _, err := Deduct(0, 1, "bad user", enum.TriggerMode_Admin); !errors.Is(err, gorm.ErrInvalidData) {
		t.Fatalf("invalid deduct: %v", err)
	}
	if _, err := QueryByUserID(2_100_000_000); !errors.Is(err, vo.ErrUserNotFound) {
		t.Fatalf("missing user query: %v", err)
	}
}
