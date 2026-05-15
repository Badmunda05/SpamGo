package modules

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"pbxgo/config"
	"pbxgo/database"

	"github.com/amarnathcjd/gogram/telegram"
)

/*
	PbxGo Raid Module
	Created By: BadMunda
*/

// ─────────────────────────────────────────────
// Active Raid Tracker — used by stopraid
// ─────────────────────────────────────────────

var (
	activeRaidChats   = make(map[int64]bool)
	activeRaidChatsMu sync.Mutex

	replyRaidWatchers   = make(map[int64]func())
	replyRaidWatchersMu sync.Mutex
)

func setRaidActive(chatID int64) {
	activeRaidChatsMu.Lock()
	activeRaidChats[chatID] = true
	activeRaidChatsMu.Unlock()
}

func setRaidStopped(chatID int64) {
	activeRaidChatsMu.Lock()
	delete(activeRaidChats, chatID)
	activeRaidChatsMu.Unlock()
}

func isRaidActive(chatID int64) bool {
	activeRaidChatsMu.Lock()
	defer activeRaidChatsMu.Unlock()
	return activeRaidChats[chatID]
}

// ─────────────────────────────────────────────
// Normal Raid Handlers
// ─────────────────────────────────────────────

func raidHandler(m *telegram.NewMessage) error {
	return genericRaid(m, config.RAID, "raid", 700*time.Millisecond)
}
func hraidHandler(m *telegram.NewMessage) error {
	return genericRaid(m, config.HRAID, "hraid", 700*time.Millisecond)
}
func eraidHandler(m *telegram.NewMessage) error {
	return genericRaid(m, config.ERAID, "eraid", 700*time.Millisecond)
}
func punraidHandler(m *telegram.NewMessage) error {
	return genericRaid(m, config.PUNRAID, "punraid", 700*time.Millisecond)
}

// ─────────────────────────────────────────────
// Reply Raid Handlers (watcher mode)
// ─────────────────────────────────────────────

func replyRaidHandler(m *telegram.NewMessage) error {
	return startReplyRaidWatcher(m, config.RAID, "raid")
}
func hreplyRaidHandler(m *telegram.NewMessage) error {
	return startReplyRaidWatcher(m, config.HRAID, "hraid")
}
func ereplyRaidHandler(m *telegram.NewMessage) error {
	return startReplyRaidWatcher(m, config.ERAID, "eraid")
}
func preplyRaidHandler(m *telegram.NewMessage) error {
	return startReplyRaidWatcher(m, config.PUNRAID, "punraid")
}

// ─────────────────────────────────────────────
// .stopraid — stops all raids, watchers, shayari
// ─────────────────────────────────────────────

func stopRaidHandler(m *telegram.NewMessage) error {
	chatID := m.ChatID()

	setRaidStopped(chatID)

	replyRaidWatchersMu.Lock()
	if cancel, ok := replyRaidWatchers[chatID]; ok {
		cancel()
		delete(replyRaidWatchers, chatID)
	}
	replyRaidWatchersMu.Unlock()

	database.DeleteRaid(chatID)
	setShayariStopped(chatID)

	_, _ = m.Delete()
	Reply(m, "🛑 <b>ᴀʟʟ ʀᴀɪᴅs sᴛᴏᴘᴘᴇᴅ!</b>")
	return nil
}

// ─────────────────────────────────────────────
// Flood Wait Helper
// ─────────────────────────────────────────────

func handleFlood(err error) bool {
	if err == nil {
		return false
	}
	errText := strings.ToUpper(err.Error())
	if strings.Contains(errText, "FLOOD_WAIT") {
		var wait int
		fmt.Sscanf(errText, "FLOOD_WAIT_%d", &wait)
		if wait <= 0 {
			wait = 5
		}
		time.Sleep(time.Duration(wait+1) * time.Second)
		return true
	}
	return false
}

// ─────────────────────────────────────────────
// Generic Normal Raid
// FIX: GetReplyMessage() called BEFORE m.Delete()
// Previously it was called after Delete which caused it to return nil,
// so all messages sent without reply even when command was used as a reply
// ─────────────────────────────────────────────

