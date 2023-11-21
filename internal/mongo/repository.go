package mongo

import (
	"context"
	"github.com/zain-saqer/twitch-chat-archive/internal/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoLogRepository(client *mongo.Client, database, collection string) *LogRepository {
	return &LogRepository{
		client:     client,
		database:   database,
		collection: collection,
	}
}

type LogRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func PrepareDatabase(ctx context.Context, client *mongo.Client, databaseName, collectionName string) error {
	database := client.Database(databaseName)
	collection := database.Collection(collectionName)
	idIndexModel := mongo.IndexModel{
		Keys:    bson.D{{"id", 1}},
		Options: options.Index().SetUnique(true),
	}
	usernameIndexModel := mongo.IndexModel{
		Keys:    bson.D{{"username", 1}},
		Options: options.Index(),
	}
	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{idIndexModel, usernameIndexModel})
	if err != nil {
		return err
	}

	return nil
}

func (r *LogRepository) SaveAll(ctx context.Context, messages []*chat.Message) error {
	if len(messages) == 0 {
		return nil
	}
	documents := make([]any, 0)
	for _, message := range messages {
		documents = append(documents, message)
	}
	tru := true
	_, err := r.client.Database(r.database).Collection(r.collection).InsertMany(ctx, documents, &options.InsertManyOptions{Ordered: &tru})
	if err != nil {
		return err
	}
	return nil
}
