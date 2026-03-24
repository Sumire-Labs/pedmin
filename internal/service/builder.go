// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"fmt"
	"log/slog"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"

	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/repository"
)

// BuilderService handles component panel business logic.
type BuilderService struct {
	store  repository.GuildStore
	client *disgobot.Client
	logger *slog.Logger
}

// NewBuilderService creates a new BuilderService.
func NewBuilderService(store repository.GuildStore, client *disgobot.Client, logger *slog.Logger) *BuilderService {
	return &BuilderService{
		store:  store,
		client: client,
		logger: logger,
	}
}

// CreatePanel creates a new empty panel.
func (s *BuilderService) CreatePanel(guildID snowflake.ID, name string) (*model.ComponentPanel, error) {
	panel := &model.ComponentPanel{
		GuildID:    guildID,
		Name:       name,
		Components: []model.PanelComponent{},
	}
	if err := s.store.CreatePanel(panel); err != nil {
		return nil, err
	}
	return panel, nil
}

// GetPanels returns all panels for a guild.
func (s *BuilderService) GetPanels(guildID snowflake.ID) ([]model.ComponentPanel, error) {
	return s.store.GetPanels(guildID)
}

// GetPanel returns a single panel.
func (s *BuilderService) GetPanel(id int64, guildID snowflake.ID) (*model.ComponentPanel, error) {
	return s.store.GetPanel(id, guildID)
}

// DeletePanel deletes a panel.
func (s *BuilderService) DeletePanel(id int64, guildID snowflake.ID) error {
	return s.store.DeletePanel(id, guildID)
}

// CountPanels returns the number of panels for a guild.
func (s *BuilderService) CountPanels(guildID snowflake.ID) (int, error) {
	return s.store.CountPanels(guildID)
}

// RenamePanel renames a panel.
func (s *BuilderService) RenamePanel(id int64, guildID snowflake.ID, newName string) (*model.ComponentPanel, error) {
	panel, err := s.store.GetPanel(id, guildID)
	if err != nil {
		return nil, err
	}
	panel.Name = newName
	if err := s.store.UpdatePanel(panel); err != nil {
		return nil, err
	}
	return panel, nil
}

// AddComponent appends a component to a panel.
func (s *BuilderService) AddComponent(id int64, guildID snowflake.ID, comp model.PanelComponent) (*model.ComponentPanel, error) {
	panel, err := s.store.GetPanel(id, guildID)
	if err != nil {
		return nil, err
	}
	if len(panel.Components) >= model.MaxComponentsPerPanel {
		return nil, fmt.Errorf("コンポーネント数が上限(%d)に達しています", model.MaxComponentsPerPanel)
	}
	panel.Components = append(panel.Components, comp)
	if err := s.store.UpdatePanel(panel); err != nil {
		return nil, err
	}
	return panel, nil
}

// RemoveComponent removes a component at the given index.
func (s *BuilderService) RemoveComponent(id int64, guildID snowflake.ID, index int) (*model.ComponentPanel, error) {
	panel, err := s.store.GetPanel(id, guildID)
	if err != nil {
		return nil, err
	}
	if index < 0 || index >= len(panel.Components) {
		return nil, fmt.Errorf("無効なインデックスです")
	}
	panel.Components = append(panel.Components[:index], panel.Components[index+1:]...)
	if err := s.store.UpdatePanel(panel); err != nil {
		return nil, err
	}
	return panel, nil
}

// MoveComponent moves a component from one index to another.
func (s *BuilderService) MoveComponent(id int64, guildID snowflake.ID, from, to int) (*model.ComponentPanel, error) {
	panel, err := s.store.GetPanel(id, guildID)
	if err != nil {
		return nil, err
	}
	if from < 0 || from >= len(panel.Components) || to < 0 || to >= len(panel.Components) {
		return nil, fmt.Errorf("無効なインデックスです")
	}
	comp := panel.Components[from]
	panel.Components = append(panel.Components[:from], panel.Components[from+1:]...)
	// Insert at 'to' position
	panel.Components = append(panel.Components[:to], append([]model.PanelComponent{comp}, panel.Components[to:]...)...)
	if err := s.store.UpdatePanel(panel); err != nil {
		return nil, err
	}
	return panel, nil
}

// RenderPanel converts a panel's components to a Discord container.
func (s *BuilderService) RenderPanel(panel *model.ComponentPanel) discord.ContainerComponent {
	var subs []discord.ContainerSubComponent

	for _, comp := range panel.Components {
		switch comp.Type {
		case model.PanelComponentText:
			subs = append(subs, discord.NewTextDisplay(comp.Content))

		case model.PanelComponentSection:
			var textDisplays []discord.SectionSubComponent
			for _, t := range comp.Texts {
				if t != "" {
					textDisplays = append(textDisplays, discord.NewTextDisplay(t))
				}
			}
			if len(textDisplays) == 0 {
				continue
			}
			section := discord.NewSection(textDisplays...)
			if comp.ThumbnailURL != "" {
				section = section.WithAccessory(discord.NewThumbnail(comp.ThumbnailURL))
			}
			subs = append(subs, section)

		case model.PanelComponentSeparator:
			if comp.Spacing == "large" {
				subs = append(subs, discord.NewLargeSeparator())
			} else {
				subs = append(subs, discord.NewSmallSeparator())
			}

		case model.PanelComponentMedia:
			var items []discord.MediaGalleryItem
			for _, item := range comp.Items {
				mgItem := discord.MediaGalleryItem{
					Media: discord.UnfurledMediaItem{URL: item.URL},
				}
				if item.Description != "" {
					mgItem.Description = item.Description
				}
				items = append(items, mgItem)
			}
			if len(items) > 0 {
				subs = append(subs, discord.NewMediaGallery(items...))
			}

		case model.PanelComponentLinks:
			var buttons []discord.InteractiveComponent
			for _, btn := range comp.Buttons {
				b := discord.NewLinkButton(btn.Label, btn.URL)
				if btn.Emoji != "" {
					b = b.WithEmoji(discord.ComponentEmoji{Name: btn.Emoji})
				}
				buttons = append(buttons, b)
			}
			if len(buttons) > 0 {
				subs = append(subs, discord.NewActionRow(buttons...))
			}
		}
	}

	if len(subs) == 0 {
		subs = append(subs, discord.NewTextDisplay("-# パネルにコンポーネントがありません"))
	}

	return discord.NewContainer(subs...)
}

// PreviewPanel renders the panel as an ephemeral message.
func (s *BuilderService) PreviewPanel(panel *model.ComponentPanel) discord.MessageCreate {
	container := s.RenderPanel(panel)
	return discord.NewMessageCreateV2(container).WithEphemeral(true)
}

// DeployPanel sends the rendered panel to a channel.
func (s *BuilderService) DeployPanel(panel *model.ComponentPanel, channelID snowflake.ID) error {
	container := s.RenderPanel(panel)
	msg := discord.NewMessageCreateV2(container)
	_, err := s.client.Rest.CreateMessage(channelID, msg)
	return err
}
