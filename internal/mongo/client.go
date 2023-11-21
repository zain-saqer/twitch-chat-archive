package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewClient(ctx context.Context, host, port, username, password string, timeout time.Duration) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf(`mongodb://%s:%s/`, host, port)
	opts := options.Client().ApplyURI(uri).
		SetAuth(options.Credential{Username: username, Password: password}).
		SetServerAPIOptions(serverAPI).
		SetTimeout(timeout).
		SetBSONOptions(&options.BSONOptions{UseJSONStructTags: true}).
		SetCompressors([]string{"snappy", "zlib", "zstd"})
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}
