package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/st2l/AlcoBot/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient struct {
	client   *mongo.Client
	database *mongo.Database
	timeout  time.Duration
}

func NewMongoClient(uri string, databaseName string, timeout time.Duration) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	// ping for the server
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	database := client.Database(databaseName)

	return &MongoClient{
		client:   client,
		database: database,
		timeout:  timeout * time.Second,
	}, nil
}

func (mc *MongoClient) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), mc.timeout)
	defer cancel()

	err := mc.client.Disconnect(ctx)
	return err
}

func (mc *MongoClient) InitializeCollection(name string) error {
	newCollection := mc.database.Collection(name)

	ctx, cancel := context.WithTimeout(context.Background(), mc.timeout)
	defer cancel()

	// check for creating test document
	_, err := newCollection.InsertOne(ctx, bson.D{
		{Key: "Test Document", Value: "Test value for test document"},
	})
	if err != nil {
		log.Fatalln(err)
		return err
	}

	// check for deleting the test document
	filter := bson.D{{Key: "Test Document", Value: "Test value for test document"}}
	_, err = newCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

func (mc *MongoClient) DeleteCollection(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), mc.timeout)
	defer cancel()

	return mc.database.Collection(name).Drop(ctx)
}

func (mc *MongoClient) CheckWorkingGroup(name string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mc.timeout)
	defer cancel()

	collections, err := mc.database.ListCollectionNames(ctx, bson.M{
		"name": name,
	})
	if err != nil {
		return false, err
	}

	return len(collections) > 0, nil
}

func (mc *MongoClient) GetOrCreateUser(collection_name string, user_telegram_id int64, firstName string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mc.timeout)
	defer cancel()

	collection := mc.database.Collection(collection_name)

	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{
		"telegram_id": user_telegram_id,
	}).Decode(&existingUser)
	// user found!!!
	if err == nil {
		return &existingUser, nil
	}

	// check if error is NOT that we just do not found user
	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	// user not found - creating new user
	now := time.Now()
	newUser := models.User{
		TelegramID:       user_telegram_id,
		FirstName:        firstName,
		LastSoberResetAt: now,
		DrinkEvents:      []models.DrinkEvent{},
	}

	resultInsertOne, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// set the generated ObjectID in the user struct
	newUser.ID = resultInsertOne.InsertedID.(primitive.ObjectID)

	return &newUser, nil
}

func (mc *MongoClient) AddDrinkEvent(group_id string, user_telegram_id int64) (*models.DrinkEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mc.timeout)
	defer cancel()

	userFilter := bson.M{
		"telegram_id": user_telegram_id,
	}
	drinkEvent := models.DrinkEvent{
		Timestamp: time.Now(),
	}
	collection := mc.database.Collection(group_id)

	_, err := collection.UpdateOne(ctx, userFilter, bson.M{
		"$push": bson.M{
			"drink_events": drinkEvent,
		},
		"$set": bson.M{
			"last_sober_reset_at": time.Now(),
		},
	})
	return &drinkEvent, err
}

func (mc *MongoClient) ListAllUsers(group_id string) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mc.timeout)
	defer cancel()

	collection := mc.database.Collection(group_id)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (mc *MongoClient) GetUserByTelegramID(collection_name string, telegram_id int64) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mc.timeout)
	defer cancel()

	collection := mc.database.Collection(collection_name)

	filter := bson.M{
		"telegram_id": telegram_id,
	}

	var user models.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
