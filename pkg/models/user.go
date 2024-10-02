package models

import (
	"context"
	"github.com/oklog/ulid/v2"
	"go-identity/pkg/requests"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID            string     `json:"id"`
	FirstName     string     `json:"firstName"`
	LastName      string     `json:"lastName"`
	Email         string     `json:"email"`
	Password      string     `json:"password"`
	MfaSecret     *string    `json:"mfaSecret"`
	MfaVerifiedAt *time.Time `json:"mfaVerifiedAt"`
}

func (m *Model) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	row := m.Conn.QueryRow(ctx, "select id, firstName, lastName, email, password, mfa_secret, mfa_verified_at FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.MfaSecret, &user.MfaVerifiedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m *Model) GetUserByToken(ctx context.Context, token string) (User, error) {
	cToken, err := m.FindToken(ctx, token)
	if err != nil {
		return User{}, err
	}

	lastUsedAt := time.Now()
	cToken.LastUsedAt = &lastUsedAt
	err = m.UpdateToken(ctx, cToken)
	if err != nil {
		return User{}, err
	}

	user, err := m.GetUserById(ctx, cToken.UserID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m *Model) GetUserById(ctx context.Context, userID string) (User, error) {
	var user User
	row := m.Conn.QueryRow(ctx, "select id, firstName, lastName, email, password, mfa_secret, mfa_verified_at FROM users WHERE id = $1", userID)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.MfaSecret, &user.MfaVerifiedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m *Model) InsertIntoUsers(ctx context.Context, request requests.SignUp) (User, error) {
	userID := ulid.Make().String()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	if err != nil {
		return User{}, err
	}

	sql := "insert into users(id, firstName, lastName, email, password) values ($1, $2, $3, $4, $5)"
	_, err = m.Conn.Exec(ctx, sql, userID, request.FirstName, request.LastName, request.Email, string(passwordHash))

	if err != nil {
		return User{}, err
	}

	user, err := m.GetUserById(ctx, userID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m *Model) UpdateUser(ctx context.Context, user User) error {
	sql := "update users set firstName = $1, lastName = $2, mfa_secret = $3, mfa_verified_at = $4, updated_at = $5 where id = $6"
	_, err := m.Conn.Exec(ctx, sql, user.FirstName, user.LastName, user.MfaSecret, user.MfaVerifiedAt, time.Now(), user.ID)
	if err != nil {
		return err
	}

	return nil
}
