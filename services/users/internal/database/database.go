package database

import (
	"context"
	"log"
	"time"
	"users/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client     *mongo.Client
	driversCollection   *mongo.Collection
	passengersCollection *mongo.Collection
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

	db := client.Database("ridenow_users")
	driversCollection := db.Collection("drivers")
	passengersCollection := db.Collection("passengers")
	
	log.Println("✅ Connecté à MongoDB")
	return &Database{
		client:              client,
		driversCollection:   driversCollection,
		passengersCollection: passengersCollection,
	}, nil
}

func (db *Database) CreateDriver(ctx context.Context, driver *types.Driver) (*primitive.ObjectID, error) {
	res, err := db.driversCollection.InsertOne(ctx, driver)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	driver.ID = id
	return &id, nil
}

func (db *Database) GetDrivers(ctx context.Context, available *bool) ([]types.Driver, error) {
	filter := bson.M{}
	if available != nil && *available {
		filter = bson.M{"is_available": true}
	}

	cursor, err := db.driversCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var drivers []types.Driver
	if err = cursor.All(ctx, &drivers); err != nil {
		return nil, err
	}

	return drivers, nil
}

func (db *Database) UpdateDriverStatus(ctx context.Context, id primitive.ObjectID, isAvailable bool) error {
	_, err := db.driversCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"is_available": isAvailable}},
	)
	return err
}

func (db *Database) CreatePassenger(ctx context.Context, passenger *types.Passenger) (*primitive.ObjectID, error) {
	now := time.Now()
	passenger.CreatedAt = now
	passenger.UpdatedAt = now
	
	res, err := db.passengersCollection.InsertOne(ctx, passenger)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	passenger.ID = id
	return &id, nil
}

func (db *Database) GetPassengers(ctx context.Context) ([]types.Passenger, error) {
	cursor, err := db.passengersCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var passengers []types.Passenger
	if err = cursor.All(ctx, &passengers); err != nil {
		return nil, err
	}

	return passengers, nil
}

func (db *Database) GetPassengerByID(ctx context.Context, id primitive.ObjectID) (*types.Passenger, error) {
	var passenger types.Passenger
	err := db.passengersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&passenger)
	if err != nil {
		return nil, err
	}
	return &passenger, nil
}

func (db *Database) UpdatePassenger(ctx context.Context, id primitive.ObjectID, passenger *types.Passenger) error {
	passenger.UpdatedAt = time.Now()
	_, err := db.passengersCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": passenger},
	)
	return err
}

func (db *Database) DeletePassenger(ctx context.Context, id primitive.ObjectID) error {
	_, err := db.passengersCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
