package models_test

import (
	"dinero/api/models"
	"testing"
)

func TestAccountValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		account  *models.Account
		expected bool
	}{
		{
			// passes the test because UserID is equal to 1
			name:     "GOOD_USER_ID_EQUAL_TO_1",
			account:  &models.Account{ID: 1, UserID: 1, Name: "Car Payment", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: true,
		},
		{
			// passes the test because UserID is greater than 1
			name:     "GOOD_USER_ID_GREATER_THAN_1",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Car Payment", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: true,
		},
		{
			// breaks the test because UserID cannot be less than 1
			name:     "BAD_USER_ID",
			account:  &models.Account{ID: 1, UserID: 0, Name: "Car Payment", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: false,
		},
		{
			// passes the test because Name contains only uppercase / lowercase letters and spaces
			name:     "GOOD_NAME",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Good name", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: true,
		},
		{
			// passes the test because Name contains at least 1 character
			name:     "GOOD_NAME_MINIMAL",
			account:  &models.Account{ID: 1, UserID: 2, Name: " ", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: true,
		},
		{
			// breaks the test because Name cannot be empty
			name:     "BAD_NAME_EMPTY",
			account:  &models.Account{ID: 1, UserID: 2, Name: "", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: false,
		},
		{
			// breaks the test because Name cannot have digits
			name:     "BAD_NAME_DIGITS",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test 123", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: false,
		},
		{
			// breaks the test because Name cannot have symbols
			name:     "BAD_NAME_SYMBOLS",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test #Bad", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: false,
		},
		{
			// breaks the test because Name cannot have whitespace other than spaces
			name:     "BAD_NAME_NEWLINE",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test \n", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: false,
		},
		{
			// passes the test because AccountType must be one of the values in the regex
			name:     "GOOD_TYPE",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: true,
		},
		{
			// breaks the test because AccountType must not by empty
			name:     "BAD_TYPE_EMPTY",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test", AccountType: "", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: false,
		},
		{
			// breaks the test because AccountType must be one of the values in the regex
			name:     "BAD_TYPE_VALUE",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test", AccountType: "bad", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: false,
		},
		{
			// passes the test because Minimum Payment must be less than or equal to the Full Amount
			name:     "GOOD_MINIMUM_PAYMENT",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: true,
		},
		{
			// breaks the test because Minimum Payment is greater than the Full Amount
			name:     "BAD_MINIMUM_PAYMENT_GREATER_THAN",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test", AccountType: "monthly", MinimumPayment: 500, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: false,
		},
		{
			// passes the test because Current Payment must be less than or equal to the Full Amount
			name:     "GOOD_CURRENT_PAYMENT",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: true,
		},
		{
			// breaks the test because Current Payment is greater than the Full Amount
			name:     "BAD_CURRENT_PAYMENT_GREATER_THAN",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 500, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: false,
		},
		{
			// passes the test because Date must be between 1 and 31
			name:     "GOOD_DATE",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "1", URL: ""},
			expected: true,
		},
		{
			// breaks the test because Date must not be empty
			name:     "BAD_DATE_EMPTY",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "", URL: ""},
			expected: false,
		},
		{
			// breaks the test because Date must be between 1 and 31
			name:     "BAD_DATE_GREATER",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "32", URL: ""},
			expected: false,
		},
		{
			// breaks the test because Date must be between 1 and 31
			name:     "BAD_DATE_LESS",
			account:  &models.Account{ID: 1, UserID: 2, Name: "Test", AccountType: "monthly", MinimumPayment: 219.95, CurrentPayment: 219.95, FullAmount: 219.95, DueDate: "0", URL: ""},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.account.Validate()

			if test.expected != result {
				t.Errorf("expected: %t, got: %t", test.expected, result)
			}
		})
	}
}
