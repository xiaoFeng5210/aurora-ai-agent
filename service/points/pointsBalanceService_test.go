package points

import (
	"aurora-agent/database"
	"aurora-agent/database/model"
	"errors"
	"testing"

	"gorm.io/gorm"

	"aurora-agent/utils/enum"
)

func TestLegacyBalanceAndCredits(t *testing.T) {
	conn, err := database.DBConnect()
	if err != nil {
		t.Fatal(err)
	}
	user := model.User{Username: "points-service-qa", Email: "points-service-qa@example.com", Password: "test"}
	if err := conn.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Where("user_id = ?", user.Id).Delete(&model.PointsBalance{})
		conn.Unscoped().Delete(&user)
	})
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
