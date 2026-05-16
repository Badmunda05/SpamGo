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
	PbxGo Powerful Admin Tools + Cross Group Support
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
// Generic Flood Handler
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

func getTargetChat(m *telegram.NewMessage) (int64, error) {
	args := strings.TrimSpace(GetArgs(m))

	if args == "" {
		return m.ChatID(), nil // Current chat
	}

	// If numeric ID ( -100123456789 )
	if id, err := strconv.ParseInt(args, 10, 64); err == nil {
		return id, nil
	}

	// If username (@groupusername)
	if strings.HasPrefix(args, "@") {
		chat, err := m.Client.ResolveUsername(strings.TrimPrefix(args, "@"))
		if err != nil {
			return 0, err
		}
		return chat.ID(), nil
	}

	return 0, fmt.Errorf("invalid group username or ID")
}

// ─────────────────────────────────────────────
// .banall
// ─────────────────────────────────────────────

func banAllHandler(m *telegram.NewMessage) error {
	targetID, err := getTargetChat(m)
	if err != nil {
		Reply(m, "❌ Invalid group username or ID")
		return nil
	}

	_, _ = m.Delete()
	setPowerActive(targetID) // Important: Use targetID for stop
	defer setPowerStopped(targetID)

	Reply(m, fmt.Sprintf("🔥 <b>ʙᴀɴᴀʟʟ sᴛᴀʀᴛᴇᴅ ɪɴ:</b> <code>%d</code>", targetID))

	participants, err := m.Client.GetParticipants(targetID, nil)
	if err != nil {
		Reply(m, "❌ Failed to get participants. Bot must be admin there.")
		return nil
	}

	count := 0
	for _, user := range participants.Users {
		if !isPowerActive(targetID) {
			break
		}

		if user.ID == m.Client.Me().ID || user.IsAdmin() || user.IsCreator() {
			continue
		}

		err := m.Client.BanMember(targetID, user.ID)
		if err != nil {
			if handlePowerFlood(err) {
				continue
			}
			time.Sleep(2 * time.Second)
			continue
		}

		count++
		if count%10 == 0 {
			Reply(m, fmt.Sprintf("✅ <b>Banned:</b> <code>%d</code> users", count))
		}
		time.Sleep(800 * time.Millisecond)
	}

	Reply(m, fmt.Sprintf("🏁 <b>ʙᴀɴᴀʟʟ ᴄᴏᴍᴘʟᴇᴛᴇᴅ!</b>\nGroup: <code>%d</code>\nTotal Banned: <code>%d</code>", targetID, count))
	return nil
}

// Similar changes applied to all other commands...

