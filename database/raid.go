package database

import (
	"context"
	"log/slog"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// RaidSession holds the state of an active raid, persisted to MongoDB.
// Count == -1 means watcher mode (replyRaid waiting for target message).
var activeRaids sync.Map

type RaidSession struct {
	ChatID       int64  `bson:"chat_id"`
	ReplyToID    int32  `bson:"reply_to_id"`    // 0 = no reply, >0 = reply to this message
	RaidType     string `bson:"raid_type"`      // "raid", "hraid", "eraid", "punraid"
	Count        int    `bson:"count"`          // -1 = watcher mode
	TargetUserID int64  `bson:"target_user_id"` // used by replyRaid watcher
}

func SaveRaid(s RaidSession) {
	activeRaids.Store(s.ChatID, s)
	if !IsConnected() {
		return
	}
	ctx := context.Background()
	col := client.Database("pbxgo").Collection("active_raids")
	filter := bson.D{{Key: "chat_id", Value: s.ChatID}}
	update := bson.D{{Key: "$set", Value: s}}
	opts := options.UpdateOne().SetUpsert(true)
	_, _ = col.UpdateOne(ctx, filter, update, opts)
}

func DeleteRaid(chatID int64) {
	activeRaids.Delete(chatID)
	if !IsConnected() {
		return
	}
	ctx := context.Background()
	col := client.Database("pbxgo").Collection("active_raids")
	_, _ = col.DeleteOne(ctx, bson.D{{Key: "chat_id", Value: chatID}})
}

// GetRaid returns the active raid session for a chat, or nil if none exists.
func GetRaid(chatID int64) *RaidSession {
	val, ok := activeRaids.Load(chatID)
	if !ok {
		return nil
	}
	s := val.(RaidSession)
	return &s
}

// LoadActiveRaids loads all persisted raid sessions from MongoDB on startup.
func LoadActiveRaids() []RaidSession {
	var list []RaidSession
	if !IsConnected() {
		return list
	}
	ctx := context.Background()
	col := client.Database("pbxgo").Collection("active_raids")
	cursor, err := col.Find(ctx, bson.D{})
	if err != nil {
		slog.Error("Failed to load active raids", "error", err)
		return list
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var s RaidSession
		if err := cursor.Decode(&s); err == nil {
			activeRaids.Store(s.ChatID, s)
			list = append(list, s)
		}
	}
	slog.Info("Active raids loaded", "count", len(list))
	return list
}
