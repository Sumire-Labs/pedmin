package logger

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
)

func BuildSettingsPanel(settings *LoggerSettings) []discord.LayoutComponent {
	// Build info text
	channelText := "未設定"
	if settings.ChannelID != 0 {
		channelText = fmt.Sprintf("<#%d>", settings.ChannelID)
	}

	var enabledNames []string
	for _, ev := range AllEvents {
		if settings.IsEventEnabled(ev.Key) {
			enabledNames = append(enabledNames, ev.Label)
		}
	}
	eventsText := "なし"
	if len(enabledNames) > 0 {
		eventsText = strings.Join(enabledNames, ", ")
	}

	infoDisplay := discord.NewTextDisplay(fmt.Sprintf(
		"**ログチャンネル:** %s\n**ログ対象:** %s",
		channelText, eventsText,
	))

	// Channel select menu
	channelSelect := discord.NewActionRow(
		discord.NewChannelSelectMenu(ModuleID+":channel", "ログチャンネルを選択...").
			WithChannelTypes(discord.ChannelTypeGuildText),
	)

	// Event select menu
	var options []discord.StringSelectMenuOption
	for _, ev := range AllEvents {
		opt := discord.StringSelectMenuOption{
			Label: ev.Label,
			Value: ev.Key,
		}
		if settings.IsEventEnabled(ev.Key) {
			opt.Default = true
		}
		options = append(options, opt)
	}
	eventSelect := discord.NewActionRow(
		discord.NewStringSelectMenu(ModuleID+":events", "ログ対象イベントを選択...", options...).
			WithMinValues(0).
			WithMaxValues(len(AllEvents)),
	)

	return []discord.LayoutComponent{
		infoDisplay,
		channelSelect,
		eventSelect,
	}
}
