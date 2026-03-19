package player

import (
	"log/slog"
	"sync"
	"time"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/module"
)

const ModuleID = "player"

type trackedMessage struct {
	channelID snowflake.ID
	messageID snowflake.ID
}

type Player struct {
	lavalink         disgolink.Client
	client           *disgobot.Client
	queues           *QueueManager
	messages         sync.Map // map[snowflake.ID]trackedMessage
	defaultVolume    int
	autoLeaveTimeout time.Duration
	leaveTimers      sync.Map // map[snowflake.ID]*time.Timer
	logger           *slog.Logger
}

func New(link disgolink.Client, client *disgobot.Client, defaultVolume int, autoLeaveTimeout time.Duration, logger *slog.Logger) *Player {
	return &Player{
		lavalink:         link,
		client:           client,
		queues:           NewQueueManager(),
		defaultVolume:    defaultVolume,
		autoLeaveTimeout: autoLeaveTimeout,
		logger:           logger,
	}
}

func (p *Player) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "ミュージックプレイヤー",
		Description: "様々なソースから音楽を再生するミュージックプレイヤー",
		AlwaysOn:    false,
	}
}

func (p *Player) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "player",
			Description: "ミュージックプレイヤーを表示",
		},
	}
}

func (p *Player) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}

func (p *Player) HandleSettingsComponent(_ *events.ComponentInteractionCreate) {}
