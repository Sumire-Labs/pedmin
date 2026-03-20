package embedfix

import (
	"context"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (ef *EmbedFix) handleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	// embedfix:translate:{screenName}:{tweetID}
	parts := strings.SplitN(customID, ":", 4)
	if len(parts) < 4 || parts[1] != "translate" {
		return
	}

	screenName := parts[2]
	tweetID := parts[3]

	_ = e.DeferUpdateMessage()

	if ef.translateClient.apiKey == "" {
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.NewMessageUpdateV2([]discord.LayoutComponent{
				discord.NewContainer(
					discord.NewTextDisplay("翻訳APIキーが設定されていないため、翻訳できません。"),
				),
			}))
		return
	}

	ref := tweetRef{ScreenName: screenName, TweetID: tweetID}

	ctx := context.Background()
	tweet, err := ef.fxClient.GetTweet(ctx, screenName, tweetID)
	if err != nil {
		ef.logger.Warn("failed to fetch tweet for translation",
			slog.String("tweet_id", tweetID),
			slog.Any("error", err),
		)
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.NewMessageUpdateV2([]discord.LayoutComponent{
				discord.NewContainer(
					discord.NewTextDisplay("ツイートの取得に失敗しました。"),
				),
			}))
		return
	}

	result, err := ef.translateClient.Translate(ctx, tweet.Text, "ja")
	if err != nil {
		ef.logger.Warn("failed to translate tweet",
			slog.String("tweet_id", tweetID),
			slog.Any("error", err),
		)
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.NewMessageUpdateV2([]discord.LayoutComponent{
				discord.NewContainer(
					discord.NewTextDisplay("翻訳に失敗しました。"),
				),
			}))
		return
	}

	ui := BuildTweetEmbedTranslated(tweet, result, ref)
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2(ui))
}