func unbanAllHandler(m *telegram.NewMessage) error {
	targetID, err := getTargetChat(m)
	if err != nil {
		Reply(m, "❌ Invalid group username or ID")
		return nil
	}

	_, _ = m.Delete()
	setPowerActive(targetID)
	defer setPowerStopped(targetID)

	Reply(m, fmt.Sprintf("🔄 <b>ᴜɴʙᴀɴᴀʟʟ sᴛᴀʀᴛᴇᴅ ɪɴ:</b> <code>%d</code>", targetID))

	participants, err := m.Client.GetBanned(targetID)
	if err != nil {
		Reply(m, "❌ Failed to get banned list.")
		return nil
	}

	count := 0
	for _, user := range participants {
		if !isPowerActive(targetID) {
			break
		}

		err := m.Client.UnbanMember(targetID, user.ID)
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

	Reply(m, fmt.Sprintf("🏁 <b>ᴜɴʙᴀɴᴀʟʟ ᴄᴏᴍᴘʟᴇᴛᴇᴅ!</b>\nTotal Unbanned: <code>%d</code>", count))
	return nil
}

func muteAllHandler(m *telegram.NewMessage) error {
	targetID, err := getTargetChat(m)
	if err != nil {
		Reply(m, "❌ Invalid group username or ID")
		return nil
	}

	_, _ = m.Delete()
	setPowerActive(targetID)
	defer setPowerStopped(targetID)

	Reply(m, fmt.Sprintf("🔇 <b>ᴍᴜᴛᴇᴀʟʟ sᴛᴀʀᴛᴇᴅ ɪɴ:</b> <code>%d</code>", targetID))

	participants, err := m.Client.GetParticipants(targetID, nil)
	if err != nil {
		Reply(m, "❌ Failed to get participants.")
		return nil
	}

	count := 0
	for _, user := range participants.Users {
		if !isPowerActive(targetID) {
			break
		}

		if user.ID == m.Client.Me().ID || user.IsAdmin() || user.IsCreator() {
			continue
		}

		err := m.Client.RestrictMember(targetID, user.ID, &telegram.ChatPermissions{CanSendMessages: false})
		if err != nil {
			if handlePowerFlood(err) {
				continue
			}
			time.Sleep(2 * time.Second)
			continue
		}

		count++
		if count%15 == 0 {
			Reply(m, fmt.Sprintf("✅ <b>Muted:</b> <code>%d</code>", count))
		}
		time.Sleep(900 * time.Millisecond)
	}

	Reply(m, fmt.Sprintf("🏁 <b>ᴍᴜᴛᴇᴀʟʟ ᴄᴏᴍᴘʟᴇᴛᴇᴅ!</b>\nTotal Muted: <code>%d</code>", count))
	return nil
}

func unmuteAllHandler(m *telegram.NewMessage) error {
	targetID, err := getTargetChat(m)
	if err != nil {
		Reply(m, "❌ Invalid group username or ID")
		return nil
	}

	_, _ = m.Delete()
	setPowerActive(targetID)
	defer setPowerStopped(targetID)

	Reply(m, fmt.Sprintf("🔊 <b>ᴜɴᴍᴜᴛᴇᴀʟʟ sᴛᴀʀᴛᴇᴅ ɪɴ:</b> <code>%d</code>", targetID))

	participants, err := m.Client.GetParticipants(targetID, nil)
	if err != nil {
		Reply(m, "❌ Failed to get participants.")
		return nil
	}

	count := 0
	for _, user := range participants.Users {
		if !isPowerActive(targetID) {
			break
		}

		err := m.Client.RestrictMember(targetID, user.ID, &telegram.ChatPermissions{
			CanSendMessages: true,
			CanSendMedia:    true,
			CanSendStickers: true,
		})
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

	Reply(m, fmt.Sprintf("🏁 <b>ᴜɴᴍᴜᴛᴇᴀʟʟ ᴄᴏᴍᴘʟᴇᴛᴇᴅ!</b>\nTotal Unmuted: <code>%d</code>", count))
	return nil
}

func kickAllHandler(m *telegram.NewMessage) error {
	targetID, err := getTargetChat(m)
	if err != nil {
		Reply(m, "❌ Invalid group username or ID")
		return nil
	}

	_, _ = m.Delete()
	setPowerActive(targetID)
	defer setPowerStopped(targetID)

	Reply(m, fmt.Sprintf("👢 <b>ᴋɪᴄᴋᴀʟʟ sᴛᴀʀᴛᴇᴅ ɪɴ:</b> <code>%d</code>", targetID))

	participants, err := m.Client.GetParticipants(targetID, nil)
	if err != nil {
		Reply(m, "❌ Failed to get participants.")
		return nil
	}

	count := 0
	for _, user := range participants.Users {
		if !isPowerActive(targetID) {
			break
		}

		if user.ID == m.Client.Me().ID || user.IsAdmin() || user.IsCreator() {
			continue
		}

		err := m.Client.KickMember(targetID, user.ID)
		if err != nil {
			if handlePowerFlood(err) {
				continue
			}
			time.Sleep(2 * time.Second)
			continue
		}

		count++
		if count%8 == 0 {
			Reply(m, fmt.Sprintf("✅ <b>Kicked:</b> <code>%d</code>", count))
		}
		time.Sleep(1000 * time.Millisecond)
	}

	Reply(m, fmt.Sprintf("🏁 <b>ᴋɪᴄᴋᴀʟʟ ᴄᴏᴍᴘʟᴇᴛᴇᴅ!</b>\nTotal Kicked: <code>%d</code>", count))
	return nil
}

// ─────────────────────────────────────────────
// .stoppower
// ─────────────────────────────────────────────

func stopPowerHandler(m *telegram.NewMessage) error {
	// Stop both current and any active
	setPowerStopped(m.ChatID())
	Reply(m, "🛑 <b>ᴘᴏᴡᴇʀ ᴀᴄᴛɪᴏɴ sᴛᴏᴘᴘᴇᴅ!</b>")
	return nil
}

// ─────────────────────────────────────────────
// REGISTER
// ─────────────────────────────────────────────

func init() {
	Register(ModuleInfo{
		Name:        "Powerful",
		Description: "BanAll, KickAll, MuteAll with Cross-Group Support",
		Commands: []CommandInfo{
			{Pattern: "banall", Handler: banAllHandler, Sudo: true},
			{Pattern: "unbanall", Handler: unbanAllHandler, Sudo: true},
			{Pattern: "muteall", Handler: muteAllHandler, Sudo: true},
			{Pattern: "unmuteall", Handler: unmuteAllHandler, Sudo: true},
			{Pattern: "kickall", Handler: kickAllHandler, Sudo: true},
			{Pattern: "stoppower", Handler: stopPowerHandler, Sudo: true},
		},
	})
}
