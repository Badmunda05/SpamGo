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
// Porn Reply Raid Watcher — separate from normal raid
// ─────────────────────────────────────────────

type pornWatcherSession struct {
	targetUserID int64
}

var (
	pornWatchers   = make(map[int64]pornWatcherSession)
	pornWatchersMu sync.Mutex
)

// ─────────────────────────────────────────────
// .pspam — porn text spam only
// ─────────────────────────────────────────────

func pspamHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	var count int
	fmt.Sscanf(args, "%d", &count)
	if count < 1 {
		count = 5
	}

	// Fetch reply ID BEFORE Delete
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
		msg := config.PORNTEXT[rng.Intn(len(config.PORNTEXT))]
		opts := &telegram.SendOptions{ParseMode: telegram.HTML}
		if replyToID > 0 {
			opts.ReplyID = replyToID
		}
		_, err := m.Client.SendMessage(m.ChatID(), msg, opts)
		if err != nil {
			if handleFlood(err) {
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
// .praid — har iteration: 1 text + 1 video dono
// ─────────────────────────────────────────────

func praidHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	var count int
	fmt.Sscanf(args, "%d", &count)
	if count < 1 {
		Reply(m, "❌ ᴄᴏᴜɴᴛ ᴍᴜsᴛ ʙᴇ ᴀᴛ ʟᴇᴀsᴛ 1.")
		return nil
	}

	// Fetch reply ID BEFORE Delete
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

		// Step 1: text bhejo
		textMsg := config.PORNTEXT[rng.Intn(len(config.PORNTEXT))]
		textOpts := &telegram.SendOptions{ParseMode: telegram.HTML}
		if replyToID > 0 {
			textOpts.ReplyID = replyToID
		}
		_, err := m.Client.SendMessage(m.ChatID(), textMsg, textOpts)
		if err != nil {
			if handleFlood(err) {
				i--
				continue
			}
			time.Sleep(2 * time.Second)
			continue
		}
		time.Sleep(400 * time.Millisecond)

		if !isRaidActive(m.ChatID()) {
			break
		}

		// Step 2: video bhejo
		videoURL := config.PORNVIDEOS[rng.Intn(len(config.PORNVIDEOS))]
		videoOpts := &telegram.MediaOptions{}
		if replyToID > 0 {
			videoOpts.ReplyID = replyToID
		}
		_, err = m.Client.SendMedia(m.ChatID(), videoURL, videoOpts)
		if err != nil {
			if handleFlood(err) {
				i--
				continue
			}
			time.Sleep(3 * time.Second)
			continue
		}
		time.Sleep(800 * time.Millisecond)
	}
	return nil
}

// ─────────────────────────────────────────────
// .preplyraid — watcher mode
// Target user jdo msg kare: 1 text + 1 video reply vich, fir band
// ─────────────────────────────────────────────

func preplyRaidHandler(m *telegram.NewMessage) error {
	if !m.IsReply() {
		Reply(m, "↩️ ʀᴇᴘʟʏ ᴛᴏ ᴛʜᴇ ᴛᴀʀɢᴇᴛ ᴜsᴇʀ's ᴍᴇssᴀɢᴇ ꜰɪʀsᴛ.")
		return nil
	}

	// Fetch BEFORE Delete
	replyMsg, err := m.GetReplyMessage()
	if err != nil || replyMsg == nil {
		Reply(m, "❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ɢᴇᴛ ʀᴇᴘʟɪᴇᴅ ᴍᴇssᴀɢᴇ.")
		return nil
	}

	targetUserID := replyMsg.SenderID()
	chatID := m.ChatID()

	// Register in porn watcher map
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

// TriggerPReplyRaidIfActive — module.go har msg te call karda hai
// Jdo target user msg kare: 1 text + 1 video reply, fir watcher band
func TriggerPReplyRaidIfActive(m *telegram.NewMessage) {
	chatID := m.ChatID()
	senderID := m.SenderID()

	pornWatchersMu.Lock()
	session, active := pornWatchers[chatID]
	pornWatchersMu.Unlock()

	if !active || senderID != session.targetUserID {
		return
	}

	// Stop watcher — fire only once
	pornWatchersMu.Lock()
	delete(pornWatchers, chatID)
	pornWatchersMu.Unlock()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 1 text reply
	textMsg := config.PORNTEXT[rng.Intn(len(config.PORNTEXT))]
	_, _ = m.Client.SendMessage(chatID, textMsg, &telegram.SendOptions{
		ParseMode: telegram.HTML,
		ReplyID:   int32(m.ID),
	})
	time.Sleep(400 * time.Millisecond)

	// 1 video reply
	videoURL := config.PORNVIDEOS[rng.Intn(len(config.PORNVIDEOS))]
	_, _ = m.Client.SendMedia(chatID, videoURL, &telegram.MediaOptions{
		ReplyID: int32(m.ID),
	})
}

// ─────────────────────────────────────────────
// REGISTER
// ─────────────────────────────────────────────

func init() {
	Register(ModuleInfo{
		Name:        "Porn",
		Description: "Porn Text Spam + Text+Video Raid",
		Commands: []CommandInfo{
			{Pattern: "pspam",      Handler: pspamHandler,     Sudo: true},
			{Pattern: "praid",      Handler: praidHandler,      Sudo: true},
			{Pattern: "preplyraid", Handler: preplyRaidHandler, Sudo: true},
		},
	})
}
