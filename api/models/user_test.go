package models_test

import (
	"dinero/api/models"
	"testing"
)

func TestUserValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		user     *models.User
		expected bool
	}{
		{
			// passes the test because the FirstName contains only uppercase / lowercase letters and spaces
			name:     "GOOD_FIRST_NAME",
			user:     &models.User{ID: 1, FirstName: "John ", LastName: "Ide", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: true,
		},
		{
			// breaks the test because the FirstName cannot be empty
			name:     "BAD_FIRST_NAME_EMPTY",
			user:     &models.User{ID: 1, FirstName: "", LastName: "Ide", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// breaks the test because the FirstName cannot contain digits
			name:     "BAD_FIRST_NAME_DIGITS",
			user:     &models.User{ID: 1, FirstName: "Test 123", LastName: "Ide", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// breaks the test because the FirstName cannot contain symbols
			name:     "BAD_FIRST_NAME_SYMBOLS",
			user:     &models.User{ID: 1, FirstName: "Test #", LastName: "Ide", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// passes the test because the LastName contains only uppercase / owercase letters and spaces
			name:     "GOOD_LAST_NAME",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide ", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: true,
		},
		{
			// breaks the test because the LastName cannot be empty
			name:     "BAD_LAST_NAME_EMPTY",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// breaks the test because the LastName cannot contain digits
			name:     "BAD_LAST_NAME_DIGITS",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Test 123", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// breaks the test because the LastName cannot contain symbols
			name:     "BAD_LAST_NAME_SYMBOLS",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Test #", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// passes the test because the FullName contains only uppercase / owercase letters and spaces
			name:     "GOOD_FULL_NAME",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide ", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: true,
		},
		{
			// breaks the test because the FullName cannot be empty
			name:     "BAD_FULL_NAME_EMPTY",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// breaks the test because the FullName cannot contain digits
			name:     "BAD_FULL_NAME_DIGITS",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "Test 123", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// breaks the test because the FullName cannot contain symbols
			name:     "BAD_FULL_NAME_SYMBOLS",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "Test #", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// passes the test because the Email matches the email regex
			name:     "GOOD_EMAIL",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: true,
		},
		{
			// breaks the test because the Email cannot by empty
			name:     "BAD_EMAIL_EMPTY",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "John Ide", Email: "", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// breaks the test because the Email doesn't have a user (before @)
			name:     "BAD_EMAIL_NO_USER",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "John Ide", Email: "@gmail.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// breaks the test because the Email has symbols in the user
			name:     "BAD_EMAIL_USER_SYMBOLS",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "John Ide", Email: "ide.#johnc@gmail.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// breaks the test because the Email has symbols in the domain
			name:     "BAD_EMAIL_DOMAIN_SYMBOLS",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "John Ide", Email: "ide.johnc@gmail#.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// breaks the test because the Email doesn't contain an @
			name:     "BAD_EMAIL_NO_@",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "John Ide", Email: "ide.johncgmail.com", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// breaks the test because the Email doesn't have a domain
			name:     "BAD_EMAIL_NO_DOMAIN",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "John Ide", Email: "ide.johnc@", BiweeklyIncome: 1860.99},
			expected: false,
		},
		{
			// breaks the test because the Email domain doesn't have an extension
			name:     "BAD_EMAIL_NO_DOMAIN_EXTENSION",
			user:     &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "John Ide", Email: "ide.johnc@gmail", BiweeklyIncome: 1860.99},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.user.Validate()

			if test.expected != result {
				t.Errorf("\nExpected: \t%t\nGot: \t%t", test.expected, result)
			}
		})
	}
}
