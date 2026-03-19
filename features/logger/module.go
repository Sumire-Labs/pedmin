package logger

import (
	"log/slog"
	"strings"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/module"
	"github.com/s12kuma01/pedmin/store"
)

const ModuleID = "logger"

type Bot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

type Logger struct {
	bot    Bot
	client *disgobot.Client
	store  store.GuildStore
	logger *slog.Logger
}

func New(bot Bot, client *disgobot.Client, guildStore store.GuildStore, logger *slog.Logger) *Logger {
	return &Logger{
		bot:    bot,
		client: client,
		store:  guildStore,
		logger: logger,
	}
}

func (l *Logger) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "Logger",
		Description: "サーバーイベントのログを記録",
		AlwaysOn:    false,
	}
}

func (l *Logger) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (l *Logger) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (l *Logger) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, action, _ := strings.Cut(customID, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	settings, err := LoadSettings(l.store, *guildID)
	if err != nil {
		l.logger.Error("failed to load logger settings", slog.Any("error", err))
		return
	}

	switch action {
	case "channel":
		data := e.Data.(discord.ChannelSelectMenuInteractionData)
		if len(data.Values) > 0 {
			settings.ChannelID = data.Values[0]
		}

	case "events":
		data := e.Data.(discord.StringSelectMenuInteractionData)
		for k := range settings.Events {
			settings.Events[k] = false
		}
		for _, v := range data.Values {
			settings.Events[v] = true
		}

	default:
		return
	}

	if err := SaveSettings(l.store, *guildID, settings); err != nil {
		l.logger.Error("failed to save logger settings", slog.Any("error", err))
	}

	_ = e.DeferUpdateMessage()
}

func (l *Logger) HandleModal(_ *events.ModalSubmitInteractionCreate) {}

func (l *Logger) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	settings, err := LoadSettings(l.store, guildID)
	if err != nil {
		l.logger.Error("failed to load logger settings", slog.Any("error", err))
		settings = &LoggerSettings{Events: make(map[string]bool)}
	}
	return BuildSettingsPanel(settings)
}

func (l *Logger) HandleSettingsComponent(_ *events.ComponentInteractionCreate) {}
