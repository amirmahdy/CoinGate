// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        int64        `json:"id"`
	Username  string       `json:"username"`
	Balance   int64        `json:"balance"`
	Coin      string       `json:"coin"`
	CreatedAt sql.NullTime `json:"created_at"`
}

type Coin struct {
	Name      string       `json:"name"`
	CreatedAt sql.NullTime `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID    `json:"id"`
	Username     string       `json:"username"`
	RefreshToken string       `json:"refresh_token"`
	UserAgent    string       `json:"user_agent"`
	ClientIp     string       `json:"client_ip"`
	IsBlocked    bool         `json:"is_blocked"`
	ExpiresAt    time.Time    `json:"expires_at"`
	CreatedAt    sql.NullTime `json:"created_at"`
}

type Transfer struct {
	ID            int64        `json:"id"`
	FromAccountID int64        `json:"from_account_id"`
	ToAccountID   int64        `json:"to_account_id"`
	Amount        int64        `json:"amount"`
	Coin          string       `json:"coin"`
	CreatedAt     sql.NullTime `json:"created_at"`
}

type User struct {
	Username          string       `json:"username"`
	HashedPassword    string       `json:"hashed_password"`
	FullName          string       `json:"full_name"`
	Email             string       `json:"email"`
	Active            sql.NullBool `json:"active"`
	PasswordChangedAt sql.NullTime `json:"password_changed_at"`
	CreatedAt         sql.NullTime `json:"created_at"`
}
