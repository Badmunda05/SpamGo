// modules/basic.go

package modules

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"pbxgo/config"
	"pbxgo/database"

	"github.com/amarnathcjd/gogram/telegram"
)

/*
	PbxSpamGo
	Created By: BadMunda
*/

var START_TEXT = `
<b>ʜᴇʏ <a href="tg://user?id=%d">%s</a>,</b>

ɪ ᴀᴍ <b>ᴘʙx sᴘᴀᴍ ɢᴏ</b>
━━━━━━━━━━━━━━━━━━━

» <b>ᴅᴇᴠᴇʟᴏᴘᴇʀ :</b> <a href="https://t.me/badnundaxd">ʙᴀᴅ</a>

» <b>ᴠᴇʀsɪᴏɴ :</b> <code>2.0.0</code>
» <b>ʟᴀɴɢᴜᴀɢᴇ :</b> <code>Go</code>
» <b>ʟɪʙʀᴀʀʏ :</b> <code>gogram</code>

━━━━━━━━━━━━━━━━━━━
🍹 <b>ᴜsᴇ /help ꜰᴏʀ ᴄᴏᴍᴍᴀɴᴅs</b>
`

var HELP_TEXT = `
★ ᴘʙxsᴘᴀᴍ ɢᴏ ★

» <b>ᴄʟɪᴄᴋ ʙᴇʟᴏᴡ ʙᴜᴛᴛᴏɴs ꜰᴏʀ ʜᴇʟᴘ</b>
» <b>ᴅᴇᴠᴇʟᴏᴘᴇʀ:</b> @BadmundaXd
`

var SPAM_TEXT = `
<b>» sᴘᴀᴍ ᴄᴏᴍᴍᴀɴᴅs:</b>

<code>.spam [count] [text]</code>
➜ ɴᴏʀᴍᴀʟ sᴘᴀᴍ

<code>.ds [delay] [count] [text]</code>
➜ ᴅᴇʟᴀʏ sᴘᴀᴍ

<code>.sspam [count]</code>
➜ sᴛɪᴄᴋᴇʀ / ᴍᴇᴅɪᴀ sᴘᴀᴍ (ʀᴇᴘʟʏ)

<code>.hang [count]</code>
➜ ʜᴀɴɢ sᴘᴀᴍ

<code>.pspam [count]</code>
➜ ᴘᴏʀɴ ᴛᴇxᴛ sᴘᴀᴍ

<code>.stopspam</code>
➜ sᴛᴏᴘ ᴀʟʟ sᴘᴀᴍ

━━━━━━━━━━━━━━━━━
© @BadmundaXd
`

var RAID_TEXT = `
<b>» ʀᴀɪᴅ ᴄᴏᴍᴍᴀɴᴅs:</b>

<code>.raid [count]</code>
<code>.hraid [count]</code>
<code>.eraid [count]</code>
<code>.punraid [count]</code>
<code>.praid [count]</code>
➜ ɴᴏʀᴍᴀʟ ʀᴀɪᴅs

━━━━━━━━━━━━━━━━━

<code>.replyraid [count]</code>
<code>.hreplyraid [count]</code>
<code>.ereplyraid [count]</code>
<code>.preplyraid [count]</code>
➜ ʀᴇᴘʟʏ ʀᴀɪᴅs (ʀᴇᴘʟʏ ᴛᴏ ᴜsᴇʀ)

<code>.shayari [count]</code>
➜ sʜᴀʏᴀʀɪ ʀᴀɪᴅ

<code>.stopraid</code>
➜ sᴛᴏᴘ ᴀʟʟ ʀᴀɪᴅs

━━━━━━━━━━━━━━━━━
© @BadmundaXd
`

