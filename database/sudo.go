package database

import (
	"context"
	"log/slog"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var sudoUsers sync.Map

type sudoEntry struct {
	UserID int64 `bson:"user_id"`
}

// LoadSudoUsers loads all sudo users from MongoDB into memory on startup.
func LoadSudoUsers() {
	if !IsConnected() {
		return
	}
	ctx := context.Background()
	cursor, err := SudoCollection.Find(ctx, bson.D{})
	if err != nil {
		slog.Error("Failed to load sudo users", "error", err)
		return
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var e sudoEntry
		if err := cursor.Decode(&e); err == nil {
			sudoUsers.Store(e.UserID, true)
			count++
		}
	}
	slog.Info("Sudo users loaded", "count", count)
}

// AddSudo adds a user to the sudo list in memory and MongoDB.
func AddSudo(userID int64) {
	sudoUsers.Store(userID, true)
	if IsConnected() {
		ctx := context.Background()
		filter := bson.D{{Key: "user_id", Value: userID}}
		update := bson.D{{Key: "$setOnInsert", Value: bson.D{{Key: "user_id", Value: userID}}}}
		opts := options.UpdateOne().SetUpsert(true)
		_, _ = SudoCollection.UpdateOne(ctx, filter, update, opts)
	}
}

// RemoveSudo removes a user from the sudo list in memory and MongoDB.
func RemoveSudo(userID int64) {
	sudoUsers.Delete(userID)
	if IsConnected() {
		ctx := context.Background()
		_, _ = SudoCollection.DeleteOne(ctx, bson.D{{Key: "user_id", Value: userID}})
	}
}

// IsSudo returns true if the user is in the sudo list.
func IsSudo(userID int64) bool {
	_, ok := sudoUsers.Load(userID)
	return ok
}

// FetchSudoList returns all current sudo user IDs.
func FetchSudoList() []int64 {
	var list []int64
	sudoUsers.Range(func(key, _ any) bool {
		list = append(list, key.(int64))
		return true
	})
	return list
}
