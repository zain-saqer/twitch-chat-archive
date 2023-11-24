package clickhouse

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
	"github.com/zain-saqer/twitch-chat-archive/internal/chat"
	"strings"
)

func NewRepository(conn driver.Conn, database string) *ChatRepository {
	return &ChatRepository{
		conn:     conn,
		database: database,
	}
}

type ChatRepository struct {
	conn     driver.Conn
	database string
}

func (r ChatRepository) PrepareDatabase(ctx context.Context) error {
	createMessageQuery := `
	CREATE TABLE IF NOT EXISTS message
	(
	   channel      LowCardinality(String),
	   username     String CODEC(ZSTD),
	   message      String CODEC(ZSTD),
	   timestamp    DateTime CODEC(Delta, ZSTD),
	   message_type UInt8 CODEC(Delta, ZSTD)
	)
	ENGINE = MergeTree()
	PRIMARY KEY (channel, username, timestamp);
`
	if err := r.conn.Exec(ctx, createMessageQuery); err != nil {
		return err
	}
	createChannelQuery := `
	CREATE TABLE IF NOT EXISTS channel
	(
	    id        UUID,
	    name      String,
	    timestamp DateTime
	)
	ENGINE = MergeTree()
	PRIMARY KEY (id);
`
	if err := r.conn.Exec(ctx, createChannelQuery); err != nil {
		return err
	}

	return nil
}

func (r ChatRepository) GetMessages(ctx context.Context, channel, username string, limit, offset int) ([]*chat.Message, error) {
	messages := make([]*chat.Message, 0)
	rows, err := r.conn.Query(ctx, `SELECT channel, username, message, timestamp, message_type FROM message WHERE (channel = ?) AND (username = ?) ORDER BY timestamp LIMIT ? OFFSET ?`, strings.ToLower(channel), strings.ToLower(username), limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		message := &chat.Message{}
		if err := rows.ScanStruct(message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r ChatRepository) SaveMessage(ctx context.Context, message *chat.Message) error {
	if err := r.conn.AsyncInsert(ctx, `INSERT INTO message (channel, username, message, timestamp, message_type) VALUES (?, ?, ?, ?, ?)`, false, message.ChannelName, message.Username, message.Message, message.Time, message.MessageType); err != nil {
		return err
	}
	return nil
}

func (r ChatRepository) GetChannels(ctx context.Context) ([]*chat.Channel, error) {
	channels := make([]*chat.Channel, 0)
	rows, err := r.conn.Query(ctx, `SELECT * FROM channel ORDER BY timestamp`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		channel := &chat.Channel{}
		if err := rows.ScanStruct(channel); err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}

	return channels, nil
}

func (r ChatRepository) GetChannel(ctx context.Context, uuid uuid.UUID) (*chat.Channel, error) {
	row := r.conn.QueryRow(ctx, `SELECT * FROM channel WHERE id = ? ORDER BY id LIMIT 1`, uuid)
	if row.Err() != nil {
		return nil, row.Err()
	}
	channel := &chat.Channel{}
	if err := row.ScanStruct(channel); err != nil {
		return nil, err
	}

	return channel, nil
}

func (r ChatRepository) SaveChannel(ctx context.Context, channel *chat.Channel) error {
	if err := r.conn.Exec(ctx, `INSERT INTO channel (id, name, timestamp) VALUES (?, ?, ?)`, channel.ID, channel.Name, channel.Time); err != nil {
		return err
	}
	return nil
}

func (r ChatRepository) DeleteChannel(ctx context.Context, channel *chat.Channel) error {
	if err := r.conn.Exec(ctx, `ALTER TABLE channel DELETE WHERE id = ?`, channel.ID); err != nil {
		return err
	}
	return nil
}
