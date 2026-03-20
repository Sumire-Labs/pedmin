package ticket

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func BuildSettingsPanel(settings *TicketSettings) []discord.LayoutComponent {
	categoryText := "未設定"
	if settings.CategoryID != 0 {
		categoryText = fmt.Sprintf("<#%d>", settings.CategoryID)
	}
	logText := "未設定"
	if settings.LogChannelID != 0 {
		logText = fmt.Sprintf("<#%d>", settings.LogChannelID)
	}
	roleText := "未設定"
	if settings.SupportRoleID != 0 {
		roleText = fmt.Sprintf("<@&%d>", settings.SupportRoleID)
	}

	infoDisplay := discord.NewTextDisplay(fmt.Sprintf(
		"**カテゴリ:** %s\n**ログチャンネル:** %s\n**サポートロール:** %s",
		categoryText, logText, roleText,
	))

	categorySelect := discord.NewActionRow(
		discord.NewChannelSelectMenu(ModuleID+":category", "カテゴリを選択...").
			WithChannelTypes(discord.ChannelTypeGuildCategory),
	)

	buttons := discord.NewActionRow(
		discord.NewSecondaryButton("ログ設定", ModuleID+":log_prompt"),
		discord.NewSecondaryButton("ロール設定", ModuleID+":role_prompt"),
		discord.NewSuccessButton("パネル設置", ModuleID+":deploy_prompt"),
	)

	return []discord.LayoutComponent{
		infoDisplay,
		categorySelect,
		buttons,
	}
}
