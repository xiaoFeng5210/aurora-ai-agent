package database

import (
	"errors"
	"testing"

	"aurora-agent/database/model"

	"gorm.io/gorm"
)

const testPointsBalanceUserID = 910000001

func TestPointsBalanceCRUD(t *testing.T) {
	_ = DeletePointsBalanceByUserID(testPointsBalanceUserID)
	t.Cleanup(func() {
		_ = DeletePointsBalanceByUserID(testPointsBalanceUserID)
	})

	created, err := CreatePointsBalance(model.PointsBalance{
		UserId:       testPointsBalanceUserID,
		BalanceAfter: 100,
		Remark:       "init",
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.Id == 0 {
		t.Fatal("expected created id > 0")
	}
	t.Logf("created: id=%d user_id=%d balance=%d remark=%q", created.Id, created.UserId, created.BalanceAfter, created.Remark)

	gotByID, err := GetPointsBalanceByID(created.Id)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if gotByID.UserId != testPointsBalanceUserID || gotByID.BalanceAfter != 100 {
		t.Fatalf("unexpected get by id result: %#v", gotByID)
	}
	t.Logf("get by id: %+v", gotByID)

	gotByUser, err := GetPointsBalanceByUserID(testPointsBalanceUserID)
	if err != nil {
		t.Fatalf("get by user id failed: %v", err)
	}
	if gotByUser.Id != created.Id {
		t.Fatalf("expected id %d, got %d", created.Id, gotByUser.Id)
	}
	t.Logf("get by user id: %+v", gotByUser)

	if err := UpdatePointsBalanceByUserID(testPointsBalanceUserID, map[string]any{
		"balance_after": 80,
		"remark":        "consume",
	}); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	updated, err := GetPointsBalanceByUserID(testPointsBalanceUserID)
	if err != nil {
		t.Fatalf("get after update failed: %v", err)
	}
	if updated.BalanceAfter != 80 || updated.Remark != "consume" {
		t.Fatalf("unexpected updated result: %#v", updated)
	}
	t.Logf("updated: %+v", updated)

	if err := DeletePointsBalanceByUserID(testPointsBalanceUserID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = GetPointsBalanceByUserID(testPointsBalanceUserID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found after delete, got %v", err)
	}
	t.Log("deleted and confirmed not found")
}
