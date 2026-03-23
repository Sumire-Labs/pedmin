package player

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

func BuildPlayerUI(player disgolink.Player, queue *Queue) discord.ContainerComponent {
	track := player.Track()
	if track == nil {
		return buildIdleUI(queue)
	}

	info := track.Info

	components := []discord.ContainerSubComponent{
		discord.NewSection(
			discord.NewTextDisplay("### ▶️ 再生中"),
			discord.NewTextDisplay(fmt.Sprintf("**%s**\nby %s", info.Title, info.Author)),
		).WithAccessory(buildThumbnail(info)),
		discord.NewLargeSeparator(),
		discord.NewTextDisplay(fmt.Sprintf(
			"%s  %s / %s  |  %s",
			buildProgressBar(player.Position(), info.Length),
			formatDuration(player.Position()),
			formatDuration(info.Length),
			queue.LoopMode().String(),
		)),
		discord.NewLargeSeparator(),
	}
	components = append(components, buildButtonRows(queue.LoopMode())...)
	return discord.NewContainer(components...)
}

func buildIdleUI(queue *Queue) discord.ContainerComponent {
	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay("### Pedmin Player"),
		discord.NewTextDisplay("再生中の曲はありません。ボタンから曲を追加してください！"),
		discord.NewLargeSeparator(),
	}
	components = append(components, buildButtonRows(queue.LoopMode())...)
	return discord.NewContainer(components...)
}

func buildButtonRows(loopMode LoopMode) []discord.ContainerSubComponent {
	modeLabel := "モード"
	switch loopMode {
	case LoopTrack:
		modeLabel = "モード: トラック"
	case LoopQueue:
		modeLabel = "モード: キュー"
	}

	return []discord.ContainerSubComponent{
		discord.NewActionRow(
			discord.NewSecondaryButton("⏪", ModuleID+":seek_back"),
			discord.NewSecondaryButton("スキップ", ModuleID+":skip"),
			discord.NewDangerButton("停止", ModuleID+":stop"),
			discord.NewSecondaryButton("⏩", ModuleID+":seek_forward"),
			discord.NewSuccessButton("追加", ModuleID+":add"),
		),
		discord.NewActionRow(
			discord.NewSecondaryButton("🔀 シャッフル", ModuleID+":shuffle"),
			discord.NewSecondaryButton("キュー", ModuleID+":queue"),
			discord.NewSecondaryButton(modeLabel, ModuleID+":loop"),
		),
	}
}
