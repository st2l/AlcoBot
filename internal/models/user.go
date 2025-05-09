package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DrinkEvent struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TelegramID       int64              `bson:"telegram_id" json:"telegram_id"`
	FirstName        string             `bson:"first_name" json:"first_name"`
	LastSoberResetAt time.Time          `bson:"last_sober_reset_at" json:"last_sober_reset_at"`
	DrinkEvents      []DrinkEvent       `bson:"drink_events" json:"drink_events"`
}
