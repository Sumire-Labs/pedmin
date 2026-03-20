package rss

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func BuildSettingsPanel(feedCount int) []discord.LayoutComponent {
	infoDisplay := discord.NewTextDisplay(
		fmt.Sprintf("**登録フィード:** %d/%d", feedCount, MaxFeedsPerGuild),
	)

	manageBtn := discord.NewSecondaryButton("フィード管理", ModuleID+":manage")
	if feedCount == 0 {
		manageBtn = manageBtn.AsDisabled()
	}

	actionRow := discord.NewActionRow(
		discord.NewPrimaryButton("フィード追加", ModuleID+":add_prompt"),
		manageBtn,
	)

	return []discord.LayoutComponent{
		infoDisplay,
		actionRow,
	}
}
