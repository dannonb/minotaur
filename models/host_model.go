package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Host struct {
	Id             primitive.ObjectID   `bson:"_id"`
	Hostname       string               `json:"hostname,omitempty" validate:"required"`
	Subscriber_ids []primitive.ObjectID `json:"subscriber_ids"`
	Created_at     time.Time            `json:"created_at"`
	Updated_at     time.Time            `json:"updated_at"`
	Host_id        string               `json:"host_id"`
}
