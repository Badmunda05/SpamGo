// modules/module.go

package modules

import (
	"strings"
	"time"

	"pbxgo/config"
	"pbxgo/database"

	"github.com/amarnathcjd/gogram/telegram"
)

type CommandInfo struct {
	Pattern string
	Handler func(*telegram.NewMessage) error
	Sudo    bool
}

type CallbackInfo struct {
	Pattern string
	Handler func(*telegram.CallbackQuery) error
}

type ModuleInfo struct {
	Name        string
	Description string
	Commands    []CommandInfo
	Callbacks   []CallbackInfo
}

var allModules []ModuleInfo

func Register(m ModuleInfo) {
	allModules = append(allModules, m)
}

// safeSendMessage sends a message with one flood wait retry, then gives up
func safeSendMessage(c *telegram.Client, chatID int64, text string, opts *telegram.SendOptions) {
	_, err := c.SendMessage(chatID, text, opts)
	if err != nil {
		errText := strings.ToUpper(err.Error())
		if strings.Contains(errText, "FLOOD_WAIT") {
			wait := extractFloodWait(errText)
			if wait <= 0 {
				wait = 5
			}
			time.Sleep(time.Duration(wait+1) * time.Second)
			_, _ = c.SendMessage(chatID, text, opts)
		}
	}
}

// isSudoOrOwner returns true if the sender is the owner or in the sudo list
func isSudoOrOwner(senderID int64) bool {
	if senderID == config.AppConfig.OwnerID {
		return true
	}
	return database.IsSudo(senderID)
}

// Load registers all message and callback handlers for a given client
func Load(c *telegram.Client) {

	c.On("message", func(m *telegram.NewMessage) error {
		if m == nil {
			return nil
		}

		// Check reply-raid watcher before parsing any command
		TriggerReplyRaidIfActive(m)

		text := m.Text()
		if text == "" {
			return nil
		}

		parts := strings.Fields(text)
		if len(parts) == 0 {
			return nil
		}

		var cmd string
		if strings.HasPrefix(text, "/") {
			cmd = strings.TrimPrefix(parts[0], "/")
		} else if strings.HasPrefix(text, ".") {
			cmd = strings.TrimPrefix(parts[0], ".")
		} else {
			return nil
		}

		// Strip @botusername suffix if present
		if idx := strings.Index(cmd, "@"); idx != -1 {
			cmd = cmd[:idx]
		}

		for _, mod := range allModules {
			for _, command := range mod.Commands {
				if strings.EqualFold(command.Pattern, cmd) {
					if command.Sudo && !isSudoOrOwner(m.SenderID()) {
						safeSendMessage(
							c,
							m.ChatID(),
							"❌ <b>ᴏɴʟʏ ᴏᴡɴᴇʀ / sᴜᴅᴏ ᴜsᴇʀs ᴄᴀɴ ᴜsᴇ ᴛʜɪs.</b>",
							&telegram.SendOptions{ParseMode: telegram.HTML},
						)
						return nil
					}
					return command.Handler(m)
				}
			}
		}

		return nil
	})

	c.On("callback", func(cb *telegram.CallbackQuery) error {
		if cb == nil {
			return nil
		}
		data := cb.DataString()
		for _, mod := range allModules {
			for _, callback := range mod.Callbacks {
				if callback.Pattern == data {
					return callback.Handler(cb)
				}
			}
		}
		return nil
	})
}
