package translator

import (
	"context"
	"log/slog"

	"github.com/disgoorg/snowflake/v2"
)

func (t *Translator) processTranslation(ctx context.Context, channelID, messageID snowflake.ID, targetLang string) {
	msg, err := t.client.Rest.GetMessage(channelID, messageID)
	if err != nil {
		t.logger.Warn("failed to fetch message for translation",
			slog.String("message_id", messageID.String()),
			slog.Any("error", err),
		)
		return
	}

	if msg.Author.Bot || msg.Content == "" {
		return
	}

	result, err := t.translateClient.Translate(ctx, msg.Content, targetLang)
	if err != nil {
		t.logger.Warn("failed to translate message",
			slog.String("message_id", messageID.String()),
			slog.Any("error", err),
		)
		return
	}

	embed := BuildTranslationEmbed(result.TranslatedText, result.DetectedLanguage, targetLang, msg.Author.ID, messageID)
	if _, err := t.client.Rest.CreateMessage(channelID, embed); err != nil {
		t.logger.Warn("failed to send translation",
			slog.String("message_id", messageID.String()),
			slog.Any("error", err),
		)
	}
}
