package tokens

import (
	"context"
	"errors"
	"fmt"
	"github.com/thanhpk/randstr"
	"go-identity/pkg/models"
	"hash/crc32"
	"time"
)

func generateTokenString() string {
	tokenEntropy := randstr.String(40)
	crc32bHash := crc32.ChecksumIEEE([]byte(tokenEntropy))

	return fmt.Sprintf("%s%x", tokenEntropy, crc32bHash)
}

func CreateToken(m models.Model, userID string) (string, error) {
	expires := time.Now().Add(24 * time.Hour * 7)
	tokenString := generateTokenString()
	request := models.CreateTokenRequest{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: expires,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err := m.InsertIntoTokens(ctx, request)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckTokenValidity(ctx context.Context, m models.Model, tokenString string) error {
	cToken, err := m.FindToken(ctx, tokenString)
	if err != nil {
		return err
	}

	if time.Now().After(cToken.ExpiresAt) {
		return errors.New("expired token")
	}

	lastUsed := time.Now()
	cToken.LastUsedAt = &lastUsed
	cToken.UpdatedAt = time.Now()
	err = m.UpdateToken(ctx, cToken)
	if err != nil {
		return err
	}

	return nil
}
