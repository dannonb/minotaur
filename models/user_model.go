package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id 		    primitive.ObjectID  `bson:"_id"`
	Username    *string             `json:"username,omitempty" validate:"required"`
	Email    	*string             `json:"email,omitempty" validate:"required"`
	Password    *string             `json:"password,omitempty" validate:"required"`
	Token         *string            `json:"token"`
    Refresh_token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
    Updated_at    time.Time          `json:"updated_at"`
    User_id       string             `json:"user_id"`
}