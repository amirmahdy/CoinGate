package token

import (
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	IssueAt  time.Time `json:"issue_at"`
	ExpireAt time.Time `json:"expire_at"`
}

func NewPayload(username string, duration time.Duration) *Payload {
	now := time.Now()
	return &Payload{
		ID:       uuid.New(),
		Username: username,
		IssueAt:  now,
		ExpireAt: now.Add(duration),
	}
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpireAt) {
		return ErrExpiredToken
	}
	return nil
}
