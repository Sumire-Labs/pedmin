package settings

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func (s *Settings) mainPanel(guildID snowflake.ID) discord.MessageCreate {
	return ephemeralV2(s.buildMainContainer(guildID))
}

func (s *Settings) mainPanelUpdate(guildID snowflake.ID) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{s.buildMainContainer(guildID)})
}

func (s *Settings) buildMainContainer(guildID snowflake.ID) discord.ContainerComponent {
	modules := s.bot.GetModules()
	var options []discord.StringSelectMenuOption
	for _, m := range modules {
		info := m.Info()
		if info.AlwaysOn {
			continue
		}
		status := "❌"
		if s.bot.IsModuleEnabled(guildID, info.ID) {
			status = "✅"
		}
		options = append(options, discord.StringSelectMenuOption{
			Label:       fmt.Sprintf("%s %s", status, info.Name),
			Value:       info.ID,
			Description: info.Description,
		})
	}

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay("## 🐸 Pedmin Settings"),
		discord.NewLargeSeparator(),
	}

	if len(options) > 0 {
		components = append(components,
			discord.NewTextDisplay("Select a module to configure:"),
			discord.NewActionRow(
				discord.NewStringSelectMenu(ModuleID+":select", "Choose a module...", options...),
			),
		)
	} else {
		components = append(components,
			discord.NewTextDisplay("No configurable modules registered."),
		)
	}

	return discord.NewContainer(components...).WithAccentColor(0x00B894)
}

func (s *Settings) modulePanel(guildID snowflake.ID, moduleID string) discord.MessageUpdate {
	modules := s.bot.GetModules()
	m, ok := modules[moduleID]
	if !ok {
		return discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay("Module not found."),
			).WithAccentColor(0xFF0000),
		})
	}

	info := m.Info()
	enabled := s.bot.IsModuleEnabled(guildID, moduleID)

	statusText := "Disabled"
	toggleLabel := "Enable"
	toggleStyle := discord.ButtonStyleSuccess
	accentColor := 0xFF6B6B
	if enabled {
		statusText = "Enabled"
		toggleLabel = "Disable"
		toggleStyle = discord.ButtonStyleDanger
		accentColor = 0x00B894
	}

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(fmt.Sprintf("## %s", info.Name)),
		discord.NewTextDisplay(info.Description),
		discord.NewLargeSeparator(),
		discord.NewTextDisplay(fmt.Sprintf("**Status:** %s", statusText)),
	}

	settingsPanel := m.SettingsPanel(guildID)
	if len(settingsPanel) > 0 {
		components = append(components, discord.NewLargeSeparator())
		for _, lc := range settingsPanel {
			if sub, ok := lc.(discord.ContainerSubComponent); ok {
				components = append(components, sub)
			}
		}
	}

	components = append(components,
		discord.NewLargeSeparator(),
		discord.NewActionRow(
			discord.NewButton(toggleStyle, toggleLabel, fmt.Sprintf("%s:toggle:%s", ModuleID, moduleID), "", 0),
			discord.NewSecondaryButton("← Back", ModuleID+":back"),
		),
	)

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(components...).WithAccentColor(accentColor),
	})
}

func ephemeralV2(components ...discord.LayoutComponent) discord.MessageCreate {
	return discord.NewMessageCreateV2(components...).WithEphemeral(true)
}
