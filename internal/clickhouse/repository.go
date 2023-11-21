package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/zain-saqer/twitch-chat-archive/internal/chat"
)

func NewRepository(conn driver.Conn, database string) *LogRepository {
	return &LogRepository{
		conn:     conn,
		database: database,
	}
}

type LogRepository struct {
	conn     driver.Conn
	database string
}

func (r LogRepository) PrepareDatabase(ctx context.Context) error {
	createQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s.message
	(
	   id           String,
	   channel      LowCardinality(String),
	   username     String,
	   message      String,
	   timestamp    DateTime,
	   message_type UInt8
	)
	ENGINE = MergeTree()
	PRIMARY KEY (channel, username, timestamp);
`, r.database)

	if err := r.conn.Exec(ctx, createQuery); err != nil {
		return err
	}

	return nil
}

func (r LogRepository) SaveAll(ctx context.Context, messages []*chat.Message) error {
	for _, message := range messages {
		if err := r.conn.AsyncInsert(ctx, `INSERT INTO message (id, channel, username, message, timestamp, message_type) VALUES (
			?, ?, ?, ?, ?, ?
		)`, false, message.ID, message.ChannelName, message.Username, message.Message, message.Time, message.MessageType); err != nil {
			return err
		}
	}
	return nil
}
