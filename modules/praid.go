package modules

import (
	"fmt"
	"math/rand"
	"time"

	"pbxgo/config"

	"github.com/amarnathcjd/gogram/telegram"
)

/*
	PbxGo Porn Module — Text + Video
	Created By: BadMunda
*/

var PORNTEXT = config.PORNTEXT

// ─────────────────────────────────────────────
// .pspam — porn text spam
// ─────────────────────────────────────────────

func pspamHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	var count int
	extraText := ""

	fmt.Sscanf(args, "%d", &count)
	if count < 1 {
		count = 5
	}

	if idx := findSpace(args); idx != -1 {
		extraText = args[idx+1:]
	}

	// Fetch reply ID before Delete
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
		msg := PORNTEXT[rng.Intn(len(PORNTEXT))]
		if extraText != "" {
			msg += "\n" + extraText
		}
		opts := &telegram.SendOptions{ParseMode: telegram.HTML}
		if replyToID > 0 {
			opts.ReplyID = replyToID
		}
		_, err := m.Client.SendMessage(m.ChatID(), msg, opts)
		if err != nil {
			handleFlood(err)
		}
		time.Sleep(700 * time.Millisecond)
	}
	return nil
}

// ─────────────────────────────────────────────
// .pornv / .pspamv — porn video spam
// ─────────────────────────────────────────────

func pornVideoSpamHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	var count int
	fmt.Sscanf(args, "%d", &count)
	if count < 1 {
		count = 5
	}

	// Fetch reply ID before Delete
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
		opts := &telegram.MediaOptions{}
		if replyToID > 0 {
			opts.ReplyID = replyToID
		}
		_, err := m.Client.SendMedia(m.ChatID(), videoURL, opts)
		if err != nil {
			handleFlood(err)
		}
		time.Sleep(1500 * time.Millisecond)
	}
	return nil
}

// ─────────────────────────────────────────────
// .praid — porn text raid
// ─────────────────────────────────────────────

func praidHandler(m *telegram.NewMessage) error {
	return genericRaid(m, config.PORNTEXT, "praid", 600*time.Millisecond)
}

// ─────────────────────────────────────────────
// .prraid — porn reply raid (watcher mode)
// ─────────────────────────────────────────────

func prreplyRaidHandler(m *telegram.NewMessage) error {
	return startReplyRaidWatcher(m, config.PORNTEXT, "praid")
}

// ─────────────────────────────────────────────
// REGISTER
// ─────────────────────────────────────────────

func init() {
	Register(ModuleInfo{
		Name:        "Porn",
		Description: "Porn Text + Video Spam & Raid",
		Commands: []CommandInfo{
			{Pattern: "pspam",  Handler: pspamHandler,         Sudo: true},
			{Pattern: "pornv",  Handler: pornVideoSpamHandler, Sudo: true},
			{Pattern: "pspamv", Handler: pornVideoSpamHandler, Sudo: true},
			{Pattern: "praid",  Handler: praidHandler,         Sudo: true},
			{Pattern: "prraid", Handler: prreplyRaidHandler,   Sudo: true},
		},
	})
}
