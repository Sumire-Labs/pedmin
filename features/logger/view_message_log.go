package logger

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func BuildMessageEditLog(user discord.User, channelID snowflake.ID, oldContent, newContent string, oldAttachments, newAttachments []discord.Attachment) discord.MessageCreate {
	title := "### ✏️ メッセージ編集"
	body := fmt.Sprintf("**ユーザー:** <@%d>\n**チャンネル:** <#%d>\n**変更前:**\n> %s\n**変更後:**\n> %s",
		user.ID, channelID, oldContent, newContent)

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(title),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(body),
	}

	removed, added := diffAttachments(oldAttachments, newAttachments)
	if len(removed) > 0 {
		components = append(components, discord.NewSmallSeparator())
		components = append(components, discord.NewTextDisplay("**削除された添付ファイル:**"))
		components = append(components, buildAttachmentComponents(removed)...)
	}
	if len(added) > 0 {
		components = append(components, discord.NewSmallSeparator())
		components = append(components, discord.NewTextDisplay("**追加された添付ファイル:**"))
		components = append(components, buildAttachmentComponents(added)...)
	}

	return discord.NewMessageCreateV2(
		discord.NewContainer(components...),
	).WithAllowedMentions(&discord.AllowedMentions{})
}

func BuildMessageDeleteLog(user *discord.User, channelID snowflake.ID, content string, attachments []discord.Attachment) discord.MessageCreate {
	userText := "*不明*"
	if user != nil {
		userText = fmt.Sprintf("<@%d>", user.ID)
	}
	contentText := content
	if contentText == "" {
		contentText = "*内容を取得できませんでした*"
	}

	title := "### 🗑️ メッセージ削除"
	body := fmt.Sprintf("**ユーザー:** %s\n**チャンネル:** <#%d>\n**内容:**\n> %s",
		userText, channelID, contentText)

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(title),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(body),
	}

	if len(attachments) > 0 {
		components = append(components, discord.NewSmallSeparator())
		components = append(components, discord.NewTextDisplay("**添付ファイル:**"))
		components = append(components, buildAttachmentComponents(attachments)...)
	}

	return discord.NewMessageCreateV2(
		discord.NewContainer(components...),
	).WithAllowedMentions(&discord.AllowedMentions{})
}
