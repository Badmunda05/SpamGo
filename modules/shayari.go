package modules

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"pbxgo/config"

	"github.com/amarnathcjd/gogram/telegram"
)

/*
	PbxGo Shayari Module
	Created By: BadMunda
*/

// ─────────────────────────────────────────────
// Active Shayari Tracker — stopped by .stopraid
// ─────────────────────────────────────────────

var (
	activeShayariChats   = make(map[int64]bool)
	activeShayariChatsMu sync.Mutex
)

func setShayariStopped(chatID int64) {
	activeShayariChatsMu.Lock()
	delete(activeShayariChats, chatID)
	activeShayariChatsMu.Unlock()
}

func isShayariActive(chatID int64) bool {
	activeShayariChatsMu.Lock()
	defer activeShayariChatsMu.Unlock()
	return activeShayariChats[chatID]
}

// ─────────────────────────────────────────────
// .shayari
// replyToID captured once before the loop
// works with or without a reply
// ─────────────────────────────────────────────

func shayariHandler(m *telegram.NewMessage) error {
	var count int
	fmt.Sscanf(GetArgs(m), "%d", &count)
	if count < 1 {
		count = 1
	}

	// Capture replyToID once, before the loop
	var replyToID int32
	if m.IsReply() {
		replyMsg, err := m.GetReplyMessage()
		if err == nil && replyMsg != nil {
			replyToID = int32(replyMsg.ID)
		}
	}

	_, _ = m.Delete()

	activeShayariChatsMu.Lock()
	activeShayariChats[m.ChatID()] = true
	activeShayariChatsMu.Unlock()
	defer setShayariStopped(m.ChatID())

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	list := config.SHAYARI

	for i := 0; i < count; i++ {
		if !isShayariActive(m.ChatID()) {
			break
		}
		shayari := list[rng.Intn(len(list))]
		opts := &telegram.SendOptions{ParseMode: telegram.HTML}
		if replyToID > 0 {
			opts.ReplyID = replyToID
		}
		_, err := m.Client.SendMessage(m.ChatID(), shayari, opts)
		if err != nil {
			errText := strings.ToUpper(err.Error())
			if strings.Contains(errText, "FLOOD_WAIT") {
				var wait int
				fmt.Sscanf(errText, "FLOOD_WAIT_%d", &wait)
				if wait <= 0 {
					wait = 5
				}
				time.Sleep(time.Duration(wait+1) * time.Second)
				i--
				continue
			}
			time.Sleep(2 * time.Second)
			continue
		}
		time.Sleep(700 * time.Millisecond)
	}
	return nil
}

// ─────────────────────────────────────────────
// REGISTER
// ─────────────────────────────────────────────

func init() {
	Register(ModuleInfo{
		Name:        "Shayari",
		Description: "Send random shayaris",
		Commands: []CommandInfo{
			{Pattern: "shayari", Handler: shayariHandler, Sudo: true},
		},
	})
}