var EXTRA_TEXT = `
<b>» ᴇxᴛʀᴀ ᴄᴏᴍᴍᴀɴᴅs:</b>

<code>.ping</code>
➜ ʙᴏᴛ ᴘɪɴɢ

<code>.restart</code>
➜ ʀᴇsᴛᴀʀᴛ ʙᴏᴛ

<code>.logs</code>
➜ ꜰᴇᴛᴄʜ ʟᴏɢs

<code>.addsudo [reply/id]</code>
➜ ᴀᴅᴅ sᴜᴅᴏ ᴜsᴇʀ

<code>.rmsudo [reply/id]</code>
➜ ʀᴇᴍᴏᴠᴇ sᴜᴅᴏ ᴜsᴇʀ

<code>.sudolist</code>
➜ ʟɪsᴛ sᴜᴅᴏ ᴜsᴇʀs

━━━━━━━━━━━━━━━━━
© @BadmundaXd
`

// ─────────────────────────────────────────────
// KEYBOARDS
// ─────────────────────────────────────────────

func startKeyboard() *telegram.ReplyInlineMarkup {
	return &telegram.ReplyInlineMarkup{
		Rows: []*telegram.KeyboardButtonRow{
			{Buttons: []telegram.KeyboardButton{telegram.Button.Data("• ᴄᴏᴍᴍᴀɴᴅs •", "help_back")}},
			{Buttons: []telegram.KeyboardButton{
				telegram.Button.URL("• ᴄʜᴀɴɴᴇʟ •", "https://t.me/PBX_UPDATE"),
				telegram.Button.URL("• sᴜᴘᴘᴏʀᴛ •", "https://t.me/PBXCHATS"),
			}},
			{Buttons: []telegram.KeyboardButton{telegram.Button.URL("• ʀᴇᴘᴏ •", "https://github.com/badmunda05")}},
		},
	}
}

func helpKeyboard() *telegram.ReplyInlineMarkup {
	return &telegram.ReplyInlineMarkup{
		Rows: []*telegram.KeyboardButtonRow{
			{Buttons: []telegram.KeyboardButton{
				telegram.Button.Data("• sᴘᴀᴍ •", "spam_help"),
				telegram.Button.Data("• ʀᴀɪᴅ •", "raid_help"),
			}},
			{Buttons: []telegram.KeyboardButton{telegram.Button.Data("• ᴇxᴛʀᴀ •", "extra_help")}},
			{Buttons: []telegram.KeyboardButton{
				telegram.Button.URL("• ᴄʜᴀɴɴᴇʟ •", "https://t.me/PBX_UPDATE"),
				telegram.Button.URL("• sᴜᴘᴘᴏʀᴛ •", "https://t.me/PBXCHATS"),
			}},
			{Buttons: []telegram.KeyboardButton{telegram.Button.Data("• ʜᴏᴍᴇ •", "go_home")}},
		},
	}
}

func homeKeyboard() *telegram.ReplyInlineMarkup {
	return &telegram.ReplyInlineMarkup{
		Rows: []*telegram.KeyboardButtonRow{
			{Buttons: []telegram.KeyboardButton{telegram.Button.Data("• ʙᴀᴄᴋ •", "back_help")}},
		},
	}
}

// ─────────────────────────────────────────────
// /start
// ─────────────────────────────────────────────

func startHandler(m *telegram.NewMessage) error {
	text := fmt.Sprintf(START_TEXT, m.SenderID(), m.Sender.FirstName)
	_, err := m.ReplyMedia(config.AppConfig.StartPic, &telegram.MediaOptions{
		Caption: text, ReplyMarkup: startKeyboard(), ParseMode: telegram.HTML,
	})
	return err
}

// ─────────────────────────────────────────────
// /help
// ─────────────────────────────────────────────

func helpHandler(m *telegram.NewMessage) error {
	_, err := m.ReplyMedia(config.AppConfig.HelpPic, &telegram.MediaOptions{
		Caption: HELP_TEXT, ReplyMarkup: helpKeyboard(), ParseMode: telegram.HTML,
	})
	return err
}

// ─────────────────────────────────────────────
// .ping
// ─────────────────────────────────────────────