func genericRaid(m *telegram.NewMessage, raidList []string, raidType string, delay time.Duration) error {
	args := GetArgs(m)
	var count int
	fmt.Sscanf(args, "%d", &count)
	if count < 1 {
		Reply(m, "❌ ᴄᴏᴜɴᴛ ᴍᴜsᴛ ʙᴇ ᴀᴛ ʟᴇᴀsᴛ 1.")
		return nil
	}

	// Fetch reply ID BEFORE m.Delete() — after delete, reply context is lost
	var replyToID int32
	if m.IsReply() {
		if replyMsg, err := m.GetReplyMessage(); err == nil && replyMsg != nil {
			replyToID = int32(replyMsg.ID)
		}
	}

	database.SaveRaid(database.RaidSession{
		ChatID:    m.ChatID(),
		ReplyToID: replyToID,
		RaidType:  raidType,
		Count:     count,
	})

	_, _ = m.Delete()
	setRaidActive(m.ChatID())
	defer func() {
		setRaidStopped(m.ChatID())
		database.DeleteRaid(m.ChatID())
	}()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < count; i++ {
		if !isRaidActive(m.ChatID()) {
			break
		}
		text := raidList[rng.Intn(len(raidList))]
		opts := &telegram.SendOptions{ParseMode: telegram.HTML}
		if replyToID > 0 {
			opts.ReplyID = replyToID
		}
		_, err := m.Client.SendMessage(m.ChatID(), text, opts)
		if err != nil {
			if handleFlood(err) {
				i--
				continue
			}
			time.Sleep(2 * time.Second)
			continue
		}
		time.Sleep(delay)
	}

	return nil
}

// ─────────────────────────────────────────────
// ReplyRaid Watcher
// Waits for target user to send a message,
// sends one random raid reply, then stops itself
// ─────────────────────────────────────────────

func startReplyRaidWatcher(m *telegram.NewMessage, raidList []string, raidType string) error {
	if !m.IsReply() {
		Reply(m, "↩️ ʀᴇᴘʟʏ ᴛᴏ ᴛʜᴇ ᴛᴀʀɢᴇᴛ ᴜsᴇʀ's ᴍᴇssᴀɢᴇ ꜰɪʀsᴛ.")
		return nil
	}

	// Fetch replied message BEFORE m.Delete() to get both ID and senderID
	replyMsg, err := m.GetReplyMessage()
	if err != nil || replyMsg == nil {
		Reply(m, "❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ɢᴇᴛ ʀᴇᴘʟɪᴇᴅ ᴍᴇssᴀɢᴇ.")
		return nil
	}

	replyToID := int32(replyMsg.ID)
	targetUserID := replyMsg.SenderID()
	chatID := m.ChatID()

	replyRaidWatchersMu.Lock()
	if cancel, exists := replyRaidWatchers[chatID]; exists {
		cancel()
		delete(replyRaidWatchers, chatID)
	}
	stopped := false
	replyRaidWatchers[chatID] = func() { stopped = true }
	_ = stopped
	replyRaidWatchersMu.Unlock()

	database.SaveRaid(database.RaidSession{
		ChatID:       chatID,
		ReplyToID:    replyToID,
		RaidType:     raidType,
		Count:        -1,
		TargetUserID: targetUserID,
	})

	_, _ = m.Delete()
	Reply(m, fmt.Sprintf(
		"👁 <b>ʀᴇᴘʟʏʀᴀɪᴅ ᴀᴄᴛɪᴠᴇ</b>\n» ᴛᴀʀɢᴇᴛ: <code>%d</code>\n» ᴡᴀɪᴛɪɴɢ ꜰᴏʀ ɴᴇxᴛ ᴍᴇssᴀɢᴇ...",
		targetUserID,
	))
	return nil
}

// ─────────────────────────────────────────────
// TriggerReplyRaidIfActive
// Called by module.go on every incoming message.
// Fires once when target user sends a message, then stops.
// ─────────────────────────────────────────────

