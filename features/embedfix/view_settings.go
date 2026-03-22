package embedfix

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
)

func BuildSettingsPanel(settings *EmbedFixSettings) []discord.LayoutComponent {
	var enabledNames []string
	for _, p := range AllPlatforms {
		if settings.IsPlatformEnabled(p.Key) {
			enabledNames = append(enabledNames, p.Label)
		}
	}
	statusText := "なし"
	if len(enabledNames) > 0 {
		statusText = strings.Join(enabledNames, ", ")
	}

	infoDisplay := discord.NewTextDisplay("**埋め込み対象:** " + statusText)

	var options []discord.StringSelectMenuOption
	for _, p := range AllPlatforms {
		opt := discord.StringSelectMenuOption{
			Label: p.Label,
			Value: string(p.Key),
		}
		if settings.IsPlatformEnabled(p.Key) {
			opt.Default = true
		}
		options = append(options, opt)
	}

	platformSelect := discord.NewActionRow(
		discord.NewStringSelectMenu(ModuleID+":platforms", "埋め込み対象を選択...", options...).
			WithMinValues(0).
			WithMaxValues(len(AllPlatforms)),
	)

	return []discord.LayoutComponent{
		infoDisplay,
		platformSelect,
	}
}
