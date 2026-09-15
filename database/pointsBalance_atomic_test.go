package database

import (
	"aurora-agent/database/model"
	"aurora-agent/utils/enum"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestPointsRegistrationAtomic(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(fmt.Sprintf("rollback=%t", fail), func(t *testing.T) {
			name := fmt.Sprintf("points-qa-%d", time.Now().UnixNano())
			user := model.User{Username: name, Email: name + "@example.com", Password: "test"}
			balance := model.PointsBalance{Remark: "注册初始化"}
			if fail {
				balance.Remark = strings.Repeat("x", 256)
			}
			err := CreateUserWithPointsBalance(user, balance)
			if fail {
				if err == nil {
					t.Fatal("expected balance insertion to fail")
				}
				var count int64
				if err := db.Model(&model.User{}).Where("username = ?", name).Count(&count).Error; err != nil {
					t.Fatal(err)
				}
				if count != 0 {
					t.Fatal("user creation was not rolled back")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			created, err := GetUserByEmail(user.Email)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				db.Where("user_id = ?", created.Id).Delete(&model.PointsBalance{})
				db.Unscoped().Delete(&created)
			})
			got, err := GetPointsBalanceByUserID(created.Id)
			if err != nil || got.BalanceAfter != 0 {
				t.Fatalf("initial balance: %+v, %v", got, err)
			}
		})
	}
}

func TestPointsConcurrentFirstCredits(t *testing.T) {
	const userID = 910000020
	if err := db.Where("user_id = ?", userID).Delete(&model.PointsBalance{}).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Where("user_id = ?", userID).Delete(&model.PointsBalance{}) })
	var wg sync.WaitGroup
	errs := make(chan error, 20)
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			balance, err := IncrementPointsBalance(userID, 5, "test credit", enum.TriggerMode_Register)
			if err == nil && (balance.UserId != userID || balance.BalanceAfter < 5) {
				err = fmt.Errorf("invalid returned balance: %+v", balance)
			}
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	got, err := GetPointsBalanceByUserID(userID)
	if err != nil || got.BalanceAfter != 100 {
		t.Fatalf("concurrent balance: %+v, %v", got, err)
	}
}
