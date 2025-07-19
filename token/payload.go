package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	UserId    string    `json:"user_id"`
	UserPlan  string    `json:"user_plan"`
	CreatedAt time.Time `json:"created_at"`
	ExpiryAt  time.Time `json:"expiry_at"`
}

func NewPayload(username, userId, userPlan string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	
	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		UserId:    userId,
		UserPlan:  userPlan,
		CreatedAt: time.Now(),
		ExpiryAt:  time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiryAt) {
		return errors.New("token has expired")
	}
	return nil
}
