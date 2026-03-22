package embedfix

import (
	"context"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (ef *EmbedFix) handleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, rest, _ := strings.Cut(rest, ":")
	if action != "translate" {
		return
	}

	_ = e.DeferUpdateMessage()

	if ef.translateClient.apiKey == "" {
		ef.respondTranslateError(e, "翻訳APIキーが設定されていないため、翻訳できません。")
		return
	}

	platform, params, _ := strings.Cut(rest, ":")
	ctx := context.Background()

	ui, err := ef.translateContent(ctx, platform, params)
	if err != nil {
		ef.respondTranslateError(e, "翻訳に失敗しました。")
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2(ui))
}

func (ef *EmbedFix) respondTranslateError(e *events.ComponentInteractionCreate, msg string) {
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay(msg),
			),
		}))
}