func TriggerReplyRaidIfActive(m *telegram.NewMessage) {
	chatID := m.ChatID()
	senderID := m.SenderID()

	replyRaidWatchersMu.Lock()
	_, watcherActive := replyRaidWatchers[chatID]
	replyRaidWatchersMu.Unlock()

	if !watcherActive {
		return
	}

	session := database.GetRaid(chatID)
	if session == nil || session.Count != -1 {
		replyRaidWatchersMu.Lock()
		delete(replyRaidWatchers, chatID)
		replyRaidWatchersMu.Unlock()
		return
	}

	if session.TargetUserID == 0 || senderID != session.TargetUserID {
		return
	}

	// Stop watcher first to prevent duplicate triggers
	replyRaidWatchersMu.Lock()
	if cancel, ok := replyRaidWatchers[chatID]; ok {
		cancel()
		delete(replyRaidWatchers, chatID)
	}
	replyRaidWatchersMu.Unlock()

	database.DeleteRaid(chatID)

	raidList := getRaidList(session.RaidType)
	if len(raidList) == 0 {
		return
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	text := raidList[rng.Intn(len(raidList))]

	_, err := m.Client.SendMessage(chatID, text, &telegram.SendOptions{
		ParseMode: telegram.HTML,
		ReplyID:   int32(m.ID),
	})
	if err != nil {
		handleFlood(err)
		_, _ = m.Client.SendMessage(chatID, text, &telegram.SendOptions{
			ParseMode: telegram.HTML,
			ReplyID:   int32(m.ID),
		})
	}
}

// ─────────────────────────────────────────────
// ResumeRaids — called on bot startup
// Normal raids resume; watcher sessions silently restore
// ─────────────────────────────────────────────

func ResumeRaids(cl *telegram.Client) {
	sessions := database.LoadActiveRaids()
	for _, s := range sessions {
		go func(sess database.RaidSession) {
			raidList := getRaidList(sess.RaidType)
			if len(raidList) == 0 {
				return
			}

			if sess.Count == -1 {
				replyRaidWatchersMu.Lock()
				replyRaidWatchers[sess.ChatID] = func() {}
				replyRaidWatchersMu.Unlock()
				return
			}

			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			activeRaidChatsMu.Lock()
			activeRaidChats[sess.ChatID] = true
			activeRaidChatsMu.Unlock()
			defer func() {
				setRaidStopped(sess.ChatID)
				database.DeleteRaid(sess.ChatID)
			}()

			for i := 0; i < sess.Count; i++ {
				if !isRaidActive(sess.ChatID) {
					break
				}
				text := raidList[rng.Intn(len(raidList))]
				opts := &telegram.SendOptions{ParseMode: telegram.HTML}
				if sess.ReplyToID > 0 {
					opts.ReplyID = sess.ReplyToID
				}
				_, err := cl.SendMessage(sess.ChatID, text, opts)
				if err != nil {
					if handleFlood(err) {
						i--
						continue
					}
					time.Sleep(2 * time.Second)
					continue
				}
				time.Sleep(800 * time.Millisecond)
			}
		}(s)
	}
}

func getRaidList(raidType string) []string {
	switch raidType {
	case "raid":
		return config.RAID
	case "hraid":
		return config.HRAID
	case "eraid":
		return config.ERAID
	case "punraid":
		return config.PUNRAID
	case "praid":
		return config.PORNTEXT
	default:
		return config.RAID
	}
}

// ─────────────────────────────────────────────
// REGISTER
// ─────────────────────────────────────────────

func init() {
	Register(ModuleInfo{
		Name:        "Raid",
		Description: "Multi Language Raid Commands",
		Commands: []CommandInfo{
			{Pattern: "raid",       Handler: raidHandler,       Sudo: true},
			{Pattern: "hraid",      Handler: hraidHandler,      Sudo: true},
			{Pattern: "eraid",      Handler: eraidHandler,      Sudo: true},
			{Pattern: "punraid",    Handler: punraidHandler,    Sudo: true},
			{Pattern: "replyraid",  Handler: replyRaidHandler,  Sudo: true},
			{Pattern: "hreplyraid", Handler: hreplyRaidHandler, Sudo: true},
			{Pattern: "ereplyraid", Handler: ereplyRaidHandler, Sudo: true},
			{Pattern: "preplyraid", Handler: preplyRaidHandler, Sudo: true},
			{Pattern: "stopraid",   Handler: stopRaidHandler,   Sudo: true},
		},
	})
}
