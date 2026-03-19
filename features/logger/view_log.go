package logger

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

const (
	colorMessageEdit   = 0xF39C12
	colorMessageDelete = 0xE74C3C
	colorMemberJoin    = 0x2ECC71
	colorMemberLeave   = 0xE67E22
	colorBan           = 0xE74C3C
	colorUnban         = 0x3498DB
	colorRole          = 0x9B59B6
	colorChannel       = 0x1ABC9C
)

func BuildMessageEditLog(user discord.User, channelID snowflake.ID, oldContent, newContent string) discord.MessageCreate {
	return logMessage(colorMessageEdit,
		"### ✏️ メッセージ編集",
		fmt.Sprintf("**ユーザー:** <@%d>\n**チャンネル:** <#%d>\n**変更前:**\n> %s\n**変更後:**\n> %s",
			user.ID, channelID, oldContent, newContent),
	)
}

func BuildMessageDeleteLog(user *discord.User, channelID snowflake.ID, content string) discord.MessageCreate {
	userText := "*不明*"
	if user != nil {
		userText = fmt.Sprintf("<@%d>", user.ID)
	}
	contentText := content
	if contentText == "" {
		contentText = "*内容を取得できませんでした*"
	}
	return logMessage(colorMessageDelete,
		"### 🗑️ メッセージ削除",
		fmt.Sprintf("**ユーザー:** %s\n**チャンネル:** <#%d>\n**内容:**\n> %s",
			userText, channelID, contentText),
	)
}

func BuildMemberJoinLog(member discord.Member) discord.MessageCreate {
	createdAt := member.User.CreatedAt().Format("2006-01-02")
	return logMessage(colorMemberJoin,
		"### 📥 メンバー参加",
		fmt.Sprintf("**ユーザー:** <@%d> (%s)\n**アカウント作成:** %s",
			member.User.ID, member.User.Username, createdAt),
	)
}

func BuildMemberLeaveLog(user discord.User) discord.MessageCreate {
	return logMessage(colorMemberLeave,
		"### 📤 メンバー退出",
		fmt.Sprintf("**ユーザー:** <@%d> (%s)",
			user.ID, user.Username),
	)
}

func BuildBanLog(user discord.User) discord.MessageCreate {
	return logMessage(colorBan,
		"### 🔨 BAN",
		fmt.Sprintf("**ユーザー:** <@%d> (ID: %d)",
			user.ID, user.ID),
	)
}

func BuildUnbanLog(user discord.User) discord.MessageCreate {
	return logMessage(colorUnban,
		"### 🔓 BAN解除",
		fmt.Sprintf("**ユーザー:** <@%d> (ID: %d)",
			user.ID, user.ID),
	)
}

func BuildRoleCreateLog(role discord.Role) discord.MessageCreate {
	return buildRoleLog("作成", role)
}

func BuildRoleUpdateLog(role discord.Role) discord.MessageCreate {
	return buildRoleLog("更新", role)
}

func BuildRoleDeleteLog(role discord.Role) discord.MessageCreate {
	return buildRoleLog("削除", role)
}

func buildRoleLog(action string, role discord.Role) discord.MessageCreate {
	colorText := "なし"
	if role.Color != 0 {
		colorText = fmt.Sprintf("#%06X", role.Color)
	}
	return logMessage(colorRole,
		fmt.Sprintf("### 🏷️ ロール%s", action),
		fmt.Sprintf("**ロール:** %s\n**色:** %s",
			role.Name, colorText),
	)
}

func BuildChannelCreateLog(channel discord.GuildChannel) discord.MessageCreate {
	return buildChannelLog("作成", channel)
}

func BuildChannelUpdateLog(channel discord.GuildChannel) discord.MessageCreate {
	return buildChannelLog("更新", channel)
}

func BuildChannelDeleteLog(channel discord.GuildChannel) discord.MessageCreate {
	return buildChannelLog("削除", channel)
}

func buildChannelLog(action string, channel discord.GuildChannel) discord.MessageCreate {
	return logMessage(colorChannel,
		fmt.Sprintf("### 📁 チャンネル%s", action),
		fmt.Sprintf("**チャンネル:** %s\n**タイプ:** %s",
			channel.Name(), channelTypeName(channel.Type())),
	)
}

func channelTypeName(t discord.ChannelType) string {
	switch t {
	case discord.ChannelTypeGuildText:
		return "テキスト"
	case discord.ChannelTypeGuildVoice:
		return "ボイス"
	case discord.ChannelTypeGuildCategory:
		return "カテゴリ"
	case discord.ChannelTypeGuildNews:
		return "ニュース"
	case discord.ChannelTypeGuildStageVoice:
		return "ステージ"
	case discord.ChannelTypeGuildForum:
		return "フォーラム"
	default:
		return "その他"
	}
}

func logMessage(color int, title, body string) discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(title),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(body),
		).WithAccentColor(color),
	)
}
