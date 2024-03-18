package handler_test

import (
    "filmoteca/internal/handler"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
    password := "secret"
    hashedPassword, err := handler.HashPassword(password)
    require.NoError(t, err)
    require.NotEmpty(t, hashedPassword)
    require.NotEqual(t, password, hashedPassword)
}

func TestComparePasswords(t *testing.T) {
    password := "secret"
    wrongPassword := "notsecret"

    hashedPassword, err := handler.HashPassword(password)
    require.NoError(t, err)

    // Тестирование с корректным паролем
    isMatch := handler.ComparePasswords(hashedPassword, password)
    require.True(t, isMatch)

    // Тестирование с некорректным паролем
    isMatch = handler.ComparePasswords(hashedPassword, wrongPassword)
    require.False(t, isMatch)
}
