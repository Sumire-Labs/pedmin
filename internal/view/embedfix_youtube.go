// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/pkg/deepl"
)

// BuildYouTubeEmbed builds a YouTube video embed message.
func BuildYouTubeEmbed(video *model.YouTubeVideo, ref model.EmbedRef) discord.MessageCreate {
	components := BuildYouTubeComponents(video, ref, "", "")
	return discord.NewMessageCreateV2(discord.NewContainer(components...))
}

// BuildYouTubeEmbedTranslated builds a translated YouTube video embed as layout components.
func BuildYouTubeEmbedTranslated(video *model.YouTubeVideo, result *deepl.TranslateResult, ref model.EmbedRef) []discord.LayoutComponent {
	duration := formatYouTubeDuration(video)
	footer := fmt.Sprintf("%s | %s | <t:%d:R> · %sから翻訳", emojiYouTube, duration, video.PublishedAt.Unix(), deepl.LangName(result.DetectedLanguage))
	components := BuildYouTubeComponents(video, ref, result.TranslatedText, footer)
	return []discord.LayoutComponent{discord.NewContainer(components...)}
}

// BuildYouTubeComponents builds the YouTube video embed sub-components.
func BuildYouTubeComponents(video *model.YouTubeVideo, ref model.EmbedRef, translatedText, footerOverride string) []discord.ContainerSubComponent {
	authorLine := fmt.Sprintf("**%s**", video.Author)

	var headerComponent discord.ContainerSubComponent
	if video.AuthorAvatar != "" {
		headerComponent = discord.NewSection(
			discord.NewTextDisplay(authorLine),
		).WithAccessory(discord.NewThumbnail(video.AuthorAvatar))
	} else {
		headerComponent = discord.NewTextDisplay(authorLine)
	}

	components := []discord.ContainerSubComponent{
		headerComponent,
		discord.NewSmallSeparator(),
	}

	if translatedText != "" {
		components = append(components, discord.NewTextDisplay(translatedText))
	} else if video.Title != "" {
		components = append(components, discord.NewTextDisplay(video.Title))
	}

	if video.ThumbnailURL != "" {
		components = append(components, discord.NewMediaGallery(
			discord.MediaGalleryItem{
				Media: discord.UnfurledMediaItem{URL: video.ThumbnailURL},
			},
		))
	}

	components = append(components, discord.NewSmallSeparator())

	stats := fmt.Sprintf("%s %s  %s %s",
		emojiPlay, FormatCount(video.ViewCount),
		emojiGood, FormatCount(video.LikeCount),
	)
	components = append(components, discord.NewTextDisplay(stats))

	duration := formatYouTubeDuration(video)
	footer := fmt.Sprintf("%s | %s | <t:%d:R>", emojiYouTube, duration, video.PublishedAt.Unix())
	if footerOverride != "" {
		footer = footerOverride
	}
	components = append(components, discord.NewTextDisplay(footer))

	if video.Title != "" && footerOverride == "" {
		customID := fmt.Sprintf("%s:translate:%s:%s", model.EmbedFixModuleID, model.PlatformYouTube, ref.Params[0])
		components = append(components,
			discord.NewActionRow(
				discord.NewSecondaryButton("\U0001f310 翻訳", customID),
			),
		)
	}

	return components
}

func formatYouTubeDuration(video *model.YouTubeVideo) string {
	if video.IsLive {
		return "LIVE"
	}
	return FormatSeconds(video.Duration)
}
