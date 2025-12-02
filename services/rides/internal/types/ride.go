package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ride struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PassengerID   string             `bson:"passenger_id" json:"passengerId"`
	DriverID      string             `bson:"driver_id" json:"driverId"`
	FromZone      string             `bson:"from_zone" json:"from_zone"`
	ToZone        string             `bson:"to_zone" json:"to_zone"`
	Price         float64            `bson:"price" json:"price"`
	Status        string             `bson:"status" json:"status"`
	PaymentStatus string             `bson:"payment_status" json:"paymentStatus"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}

