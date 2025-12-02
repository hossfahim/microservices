package database

import (
	"context"
	"log"
	"time"
	"rides/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client         *mongo.Client
	ridesCollection *mongo.Collection
}

func InitMongoDB(mongoURI string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var err error
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database("ridenow_rides")
	ridesCollection := db.Collection("rides")

	log.Println("✅ Connecté à MongoDB")
	return &Database{
		client:          client,
		ridesCollection: ridesCollection,
	}, nil
}

func (db *Database) CreateRide(ctx context.Context, ride *types.Ride) (*primitive.ObjectID, error) {
	res, err := db.ridesCollection.InsertOne(ctx, ride)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	ride.ID = id
	return &id, nil
}

func (db *Database) GetRideByID(ctx context.Context, id primitive.ObjectID) (*types.Ride, error) {
	var ride types.Ride
	err := db.ridesCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&ride)
	if err != nil {
		return nil, err
	}
	return &ride, nil
}

func (db *Database) UpdateRideStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	_, err := db.ridesCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		}},
	)
	return err
}

func (db *Database) UpdateRidePaymentStatus(ctx context.Context, id primitive.ObjectID, paymentStatus string) error {
	_, err := db.ridesCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"payment_status": paymentStatus,
			"updated_at":     time.Now(),
		}},
	)
	return err
}

func (db *Database) UpdateRide(ctx context.Context, id primitive.ObjectID, ride *types.Ride) error {
	ride.UpdatedAt = time.Now()
	_, err := db.ridesCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": ride},
	)
	return err
}

