package modules

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"pbxgo/config"

	"github.com/amarnathcjd/gogram/telegram"
)

/*
PbxGo Porn Module
Created By: BadMunda
*/

// ─────────────────────────────────────────────
// Porn Reply Raid Watcher
// ─────────────────────────────────────────────

type pornWatcherSession struct {
	targetUserID int64
}

var (
	pornWatchers   = make(map[int64]pornWatcherSession)
	pornWatchersMu sync.Mutex
)

// ─────────────────────────────────────────────
// .pspam — ONLY Porn Video Spam (No Text)
// ─────────────────────────────────────────────

func pspamHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	var count int
	fmt.Sscanf(args, "%d", &count)
	if count < 1 {
		count = 5
	}

	// Fetch reply ID
	var replyToID int32
	if m.IsReply() {
		if replyMsg, err := m.GetReplyMessage(); err == nil && replyMsg != nil {
			replyToID = int32(replyMsg.ID)
		}
	}

	_, _ = m.Delete()
	setSpamActive(m.ChatID())
	defer setSpamStopped(m.ChatID())

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < count; i++ {
		if !isSpamActive(m.ChatID()) {
			break
		}

		videoURL := config.PORNVIDEOS[rng.Intn(len(config.PORNVIDEOS))]

		videoOpts := &telegram.MediaOptions{}
		if replyToID > 0 {
			videoOpts.ReplyID = replyToID
		}

		_, err := m.Client.SendMedia(m.ChatID(), videoURL, videoOpts)
		if err != nil {
			if handleFlood(err) {
				i--
				continue
			}
			time.Sleep(3 * time.Second)
			continue
		}

		time.Sleep(900 * time.Millisecond) // good delay for video spam
	}
	return nil
}

// ─────────────────────────────────────────────
// .praid — Video + Text (Caption) Combined
// ─────────────────────────────────────────────

func praidHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	var count int
	fmt.Sscanf(args, "%d", &count)
	if count < 1 {
		Reply(m, "❌ ᴄᴏᴜɴᴛ ᴍᴜsᴛ ʙᴇ ᴀᴛ ʟᴇᴀsᴛ 1.")
		return nil
	}

	// Fetch reply ID
	var replyToID int32
	if m.IsReply() {
		if replyMsg, err := m.GetReplyMessage(); err == nil && replyMsg != nil {
			replyToID = int32(replyMsg.ID)
		}
	}

	_, _ = m.Delete()
	setRaidActive(m.ChatID())
	defer setRaidStopped(m.ChatID())

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < count; i++ {
		if !isRaidActive(m.ChatID()) {
			break
		}

		text := config.PORNTEXT[rng.Intn(len(config.PORNTEXT))]
		videoURL := config.PORNVIDEOS[rng.Intn(len(config.PORNVIDEOS))]

		videoOpts := &telegram.MediaOptions{
			Caption:   text,
			ParseMode: telegram.HTML,
		}

		if replyToID > 0 {
			videoOpts.ReplyID = replyToID
		}

		_, err := m.Client.SendMedia(m.ChatID(), videoURL, videoOpts)
		if err != nil {
			if handleFlood(err) {
				i--
				continue
			}
			time.Sleep(3 * time.Second)
			continue
		}

		time.Sleep(1300 * time.Millisecond)
	}
	return nil
}

// ─────────────────────────────────────────────
// .preplyraid — Watcher (Video + Text Reply)
// ─────────────────────────────────────────────

func pornReplyRaidHandler(m *telegram.NewMessage) error {
	if !m.IsReply() {
		Reply(m, "↩️ ʀᴇᴘʟʏ ᴛᴏ ᴛʜᴇ ᴛᴀʀɢᴇᴛ ᴜsᴇʀ's ᴍᴇssᴀɢᴇ ꜰɪʀsᴛ.")
		return nil
	}

	replyMsg, err := m.GetReplyMessage()
	if err != nil || replyMsg == nil {
		Reply(m, "❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ɢᴇᴛ ʀᴇᴘʟɪᴇᴅ ᴍᴇssᴀɢᴇ.")
		return nil
	}

	targetUserID := replyMsg.SenderID()
	chatID := m.ChatID()

	pornWatchersMu.Lock()
	pornWatchers[chatID] = pornWatcherSession{targetUserID: targetUserID}
	pornWatchersMu.Unlock()

	_, _ = m.Delete()
	Reply(m, fmt.Sprintf(
		"👁 <b>ᴘʀᴇᴘʟʏʀᴀɪᴅ ᴀᴄᴛɪᴠᴇ</b>\n» ᴛᴀʀɢᴇᴛ: <code>%d</code>\n» ᴡᴀɪᴛɪɴɢ ꜰᴏʀ ɴᴇxᴛ ᴍᴇssᴀɢᴇ...",
		targetUserID,
	))
	return nil
}

// Trigger on target message
func TriggerPReplyRaidIfActive(m *telegram.NewMessage) {
	chatID := m.ChatID()
	senderID := m.SenderID()

	pornWatchersMu.Lock()
	session, active := pornWatchers[chatID]
	pornWatchersMu.Unlock()

	if !active || senderID != session.targetUserID {
		return
	}

	// Stop watcher
	pornWatchersMu.Lock()
	delete(pornWatchers, chatID)
	pornWatchersMu.Unlock()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	text := config.PORNTEXT[rng.Intn(len(config.PORNTEXT))]
	videoURL := config.PORNVIDEOS[rng.Intn(len(config.PORNVIDEOS))]

	// Video with caption
	_, _ = m.Client.SendMedia(chatID, videoURL, &telegram.MediaOptions{
		Caption:   text,
		ParseMode: telegram.HTML,
		ReplyID:   int32(m.ID),
	})
}

// ─────────────────────────────────────────────
// REGISTER
// ─────────────────────────────────────────────

func init() {
	Register(ModuleInfo{
		Name:        "Porn",
		Description: "Porn Video Spam + Text+Video Raid",
		Commands: []CommandInfo{
			{Pattern: "pspam", Handler: pspamHandler, Sudo: true},
			{Pattern: "praid", Handler: praidHandler, Sudo: true},
			{Pattern: "preplyraid", Handler: pornReplyRaidHandler, Sudo: true},
		},
	})
}
