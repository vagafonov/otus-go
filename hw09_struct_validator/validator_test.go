package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	CreditCard struct {
		CVV    string `validate:"regexp:\\d+|len:3"`
		Number string `validate:"undefined:42"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr []error
	}{
		{
			User{
				ID:     "1000",
				Name:   "Test",
				Age:    30,
				Email:  "test@test.ru",
				Role:   UserRole("admin"),
				Phones: []string{"11111111111"},
				meta:   nil,
			},
			[]error{
				ErrValidationLengthString,
			},
		},
		{
			User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "Test",
				Age:    17,
				Email:  "test@test.ru",
				Role:   UserRole("admin"),
				Phones: []string{"11111111111"},
				meta:   nil,
			},
			[]error{
				ErrValidationMin,
			},
		},
		{
			User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "Test",
				Age:    51,
				Email:  "test@test.ru",
				Role:   UserRole("admin"),
				Phones: []string{"11111111111"},
				meta:   nil,
			},
			[]error{
				ErrValidationMax,
			},
		},
		{
			User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "Test",
				Age:    30,
				Email:  "test@test",
				Role:   UserRole("admin"),
				Phones: []string{"11111111111"},
				meta:   nil,
			},
			[]error{
				ErrValidationRegexp,
			},
		},
		{
			User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "Test",
				Age:    30,
				Email:  "test@test.ru",
				Role:   UserRole("guest"),
				Phones: []string{"11111111111"},
				meta:   nil,
			},
			[]error{
				ErrValidationIn,
			},
		},
		{
			User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "Test",
				Age:    30,
				Email:  "test@test.ru",
				Role:   UserRole("admin"),
				Phones: []string{"1", "2"},
				meta:   nil,
			},
			[]error{
				ErrValidationLengthSliceString,
			},
		},
		{
			App{
				Version: "1",
			},
			[]error{
				ErrValidationLengthString,
			},
		},
		{
			Response{
				Code: 300,
				Body: "{}",
			},
			[]error{
				ErrValidationIn,
			},
		},
		{
			CreditCard{
				CVV: "ab", // Multiple errors for one field
			},
			[]error{
				ErrValidationRegexp,
				ErrValidationLengthString,
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			result, err := Validate(tt.in)
			if err != nil {
				fmt.Println(err)
			}

			if len(result) != len(tt.expectedErr) {
				require.Failf(
					t,
					"Test fail",
					"Number of errors does not match. Expected errors: %v (%v) Result errors: %v (%v) ",
					len(tt.expectedErr),
					tt.expectedErr,
					len(result),
					result,
				)
			}

			for k, v := range result {
				require.ErrorIs(t, v.Err, tt.expectedErr[k])
			}
		})
	}
}

func TestWithoutValidationErrors(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr []error
	}{
		{
			User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "Test",
				Age:    30,
				Email:  "test@test.ru",
				Role:   UserRole("admin"),
				Phones: []string{"11111111111", "22222222222"},
				meta:   nil,
			},
			[]error{},
		},
		{
			App{
				Version: "12345",
			},
			[]error{},
		},
		{
			Token{
				Header:    []byte{1},
				Payload:   []byte{2},
				Signature: []byte{3},
			},
			[]error{},
		},
		{
			Response{
				Code: 200,
				Body: "{}",
			},
			[]error{},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			result, err := Validate(tt.in)
			if err != nil {
				fmt.Println(err)
			}

			if len(result) != len(tt.expectedErr) {
				require.Failf(
					t,
					"Test fail",
					"Number of errors does not match. Expected errors: %v (%v) Result errors: %v (%v)",
					len(tt.expectedErr),
					tt.expectedErr,
					len(result),
					result,
				)
			}

			for k, v := range result {
				require.ErrorIs(t, v.Err, tt.expectedErr[k])
			}
		})
	}
}

func TestErrors(t *testing.T) {
	t.Run("is structure", func(t *testing.T) {
		result, err := Validate(1)
		if err != nil {
			fmt.Println(err)
		}
		_ = result

		require.ErrorContains(t, err, "input value myst be structure")
	})
}
