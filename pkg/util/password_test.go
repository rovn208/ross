package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "secretPassword12345"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	password2 := "secretPassword123452"
	hashedPassword2, err := HashPassword(password2)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)

	require.NotEqual(t, hashedPassword, hashedPassword2)
}

func TestCheckPassword(t *testing.T) {
	password := "secretPassword12345"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := "wrongPassword12345"
	err = CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
