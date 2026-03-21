package service

import "testing"

func TestParseBirthday(t *testing.T) {
	t.Run("empty birthday", func(t *testing.T) {
		birthday, err := parseBirthday("")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if birthday != nil {
			t.Fatalf("expected nil birthday, got %v", birthday)
		}
	})

	t.Run("valid birthday", func(t *testing.T) {
		birthday, err := parseBirthday("2024-10-01")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if birthday == nil {
			t.Fatal("expected parsed birthday, got nil")
		}
		if birthday.Format(birthdayLayout) != "2024-10-01" {
			t.Fatalf("unexpected birthday value: %s", birthday.Format(birthdayLayout))
		}
	})

	t.Run("invalid birthday", func(t *testing.T) {
		_, err := parseBirthday("2024/10/01")
		if err != ErrBirthdayFormat {
			t.Fatalf("expected %v, got %v", ErrBirthdayFormat, err)
		}
	})
}

func TestHashAndVerifyPassword(t *testing.T) {
	password := "aurora-secret"

	hash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("hash password failed: %v", err)
	}

	if hash == password {
		t.Fatal("expected hashed password to differ from raw password")
	}

	if err := verifyPassword(hash, password); err != nil {
		t.Fatalf("verify password failed: %v", err)
	}

	if err := verifyPassword(hash, "wrong-password"); err == nil {
		t.Fatal("expected wrong password verification to fail")
	}
}

func TestValidatePassword(t *testing.T) {
	if err := validatePassword("12345"); err != ErrPasswordTooShort {
		t.Fatalf("expected %v, got %v", ErrPasswordTooShort, err)
	}

	if err := validatePassword("123456"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
