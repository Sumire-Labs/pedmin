package embedfix

import (
	"context"
	"log/slog"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func SetupListeners(client *disgobot.Client, ef *EmbedFix) {
	client.AddEventListeners(
		disgobot.NewListenerFunc(ef.onMessageCreate),
	)
}

func (ef *EmbedFix) onMessageCreate(e *events.GuildMessageCreate) {
	if e.Message.Author.Bot {
		return
	}

	guildID := e.GuildID
	if !ef.bot.IsModuleEnabled(guildID, ModuleID) {
		return
	}

	refs := extractTweetURLs(e.Message.Content)
	if len(refs) == 0 {
		return
	}

	// Suppress embeds on the original message (best-effort)
	_, err := ef.client.Rest.UpdateMessage(e.ChannelID, e.MessageID,
		discord.NewMessageUpdate().WithSuppressEmbeds(true))
	if err != nil {
		ef.logger.Debug("failed to suppress embeds",
			slog.Any("error", err),
			slog.String("message_id", e.MessageID.String()),
		)
	}

	ctx := context.Background()
	for _, ref := range refs {
		tweet, err := ef.fxClient.GetTweet(ctx, ref.ScreenName, ref.TweetID)
		if err != nil {
			ef.logger.Warn("failed to fetch tweet",
				slog.String("screen_name", ref.ScreenName),
				slog.String("tweet_id", ref.TweetID),
				slog.Any("error", err),
			)
			continue
		}

		msg := BuildTweetEmbed(tweet, ref)
		_, err = ef.client.Rest.CreateMessage(e.ChannelID,
			msg.WithMessageReferenceByID(e.MessageID))
		if err != nil {
			ef.logger.Warn("failed to send tweet embed",
				slog.String("tweet_id", ref.TweetID),
				slog.Any("error", err),
			)
		}
	}
}
