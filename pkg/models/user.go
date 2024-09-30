package models

import (
	"context"
	"github.com/oklog/ulid/v2"
	"go-identity/pkg/requests"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (m *Model) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	row := m.Conn.QueryRow(ctx, "select id, firstName, lastname, email, password FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m *Model) GetUserById(ctx context.Context, userID string) (User, error) {
	var user User
	row := m.Conn.QueryRow(ctx, "select id, firstName, lastname, email, password FROM users WHERE id = $1", userID)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
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
