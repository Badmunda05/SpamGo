package modules

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

/*
	PbxGo Powerful Admin Tools
	Created By: BadMunda
*/

// ─────────────────────────────────────────────
// Active Power Action Tracker
// ─────────────────────────────────────────────

var (
	activePowerActions = make(map[int64]bool)
	activePowerMu      sync.Mutex
)

func setPowerActive(chatID int64) {
	activePowerMu.Lock()
	activePowerActions[chatID] = true
	activePowerMu.Unlock()
}

func setPowerStopped(chatID int64) {
	activePowerMu.Lock()
	delete(activePowerActions, chatID)
	activePowerMu.Unlock()
}

func isPowerActive(chatID int64) bool {
	activePowerMu.Lock()
	defer activePowerMu.Unlock()
	return activePowerActions[chatID]
}

// ─────────────────────────────────────────────
// Flood Wait Helper
// ─────────────────────────────────────────────

func handlePowerFlood(err error) bool {
	if err == nil {
		return false
	}
	errText := strings.ToUpper(err.Error())
	if strings.Contains(errText, "FLOOD_WAIT") {
		var wait int
		fmt.Sscanf(errText, "FLOOD_WAIT_%d", &wait)
		if wait < 1 {
			wait = 5
		}
		time.Sleep(time.Duration(wait+2) * time.Second)
		return true
	}
	return false
}

// ─────────────────────────────────────────────
// Helper: Resolve Target Chat ID
// ─────────────────────────────────────────────

func getTargetChatID(m *telegram.NewMessage) (int64, error) {
	args := strings.TrimSpace(GetArgs(m))
	if args == "" {
		return m.ChatID(), nil
	}
	if id, err := strconv.ParseInt(args, 10, 64); err == nil {
		return id, nil
	}
	if strings.HasPrefix(args, "@") {
		peer, err := m.Client.ResolvePeer(strings.TrimPrefix(args, "@"))
		if err != nil {
			return 0, fmt.Errorf("failed to resolve: %v", err)
		}
		switch p := peer.(type) {
		case *telegram.InputPeerChannel:
			return p.ChannelID, nil
		case *telegram.InputPeerChat:
			return p.ChatID, nil
		case *telegram.InputPeerUser:
			return p.UserID, nil
		}
	}
	return 0, fmt.Errorf("invalid group username or ID")
}

// ─────────────────────────────────────────────
// .banall
// ─────────────────────────────────────────────

func banAllHandler(m *telegram.NewMessage) error {
	targetID, err := getTargetChatID(m)
	if err != nil {
		Reply(m, fmt.Sprintf("❌ %s", err.Error()))
		return nil
	}

	_, _ = m.Delete()
	setPowerActive(targetID)
	defer setPowerStopped(targetID)

	Reply(m, fmt.Sprintf("🔥 <b>ʙᴀɴᴀʟʟ sᴛᴀʀᴛᴇᴅ</b>\n» ᴄʜᴀᴛ: <code>%d</code>", targetID))

	participants, err := m.Client.GetParticipants(targetID, nil)
	if err != nil {
		Reply(m, "❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ɢᴇᴛ ᴘᴀʀᴛɪᴄɪᴘᴀɴᴛs.")
		return nil
	}

	count := 0
	for _, user := range participants.Users {
		if !isPowerActive(targetID) {
			break
		}
		if user.Self || user.Bot {
			continue
		}
		_, err := m.Client.EditBannedParticipant(targetID, user.ID, &telegram.ChatBannedRights{
			ViewMessages: true,
			UntilDate:    0,
		})
		if err != nil {
			if handlePowerFlood(err) {
				continue
			}
			time.Sleep(2 * time.Second)
			continue
		}
		count++
		if count%10 == 0 {
			Reply(m, fmt.Sprintf("✅ <b>ʙᴀɴɴᴇᴅ:</b> <code>%d</code>", count))
		}
		time.Sleep(800 * time.Millisecond)
	}
	Reply(m, fmt.Sprintf("🏁 <b>ʙᴀɴᴀʟʟ ᴅᴏɴᴇ!</b>\n» ᴛᴏᴛᴀʟ: <code>%d</code>", count))
	return nil
}

