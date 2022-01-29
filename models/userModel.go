package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id"`
	Email *string `json:"email" validate:"email,required"`
	Password *string `json:"password" validate:"required,min=6"`
	Token *string `json:"token"`
	RefreshToken *string `json:"refreshToken"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	User_id	string `json:"user_id"`
}