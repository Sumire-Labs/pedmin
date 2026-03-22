package translator

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func BuildTranslationEmbed(translatedText string, sourceLang, targetLang string, authorID snowflake.ID, messageID snowflake.ID) discord.MessageCreate {
	header := fmt.Sprintf("🌐 **翻訳** (%s → %s)", langName(sourceLang), langName(targetLang))

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(header),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(translatedText),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(fmt.Sprintf("-# 翻訳元: <@%d>", authorID)),
	}

	return discord.NewMessageCreateV2(
		discord.NewContainer(components...),
	).WithMessageReferenceByID(messageID)
}
