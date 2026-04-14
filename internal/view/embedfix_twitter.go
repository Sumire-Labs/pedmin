// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/pkg/deepl"
)

// BuildTweetEmbed builds a tweet embed message.
func BuildTweetEmbed(tweet *model.Tweet, ref model.EmbedRef) discord.MessageCreate {
	components := BuildTweetComponents(tweet, ref, tweet.Text, "", "")
	return discord.NewMessageCreateV2(discord.NewContainer(components...))
}

// BuildTweetEmbedTranslated builds a translated tweet embed as layout components.
// quoteText is the translated quoted tweet text (empty if no quote or not translated).
func BuildTweetEmbedTranslated(tweet *model.Tweet, result *deepl.TranslateResult, quoteText string, ref model.EmbedRef) []discord.LayoutComponent {
	footer := fmt.Sprintf("%s | <t:%d:R> · %sから翻訳", emojiX, tweet.CreatedAt.Unix(), deepl.LangName(result.DetectedLanguage))
	components := BuildTweetComponents(tweet, ref, result.TranslatedText, quoteText, footer)
	return []discord.LayoutComponent{discord.NewContainer(components...)}
}

// BuildTweetComponents builds the tweet embed sub-components.
// If quoteText is non-empty, it replaces the original quoted tweet text.
func BuildTweetComponents(tweet *model.Tweet, ref model.EmbedRef, text, quoteText, footerOverride string) []discord.ContainerSubComponent {
	components := []discord.ContainerSubComponent{
		discord.NewSection(
			discord.NewTextDisplay(fmt.Sprintf("**%s** [@%s](https://x.com/%s)", tweet.Author.Name, tweet.Author.ScreenName, tweet.Author.ScreenName)),
		).WithAccessory(discord.NewThumbnail(tweet.Author.AvatarURL)),
		discord.NewSmallSeparator(),
	}

	if text != "" {
		components = append(components, discord.NewTextDisplay(text))
	}

	if len(tweet.Media) > 0 {
		items := make([]discord.MediaGalleryItem, 0, len(tweet.Media))
		for _, m := range tweet.Media {
			items = append(items, discord.MediaGalleryItem{
				Media: discord.UnfurledMediaItem{URL: m.URL},
			})
		}
		components = append(components, discord.NewMediaGallery(items...))
	}

	if tweet.Quote != nil {
		components = append(components,
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(fmt.Sprintf("%s 引用", emojiRepost)),
			discord.NewTextDisplay(fmt.Sprintf("**%s** @%s", tweet.Quote.Author.Name, tweet.Quote.Author.ScreenName)),
		)
		displayQuoteText := tweet.Quote.Text
		if quoteText != "" {
			displayQuoteText = quoteText
		}
		if displayQuoteText != "" {
			components = append(components, discord.NewTextDisplay(displayQuoteText))
		}
		if len(tweet.Quote.Media) > 0 {
			items := make([]discord.MediaGalleryItem, 0, len(tweet.Quote.Media))
			for _, m := range tweet.Quote.Media {
				items = append(items, discord.MediaGalleryItem{
					Media: discord.UnfurledMediaItem{URL: m.URL},
				})
			}
			components = append(components, discord.NewMediaGallery(items...))
		}
	}

	components = append(components, discord.NewSmallSeparator())

	stats := fmt.Sprintf("%s %s  %s %s  %s %s  %s %s",
		emojiMessages, FormatCount(tweet.Replies),
		emojiRepost, FormatCount(tweet.Retweets),
		emojiLike, FormatCount(tweet.Likes),
		emojiGraph, FormatCount(tweet.Views),
	)
	components = append(components, discord.NewTextDisplay(stats))

	footer := fmt.Sprintf("%s | <t:%d:R>", emojiX, tweet.CreatedAt.Unix())
	if footerOverride != "" {
		footer = footerOverride
	}
	components = append(components, discord.NewTextDisplay(footer))

	// Show translate button only for non-Japanese tweets and when not already translated
	if tweet.Lang != "ja" && footerOverride == "" {
		customID := fmt.Sprintf("%s:translate:%s:%s:%s", model.EmbedFixModuleID, model.PlatformTwitter, ref.Params[0], ref.Params[1])
		components = append(components,
			discord.NewActionRow(
				discord.NewSecondaryButton("\U0001f310 翻訳", customID),
			),
		)
	}

	// Show revert button when translated
	if footerOverride != "" {
		revertID := fmt.Sprintf("%s:revert:%s:%s:%s", model.EmbedFixModuleID, model.PlatformTwitter, ref.Params[0], ref.Params[1])
		components = append(components,
			discord.NewActionRow(
				discord.NewSecondaryButton("↩ 原文に戻す", revertID),
			),
		)
	}

	return components
}

// BuildTweetEmbedOriginal builds the original tweet embed as layout components (for revert).
func BuildTweetEmbedOriginal(tweet *model.Tweet, ref model.EmbedRef) []discord.LayoutComponent {
	components := BuildTweetComponents(tweet, ref, tweet.Text, "", "")
	return []discord.LayoutComponent{discord.NewContainer(components...)}
}
