package database

import (
	"errors"
	"testing"

	"aurora-agent/database/model"
	"aurora-agent/utils/enum"

	"gorm.io/gorm"
)

// 集成测试：直连真实库。每个用例用独立 user_id，开头清掉该用户旧数据，结束后不删，方便查库。
// 单跑示例：go test ./database -run TestCreatePointsBalance -v

const (
	testPointsBalanceCreateUserID     = 910000001
	testPointsBalanceGetByIDUserID    = 910000002
	testPointsBalanceGetByUserID      = 910000003
	testPointsBalanceUpdateUserID     = 910000004
	testPointsBalanceDeleteUserID     = 910000005
	testPointsBalanceDuplicateUserID  = 910000006
	testPointsBalanceMissingGetUserID = 910000007
	testPointsBalanceMissingUpdUserID = 910000008
	testPointsBalanceMissingDelUserID = 910000009
	testPointsBalanceDecrementUserID  = 910000010
	testPointsBalanceOverdraftUserID  = 910000011
)

func resetPointsBalance(t *testing.T, userID int) {
	t.Helper()
	_ = DeletePointsBalanceByUserID(userID)
}

func seedPointsBalance(t *testing.T, userID int, balance int, remark string) model.PointsBalance {
	t.Helper()
	resetPointsBalance(t, userID)

	created, err := CreatePointsBalance(model.PointsBalance{
		UserId:       userID,
		BalanceAfter: balance,
		Remark:       remark,
	})
	if err != nil {
		t.Fatalf("seed create failed: %v", err)
	}
	if created.Id == 0 {
		t.Fatal("expected seeded id > 0")
	}
	t.Logf("seeded: id=%d user_id=%d balance=%d remark=%q", created.Id, created.UserId, created.BalanceAfter, created.Remark)
	return created
}

func TestCreatePointsBalance(t *testing.T) {
	resetPointsBalance(t, testPointsBalanceCreateUserID)

	created, err := CreatePointsBalance(model.PointsBalance{
		UserId:       testPointsBalanceCreateUserID,
		BalanceAfter: 100,
		Remark:       "init",
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.Id == 0 {
		t.Fatal("expected created id > 0")
	}
	if created.UserId != testPointsBalanceCreateUserID || created.BalanceAfter != 100 || created.Remark != "init" {
		t.Fatalf("unexpected created result: %#v", created)
	}
	t.Logf("created: id=%d user_id=%d balance=%d remark=%q", created.Id, created.UserId, created.BalanceAfter, created.Remark)
}

func TestGetPointsBalanceByID(t *testing.T) {
	created := seedPointsBalance(t, testPointsBalanceGetByIDUserID, 100, "get-by-id")

	got, err := GetPointsBalanceByID(created.Id)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if got.UserId != testPointsBalanceGetByIDUserID || got.BalanceAfter != 100 {
		t.Fatalf("unexpected get by id result: %#v", got)
	}
	t.Logf("get by id: %+v", got)
}

func TestGetPointsBalanceByUserID(t *testing.T) {
	created := seedPointsBalance(t, testPointsBalanceGetByUserID, 100, "get-by-user")

	got, err := GetPointsBalanceByUserID(testPointsBalanceGetByUserID)
	if err != nil {
		t.Fatalf("get by user id failed: %v", err)
	}
	if got.Id != created.Id {
		t.Fatalf("expected id %d, got %d", created.Id, got.Id)
	}
	t.Logf("get by user id: %+v", got)
}

func TestUpdatePointsBalanceByUserID(t *testing.T) {
	seedPointsBalance(t, testPointsBalanceUpdateUserID, 100, "before-update")

	if err := UpdatePointsBalanceByUserID(testPointsBalanceUpdateUserID, map[string]any{
		"balance_after": 80,
		"remark":        "consume",
	}); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	updated, err := GetPointsBalanceByUserID(testPointsBalanceUpdateUserID)
	if err != nil {
		t.Fatalf("get after update failed: %v", err)
	}
	if updated.BalanceAfter != 80 || updated.Remark != "consume" {
		t.Fatalf("unexpected updated result: %#v", updated)
	}
	t.Logf("updated: %+v", updated)
}

func TestDeletePointsBalanceByUserID(t *testing.T) {
	seedPointsBalance(t, testPointsBalanceDeleteUserID, 100, "to-delete")

	if err := DeletePointsBalanceByUserID(testPointsBalanceDeleteUserID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err := GetPointsBalanceByUserID(testPointsBalanceDeleteUserID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found after delete, got %v", err)
	}
	t.Log("deleted and confirmed not found")
}

func TestCreatePointsBalanceDuplicateUserID(t *testing.T) {
	seedPointsBalance(t, testPointsBalanceDuplicateUserID, 100, "first")

	_, err := CreatePointsBalance(model.PointsBalance{
		UserId:       testPointsBalanceDuplicateUserID,
		BalanceAfter: 200,
		Remark:       "second",
	})
	if err == nil {
		t.Fatal("expected duplicate user_id to fail")
	}
	t.Logf("duplicate create error: %v", err)
}

func TestGetPointsBalanceByUserIDNotFound(t *testing.T) {
	resetPointsBalance(t, testPointsBalanceMissingGetUserID)

	_, err := GetPointsBalanceByUserID(testPointsBalanceMissingGetUserID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}
	t.Log("get missing user confirmed not found")
}

func TestUpdatePointsBalanceByUserIDNotFound(t *testing.T) {
	resetPointsBalance(t, testPointsBalanceMissingUpdUserID)

	err := UpdatePointsBalanceByUserID(testPointsBalanceMissingUpdUserID, map[string]any{
		"balance_after": 1,
	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}
	t.Log("update missing user confirmed not found")
}

func TestDeletePointsBalanceByUserIDNotFound(t *testing.T) {
	resetPointsBalance(t, testPointsBalanceMissingDelUserID)

	err := DeletePointsBalanceByUserID(testPointsBalanceMissingDelUserID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}
	t.Log("delete missing user confirmed not found")
}

func TestDecrementPointsBalance(t *testing.T) {
	seedPointsBalance(t, testPointsBalanceDecrementUserID, 20, "before-debit")
	conn, err := DBConnect()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Where("user_id = ?", testPointsBalanceDecrementUserID).Delete(&model.PointsRecord{})
		resetPointsBalance(t, testPointsBalanceDecrementUserID)
	})

	got, err := DecrementPointsBalance(testPointsBalanceDecrementUserID, 8, "admin deduct", enum.TriggerMode_Admin)
	if err != nil || got.BalanceAfter != 12 {
		t.Fatalf("decrement = %+v, %v", got, err)
	}
}

func TestDecrementPointsBalanceInsufficient(t *testing.T) {
	seedPointsBalance(t, testPointsBalanceOverdraftUserID, 5, "low-balance")
	conn, err := DBConnect()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Where("user_id = ?", testPointsBalanceOverdraftUserID).Delete(&model.PointsRecord{})
		resetPointsBalance(t, testPointsBalanceOverdraftUserID)
	})

	_, err = DecrementPointsBalance(testPointsBalanceOverdraftUserID, 6, "overdraft", enum.TriggerMode_Admin)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected insufficient debit to miss the row, got %v", err)
	}
	got, err := GetPointsBalanceByUserID(testPointsBalanceOverdraftUserID)
	if err != nil || got.BalanceAfter != 5 {
		t.Fatalf("balance mutated on failed debit: %+v, %v", got, err)
	}
}