// ─────────────────────────────────────────────
// .unbanall
// ─────────────────────────────────────────────

func unbanAllHandler(m *telegram.NewMessage) error {
	targetID, err := getTargetChatID(m)
	if err != nil {
		Reply(m, fmt.Sprintf("❌ %s", err.Error()))
		return nil
	}

	_, _ = m.Delete()
	setPowerActive(targetID)
	defer setPowerStopped(targetID)

	Reply(m, fmt.Sprintf("🔄 <b>ᴜɴʙᴀɴᴀʟʟ sᴛᴀʀᴛᴇᴅ</b>\n» ᴄʜᴀᴛ: <code>%d</code>", targetID))

	// Get banned list using channel participants banned filter
	participants, err := m.Client.GetParticipants(targetID, &telegram.ParticipantsOptions{
		Filter: &telegram.ChannelParticipantsBanned{},
	})
	if err != nil {
		Reply(m, "❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ɢᴇᴛ ʙᴀɴɴᴇᴅ ʟɪsᴛ.")
		return nil
	}

	count := 0
	for _, user := range participants.Users {
		if !isPowerActive(targetID) {
			break
		}
		_, err := m.Client.EditBannedParticipant(targetID, user.ID, &telegram.ChatBannedRights{})
		if err != nil {
			if handlePowerFlood(err) {
				continue
			}
			time.Sleep(2 * time.Second)
			continue
		}
		count++
		time.Sleep(700 * time.Millisecond)
	}
	Reply(m, fmt.Sprintf("🏁 <b>ᴜɴʙᴀɴᴀʟʟ ᴅᴏɴᴇ!</b>\n» ᴛᴏᴛᴀʟ: <code>%d</code>", count))
	return nil
}

// ─────────────────────────────────────────────
// .muteall
// ─────────────────────────────────────────────

func muteAllHandler(m *telegram.NewMessage) error {
	targetID, err := getTargetChatID(m)
	if err != nil {
		Reply(m, fmt.Sprintf("❌ %s", err.Error()))
		return nil
	}

	_, _ = m.Delete()
	setPowerActive(targetID)
	defer setPowerStopped(targetID)

	Reply(m, fmt.Sprintf("🔇 <b>ᴍᴜᴛᴇᴀʟʟ sᴛᴀʀᴛᴇᴅ</b>\n» ᴄʜᴀᴛ: <code>%d</code>", targetID))

	participants, err := m.Client.GetParticipants(targetID, nil)
	if err != nil {
		Reply(m, "❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ɢᴇᴛ ᴘᴀʀᴛɪᴄɪᴘᴀɴᴛs.")
		return nil
	}

	count := 0
	for _, user := range participants.Users {
		if !isPowerActive(targetID) {
			break
		}
		if user.Self || user.Bot {
			continue
		}
		_, err := m.Client.EditBannedParticipant(targetID, user.ID, &telegram.ChatBannedRights{
			SendMessages: true,
			UntilDate:    0,
		})
		if err != nil {
			if handlePowerFlood(err) {
				continue
			}
			time.Sleep(2 * time.Second)
			continue
		}
		count++
		if count%15 == 0 {
			Reply(m, fmt.Sprintf("✅ <b>ᴍᴜᴛᴇᴅ:</b> <code>%d</code>", count))
		}
		time.Sleep(900 * time.Millisecond)
	}
	Reply(m, fmt.Sprintf("🏁 <b>ᴍᴜᴛᴇᴀʟʟ ᴅᴏɴᴇ!</b>\n» ᴛᴏᴛᴀʟ: <code>%d</code>", count))
	return nil
}

// ─────────────────────────────────────────────
// .unmuteall
// ─────────────────────────────────────────────

