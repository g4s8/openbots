package logger

import (
	"fmt"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type logger struct {
	ogger zerolog.Logger
	evel  zerolog.Level
}

// Wrap zerolog.Logger to implement telegram.BotLogger interface.
func Wrap(log zerolog.Logger, lvl zerolog.Level) telegram.BotLogger {
	return &logger{ogger: log, evel: lvl}
}

func (l *logger) Println(v ...interface{}) {
	l.ogger.WithLevel(l.evel).Msg(strings.TrimSpace(fmt.Sprintln(v...)))
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.ogger.WithLevel(l.evel).Msg(strings.TrimSpace(fmt.Sprintf(format, v...)))
}
