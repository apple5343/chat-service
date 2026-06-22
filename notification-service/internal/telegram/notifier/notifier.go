package notifier

import (
	tele "gopkg.in/telebot.v3"

	"fmt"
	"time"
)

// Notifier представляет телеграм бота, уведомляющего о событиях.
type Notifier struct {
	bot *tele.Bot
}

// NewNotifier инициализирует Notifier.
func NewNotifier(botToken string) (*Notifier, error) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &Notifier{bot: bot}, nil
}

// Notify отправляет сообщение в чат с топиком.
func (n *Notifier) Notify(chatID int64, message string) error {
	_, err := n.bot.Send(&tele.Chat{ID: chatID}, message, &tele.SendOptions{
		ParseMode: tele.ModeHTML,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
