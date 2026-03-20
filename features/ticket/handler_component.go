package ticket

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (t *Ticket) handleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	parts := strings.SplitN(customID, ":", 3)
	if len(parts) < 2 {
		return
	}
	action := parts[1]

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch action {
	case "category":
		t.handleCategorySelect(e, *guildID)
	case "log_prompt":
		t.handleLogPrompt(e)
	case "log_channel":
		t.handleLogChannelSelect(e, *guildID)
	case "role_prompt":
		t.handleRolePrompt(e)
	case "role":
		t.handleRoleSelect(e, *guildID)
	case "deploy_prompt":
		t.handleDeployPrompt(e)
	case "deploy_channel":
		t.handleDeployChannelSelect(e)
	case "deploy_confirm":
		if len(parts) < 3 {
			return
		}
		t.handleDeployConfirm(e, parts[2])
	case "deploy_cancel":
		_ = e.DeferUpdateMessage()
	case "create":
		if !t.bot.IsModuleEnabled(*guildID, ModuleID) {
			return
		}
		_ = e.Modal(discord.ModalCreate{
			CustomID: ModuleID + ":create_modal",
			Title:    "チケットを作成",
			Components: []discord.LayoutComponent{
				discord.NewLabel("件名",
					discord.NewShortTextInput(ModuleID+":subject").
						WithPlaceholder("チケットの件名を入力").
						WithRequired(true).
						WithMaxLength(100),
				),
				discord.NewLabel("説明",
					discord.NewParagraphTextInput(ModuleID+":description").
						WithPlaceholder("詳しい内容を入力してください").
						WithRequired(false).
						WithMaxLength(1000),
				),
			},
		})
	case "close":
		t.archiveTicket(e, *guildID)
	case "delete":
		t.deleteTicket(e, *guildID)
	}
}
