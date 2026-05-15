package modules

import (
	"fmt"
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

var IsBotMode bool

// Reply sends a reply with flood wait handling (max 3 retries)
func Reply(m *telegram.NewMessage, text string) (*telegram.NewMessage, error) {
	retries := 0
	for retries < 3 {
		msg, err := m.Reply(text, &telegram.SendOptions{
			ParseMode: telegram.HTML,
		})
		if err != nil {
			errText := strings.ToUpper(err.Error())
			if strings.Contains(errText, "FLOOD_WAIT") {
				wait := extractFloodWait(errText)
				if wait <= 0 {
					wait = 5
				}
				time.Sleep(time.Duration(wait+1) * time.Second)
				retries++
				continue
			}
			// Reply failed — fallback to SendMessage
			_, _ = m.Client.SendMessage(m.ChatID(), text, &telegram.SendOptions{
				ParseMode: telegram.HTML,
			})
			return nil, err
		}
		return msg, nil
	}
	return nil, fmt.Errorf("reply failed after retries")
}

// GetArgs returns everything after the first word/command
func GetArgs(m *telegram.NewMessage) string {
	text := m.Text()
	if idx := findSpace(text); idx != -1 {
		return text[idx+1:]
	}
	return ""
}

func findSpace(s string) int {
	for i, c := range s {
		if c == ' ' {
			return i
		}
	}
	return -1
}

// extractFloodWait parses the wait seconds from a FLOOD_WAIT error string
func extractFloodWait(errText string) int {
	var wait int
	_, _ = fmt.Sscanf(errText, "FLOOD_WAIT_%d", &wait)
	return wait
}
