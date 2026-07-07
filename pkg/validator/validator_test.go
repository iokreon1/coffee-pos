package validator_test

import (
	"testing"

	"github.com/iokreon1/coffee-pos/pkg/validator"
)

type TestInput struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Price int64  `json:"price" validate:"required,min=1"`
}

func TestValidator_Validate(t *testing.T) {
	v := validator.New()

	// Test 1: input valid → Validate() return nil
	t.Run("Test 1: input valid", func(t *testing.T) {
		input := TestInput{
			Name:  "John Doe",
			Email: "john@example.com",
			Price: 10000,
		}
		errs := v.Validate(input)
		if errs != nil {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	// Test 2: semua field kosong → Validate() return map dengan 3 keys: "name", "email", "price"
	t.Run("Test 2: all fields empty", func(t *testing.T) {
		input := TestInput{}
		errs := v.Validate(input)
		if errs == nil {
			t.Error("expected validation errors, got nil")
			return
		}

		expectedKeys := []string{"name", "email", "price"}
		for _, key := range expectedKeys {
			if _, ok := errs[key]; !ok {
				t.Errorf("expected error for key %q, but not found", key)
			}
		}

		if len(errs) != len(expectedKeys) {
			t.Errorf("expected exactly %d error keys, got %d: %v", len(expectedKeys), len(errs), errs)
		}
	})

	// Test 3: email format salah → error message mengandung "format email tidak valid"
	t.Run("Test 3: invalid email format", func(t *testing.T) {
		input := TestInput{
			Name:  "John Doe",
			Email: "invalid-email",
			Price: 10000,
		}
		errs := v.Validate(input)
		if errs == nil {
			t.Error("expected validation errors, got nil")
			return
		}

		emailErr, exists := errs["email"]
		if !exists {
			t.Error("expected error for key 'email', but not found")
			return
		}

		expectedSubstr := "format email tidak valid"
		if emailErr != expectedSubstr {
			t.Errorf("expected email error message to be %q, got %q", expectedSubstr, emailErr)
		}
	})

	// Test 4: name terlalu pendek (1 karakter) → error message mengandung "minimal 2 karakter"
	t.Run("Test 4: name too short", func(t *testing.T) {
		input := TestInput{
			Name:  "A",
			Email: "john@example.com",
			Price: 10000,
		}
		errs := v.Validate(input)
		if errs == nil {
			t.Error("expected validation errors, got nil")
			return
		}

		nameErr, exists := errs["name"]
		if !exists {
			t.Error("expected error for key 'name', but not found")
			return
		}

		expectedSubstr := "minimal 2 karakter"
		if nameErr != expectedSubstr {
			t.Errorf("expected name error message to be %q, got %q", expectedSubstr, nameErr)
		}
	})
}
