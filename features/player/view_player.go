package player

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

const accentPlaying = 0x00B894
const accentPaused = 0xFDCB6E
const accentIdle = 0x636E72

func BuildPlayerUI(player disgolink.Player, queue *Queue) discord.ContainerComponent {
	track := player.Track()
	paused := player.Paused()

	if track == nil {
		return buildIdleUI(queue)
	}

	accentColor := accentPlaying
	if paused {
		accentColor = accentPaused
	}

	info := track.Info

	statusIcon := "▶️"
	if paused {
		statusIcon = "⏸️"
	}

	components := []discord.ContainerSubComponent{
		discord.NewSection(
			discord.NewTextDisplay(fmt.Sprintf("### %s 再生中", statusIcon)),
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
		buildControlRow(paused),
		buildSecondaryRow(queue.LoopMode(), player.Volume()),
	}

	return discord.NewContainer(components...).WithAccentColor(accentColor)
}

func buildIdleUI(queue *Queue) discord.ContainerComponent {
	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay("### Pedmin Player"),
		discord.NewTextDisplay("再生中の曲はありません。ボタンから曲を追加してください！"),
		discord.NewLargeSeparator(),
		buildControlRow(false),
		buildSecondaryRow(queue.LoopMode(), 100),
	}
	return discord.NewContainer(components...).WithAccentColor(accentIdle)
}

func buildControlRow(paused bool) discord.ActionRowComponent {
	playPauseLabel := "⏸"
	playPauseID := ModuleID + ":pause"
	if paused {
		playPauseLabel = "▶"
		playPauseID = ModuleID + ":resume"
	}

	return discord.NewActionRow(
		discord.NewPrimaryButton(playPauseLabel, playPauseID),
		discord.NewSecondaryButton("⏮", ModuleID+":previous"),
		discord.NewSecondaryButton("⏭", ModuleID+":skip"),
		discord.NewDangerButton("⏹", ModuleID+":stop"),
		discord.NewDangerButton("🔌", ModuleID+":disconnect"),
	)
}

func buildSecondaryRow(loopMode LoopMode, volume int) discord.ActionRowComponent {
	loopLabel := "🔁"
	switch loopMode {
	case LoopTrack:
		loopLabel = "🔂"
	case LoopQueue:
		loopLabel = "🔁"
	default:
		loopLabel = "➡️"
	}

	return discord.NewActionRow(
		discord.NewSecondaryButton(loopLabel, ModuleID+":loop"),
		discord.NewSecondaryButton("🔉", ModuleID+":vol_down"),
		discord.NewSecondaryButton(fmt.Sprintf("🔊 %d%%", volume), ModuleID+":vol_up"),
		discord.NewSuccessButton("➕ 追加", ModuleID+":add"),
		discord.NewSecondaryButton("📜 キュー", ModuleID+":queue"),
	)
}