func pingHandler(m *telegram.NewMessage) error {
	start := time.Now()
	msg, err := m.Reply("🍹 ᴘɪɴɢɪɴɢ...", &telegram.SendOptions{ParseMode: telegram.HTML})
	speed := time.Since(start).Milliseconds()
	pingText := fmt.Sprintf(
		"•[ 🍹 ᴘʙx sᴘᴀᴍ ɢᴏ 🍹 ]•\n\n"+
			"» ᴘɪɴɢ  ➜ <code>%d ᴍs</code>\n"+
			"» sᴛᴀᴛᴜs ➜ <code>ᴏɴʟɪɴᴇ ✅</code>", speed)
	if err != nil || msg == nil {
		_, _ = m.Client.SendMessage(m.ChatID(), pingText, &telegram.SendOptions{ParseMode: telegram.HTML})
		return nil
	}
	if _, editErr := msg.Edit(pingText, &telegram.SendOptions{ParseMode: telegram.HTML}); editErr != nil {
		_, _ = m.Client.SendMessage(m.ChatID(), pingText, &telegram.SendOptions{ParseMode: telegram.HTML})
	}
	return nil
}

// ─────────────────────────────────────────────
// .restart
// ─────────────────────────────────────────────

func restartHandler(m *telegram.NewMessage) error {
	_, _ = Reply(m, "🔄 <b>ʀᴇsᴛᴀʀᴛɪɴɢ ᴘʙx sᴘᴀᴍ ɢᴏ...</b>")
	time.Sleep(1 * time.Second)
	exe, err := os.Executable()
	if err != nil {
		exe = os.Args[0]
	}
	cmd := exec.Command(exe, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	_ = cmd.Start()
	os.Exit(0)
	return nil
}

// ─────────────────────────────────────────────
// .logs
// ─────────────────────────────────────────────

func logsHandler(m *telegram.NewMessage) error {
	start := time.Now()
	msg, _ := Reply(m, "⚡ ꜰᴇᴛᴄʜɪɴɢ ʟᴏɢs...")
	cmd := exec.Command("bash", "-c", "pm2 logs --lines 150 --nostream 2>&1")
	output, err := cmd.CombinedOutput()
	if err != nil && len(output) == 0 {
		if msg != nil {
			_, _ = msg.Edit(fmt.Sprintf("❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ꜰᴇᴛᴄʜ ʟᴏɢs.\n\n<code>%s</code>", err.Error()), nil)
		}
		return nil
	}
	fileName := "PbxLogs.txt"
	logText := fmt.Sprintf("⚡ PBXGO BOT LOGS ⚡\n\n%s", string(output))
	if writeErr := os.WriteFile(fileName, []byte(logText), 0644); writeErr != nil {
		if msg != nil {
			_, _ = msg.Edit("❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ᴡʀɪᴛᴇ ʟᴏɢs ꜰɪʟᴇ.", nil)
		}
		return nil
	}
	defer os.Remove(fileName)
	taken := time.Since(start).Seconds()
	caption := fmt.Sprintf("⚡ <b>ᴘʙxɢᴏ ʟᴏɢs</b> ⚡\n» <b>ᴛɪᴍᴇ ᴛᴀᴋᴇɴ:</b> <code>%.0f sᴇᴄ</code>", taken)
	_, sendErr := m.Client.SendMedia(m.ChatID(), fileName, &telegram.MediaOptions{Caption: caption, ParseMode: telegram.HTML})
	if sendErr != nil && strings.Contains(strings.ToUpper(sendErr.Error()), "FLOOD_WAIT") {
		time.Sleep(5 * time.Second)
		_, _ = m.Client.SendMedia(m.ChatID(), fileName, &telegram.MediaOptions{Caption: caption, ParseMode: telegram.HTML})
	}
	if msg != nil {
		_, _ = msg.Delete()
	}
	return nil
}

// ─────────────────────────────────────────────
// .addsudo
// ─────────────────────────────────────────────

func addSudoHandler(m *telegram.NewMessage) error {
	if m.SenderID() != config.AppConfig.OwnerID {
		_, _ = Reply(m, "❌ <b>ᴏɴʟʏ ᴏᴡɴᴇʀ ᴄᴀɴ ᴀᴅᴅ sᴜᴅᴏ ᴜsᴇʀs.</b>")
		return nil
	}
	var targetID int64
	if m.IsReply() {
		replyMsg, err := m.GetReplyMessage()
		if err != nil || replyMsg == nil {
			_, _ = Reply(m, "❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ɢᴇᴛ ʀᴇᴘʟɪᴇᴅ ᴍᴇssᴀɢᴇ.")
			return nil
		}
		targetID = replyMsg.SenderID()
	} else {
		args := GetArgs(m)
		if args == "" {
			_, _ = Reply(m, "⚠️ <b>ᴜsᴀɢᴇ:</b>\n<code>.addsudo [user_id]</code>\nᴏʀ ʀᴇᴘʟʏ ᴛᴏ ᴀ ᴜsᴇʀ.")
			return nil
		}
		parsed, err := strconv.ParseInt(strings.TrimSpace(args), 10, 64)
		if err != nil {
			_, _ = Reply(m, "❌ ɪɴᴠᴀʟɪᴅ ᴜsᴇʀ ɪᴅ.")
			return nil
		}
		targetID = parsed
	}
	if targetID == config.AppConfig.OwnerID {
		_, _ = Reply(m, "ℹ️ ᴏᴡɴᴇʀ ɪs ᴀʟʀᴇᴀᴅʏ sᴜᴘᴇʀ ᴀᴅᴍɪɴ.")
		return nil
	}
	database.AddSudo(targetID)
	_, _ = Reply(m, fmt.Sprintf("✅ <b>ᴜsᴇʀ <code>%d</code> ᴀᴅᴅᴇᴅ ᴛᴏ sᴜᴅᴏ ʟɪsᴛ.</b>", targetID))
	return nil
}

// ─────────────────────────────────────────────
// .rmsudo
// ─────────────────────────────────────────────

func rmSudoHandler(m *telegram.NewMessage) error {
	if m.SenderID() != config.AppConfig.OwnerID {
		_, _ = Reply(m, "❌ <b>ᴏɴʟʏ ᴏᴡɴᴇʀ ᴄᴀɴ ʀᴇᴍᴏᴠᴇ sᴜᴅᴏ ᴜsᴇʀs.</b>")
		return nil
	}
	var targetID int64
	if m.IsReply() {
		replyMsg, err := m.GetReplyMessage()
		if err != nil || replyMsg == nil {
			_, _ = Reply(m, "❌ ꜰᴀɪʟᴇᴅ ᴛᴏ ɢᴇᴛ ʀᴇᴘʟɪᴇᴅ ᴍᴇssᴀɢᴇ.")
			return nil
		}
		targetID = replyMsg.SenderID()
	} else {
		args := GetArgs(m)
		if args == "" {
			_, _ = Reply(m, "⚠️ <b>ᴜsᴀɢᴇ:</b>\n<code>.rmsudo [user_id]</code>\nᴏʀ ʀᴇᴘʟʏ ᴛᴏ ᴀ ᴜsᴇʀ.")
			return nil
		}
		parsed, err := strconv.ParseInt(strings.TrimSpace(args), 10, 64)
		if err != nil {
			_, _ = Reply(m, "❌ ɪɴᴠᴀʟɪᴅ ᴜsᴇʀ ɪᴅ.")
			return nil
		}
		targetID = parsed
	}
	if targetID == config.AppConfig.OwnerID {
		_, _ = Reply(m, "❌ ᴄᴀɴɴᴏᴛ ʀᴇᴍᴏᴠᴇ ᴏᴡɴᴇʀ ꜰʀᴏᴍ sᴜᴅᴏ.")
		return nil
	}
	database.RemoveSudo(targetID)
	_, _ = Reply(m, fmt.Sprintf("✅ <b>ᴜsᴇʀ <code>%d</code> ʀᴇᴍᴏᴠᴇᴅ ꜰʀᴏᴍ sᴜᴅᴏ ʟɪsᴛ.</b>", targetID))
	return nil
}

// ─────────────────────────────────────────────
// .sudolist
// ─────────────────────────────────────────────

func sudoListHandler(m *telegram.NewMessage) error {
	list := database.FetchSudoList()
	if len(list) == 0 {
		_, _ = Reply(m, "ℹ️ <b>ɴᴏ sᴜᴅᴏ ᴜsᴇʀs ᴀᴅᴅᴇᴅ ʏᴇᴛ.</b>")
		return nil
	}
	text := "👑 <b>sᴜᴅᴏ ᴜsᴇʀs ʟɪsᴛ:</b>\n\n"
	for i, uid := range list {
		text += fmt.Sprintf("%d. <code>%d</code>\n", i+1, uid)
	}
	text += fmt.Sprintf("\n<b>ᴛᴏᴛᴀʟ:</b> <code>%d</code>", len(list))
	_, _ = Reply(m, text)
	return nil
}

// ─────────────────────────────────────────────
// CALLBACKS
// ─────────────────────────────────────────────

func callbackHandler(c *telegram.CallbackQuery) error {
	data := c.DataString()
	switch data {
	case "help_back":
		c.Edit(HELP_TEXT, &telegram.SendOptions{ParseMode: telegram.HTML, ReplyMarkup: helpKeyboard()})
		c.Answer("")
	case "spam_help":
		c.Edit(SPAM_TEXT, &telegram.SendOptions{ParseMode: telegram.HTML, ReplyMarkup: homeKeyboard()})
		c.Answer("")
	case "raid_help":
		c.Edit(RAID_TEXT, &telegram.SendOptions{ParseMode: telegram.HTML, ReplyMarkup: homeKeyboard()})
		c.Answer("")
	case "extra_help":
		c.Edit(EXTRA_TEXT, &telegram.SendOptions{ParseMode: telegram.HTML, ReplyMarkup: homeKeyboard()})
		c.Answer("")
	case "back_help":
		c.Edit(HELP_TEXT, &telegram.SendOptions{ParseMode: telegram.HTML, ReplyMarkup: helpKeyboard()})
		c.Answer("")
	case "go_home":
		sender, _ := c.GetSender()
		name := "ᴜsᴇʀ"
		id := c.SenderID
		if sender != nil {
			name = sender.FirstName
		}
		c.Edit(fmt.Sprintf(START_TEXT, id, name), &telegram.SendOptions{ParseMode: telegram.HTML, ReplyMarkup: startKeyboard()})
		c.Answer("")
	}
	return nil
}

// ─────────────────────────────────────────────
// REGISTER
// ─────────────────────────────────────────────

func init() {
	Register(ModuleInfo{
		Name:        "Basic",
		Description: "Start Help Ping Restart Logs Sudo",
		Commands: []CommandInfo{
			{Pattern: "start",    Handler: startHandler,    Sudo: false},
			{Pattern: "help",     Handler: helpHandler,     Sudo: false},
			{Pattern: "ping",     Handler: pingHandler,     Sudo: true},
			{Pattern: "restart",  Handler: restartHandler,  Sudo: true},
			{Pattern: "logs",     Handler: logsHandler,     Sudo: true},
			{Pattern: "addsudo",  Handler: addSudoHandler,  Sudo: false},
			{Pattern: "rmsudo",   Handler: rmSudoHandler,   Sudo: false},
			{Pattern: "sudolist", Handler: sudoListHandler, Sudo: true},
		},
		Callbacks: []CallbackInfo{
			{Pattern: "help_back",  Handler: callbackHandler},
			{Pattern: "spam_help",  Handler: callbackHandler},
			{Pattern: "raid_help",  Handler: callbackHandler},
			{Pattern: "extra_help", Handler: callbackHandler},
			{Pattern: "back_help",  Handler: callbackHandler},
			{Pattern: "go_home",    Handler: callbackHandler},
		},
	})
}