func unmuteAllHandler(m *telegram.NewMessage) error {
	targetID, err := getTargetChatID(m)
	if err != nil {
		Reply(m, fmt.Sprintf("❌ %s", err.Error()))
		return nil
	}

	_, _ = m.Delete()
	setPowerActive(targetID)
	defer setPowerStopped(targetID)

	Reply(m, fmt.Sprintf("🔊 <b>ᴜɴᴍᴜᴛᴇᴀʟʟ sᴛᴀʀᴛᴇᴅ</b>\n» ᴄʜᴀᴛ: <code>%d</code>", targetID))

	participants, err := m.Client.GetParticipants(targetID, nil)
	if err != nil {
		Reply(m, "❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ɢᴇᴛ ᴘᴀʀᴛɪᴄɪᴘᴀɴᴛs.")
		return nil
	}

	count := 0
	for _, user := range participants.Users {
		if !isPowerActive(targetID) {
			break
		}
		_, err := m.Client.EditBannedParticipant(targetID, user.ID, &telegram.ChatBannedRights{})
		if err != nil {
			if handlePowerFlood(err) {
				continue
			}
			time.Sleep(2 * time.Second)
			continue
		}
		count++
		time.Sleep(800 * time.Millisecond)
	}
	Reply(m, fmt.Sprintf("🏁 <b>ᴜɴᴍᴜᴛᴇᴀʟʟ ᴅᴏɴᴇ!</b>\n» ᴛᴏᴛᴀʟ: <code>%d</code>", count))
	return nil
}

// ─────────────────────────────────────────────
// .kickall
// ─────────────────────────────────────────────

func kickAllHandler(m *telegram.NewMessage) error {
	targetID, err := getTargetChatID(m)
	if err != nil {
		Reply(m, fmt.Sprintf("❌ %s", err.Error()))
		return nil
	}

	_, _ = m.Delete()
	setPowerActive(targetID)
	defer setPowerStopped(targetID)

	Reply(m, fmt.Sprintf("👢 <b>ᴋɪᴄᴋᴀʟʟ sᴛᴀʀᴛᴇᴅ</b>\n» ᴄʜᴀᴛ: <code>%d</code>", targetID))

	participants, err := m.Client.GetParticipants(targetID, nil)
	if err != nil {
		Reply(m, "❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ɢᴇᴛ ᴘᴀʀᴛɪᴄɪᴘᴀɴᴛs.")
		return nil
	}

	count := 0
	for _, user := range participants.Users {
		if !isPowerActive(targetID) {
			break
		}
		if user.Self || user.Bot {
			continue
		}
		err := m.Client.KickParticipant(targetID, user.ID)
		if err != nil {
			if handlePowerFlood(err) {
				continue
			}
			time.Sleep(2 * time.Second)
			continue
		}
		count++
		if count%8 == 0 {
			Reply(m, fmt.Sprintf("✅ <b>ᴋɪᴄᴋᴇᴅ:</b> <code>%d</code>", count))
		}
		time.Sleep(1000 * time.Millisecond)
	}
	Reply(m, fmt.Sprintf("🏁 <b>ᴋɪᴄᴋᴀʟʟ ᴅᴏɴᴇ!</b>\n» ᴛᴏᴛᴀʟ: <code>%d</code>", count))
	return nil
}

// ─────────────────────────────────────────────
// .stoppower
// ─────────────────────────────────────────────

func stopPowerHandler(m *telegram.NewMessage) error {
	setPowerStopped(m.ChatID())
	_, _ = m.Delete()
	Reply(m, "🛑 <b>ᴘᴏᴡᴇʀ ᴀᴄᴛɪᴏɴ sᴛᴏᴘᴘᴇᴅ!</b>")
	return nil
}

// ─────────────────────────────────────────────
// REGISTER
// ─────────────────────────────────────────────

func init() {
	Register(ModuleInfo{
		Name:        "Powerful",
		Description: "BanAll KickAll MuteAll with Cross-Group Support",
		Commands: []CommandInfo{
			{Pattern: "banall",    Handler: banAllHandler,    Sudo: true},
			{Pattern: "unbanall",  Handler: unbanAllHandler,  Sudo: true},
			{Pattern: "muteall",   Handler: muteAllHandler,   Sudo: true},
			{Pattern: "unmuteall", Handler: unmuteAllHandler, Sudo: true},
			{Pattern: "kickall",   Handler: kickAllHandler,   Sudo: true},
			{Pattern: "stoppower", Handler: stopPowerHandler, Sudo: true},
		},
	})
}
