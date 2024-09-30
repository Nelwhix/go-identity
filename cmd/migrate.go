package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://nelwhix:admin@localhost:5432/go_identity")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id CHAR(26) PRIMARY KEY,
		firstName VARCHAR(255) NOT NULL,
		lastName VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	createTokensTable := `
	CREATE TABLE IF NOT EXISTS personal_access_tokens (
		id SERIAL PRIMARY KEY,
		user_id CHAR(26) NOT NULL,
		token VARCHAR(64) NOT NULL,
		last_used_at TIMESTAMP null,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	addMfaMigrations := `
	ALTER TABLE users ADD mfa_secret VARCHAR(255) NULL,
	ADD mfa_verified_at TIMESTAMP NULL,
	ADD mfa_recovery_codes TEXT NULL
	`

	tableQueries := []string{createUsersTable, createTokensTable, addMfaMigrations}

	for _, tableQuery := range tableQueries {
		_, err = conn.Exec(context.Background(), tableQuery)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to execute query: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Migrations ran successfully!")
}
