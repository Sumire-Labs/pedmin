// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/ui"
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *BuilderHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch action {
	case "list":
		h.handleList(e)
	case "create_prompt":
		h.handleCreatePrompt(e)
	case "select":
		h.handleSelect(e, extra)
	case "add_text":
		h.handleAddTextPrompt(e, extra)
	case "add_section":
		h.handleAddSectionPrompt(e, extra)
	case "add_separator":
		h.handleAddSeparator(e, extra)
	case "sep_select":
		h.handleSepSelect(e, extra)
	case "add_media":
		h.handleAddMediaPrompt(e, extra)
	case "add_links":
		h.handleAddLinksPrompt(e, extra)
	case "manage":
		h.handleManage(e, extra)
	case "manage_select":
		h.handleManageSelect(e, extra)
	case "delete_comp":
		h.handleDeleteComponent(e, extra)
	case "move_up":
		h.handleMoveUp(e, extra)
	case "move_down":
		h.handleMoveDown(e, extra)
	case "preview":
		h.handlePreview(e, extra)
	case "deploy_prompt":
		h.handleDeployPrompt(e, extra)
	case "deploy_channel":
		h.handleDeployChannel(e, extra)
	case "deploy_confirm":
		h.handleDeployConfirm(e, extra)
	case "delete_panel":
		h.handleDeletePanel(e, extra)
	case "delete_confirm":
		h.handleDeleteConfirm(e, extra)
	case "rename":
		h.handleRenamePrompt(e, extra)
	case "back":
		h.handleList(e)
	}
}

// --- List ---

func (h *BuilderHandler) handleList(e *events.ComponentInteractionCreate) {
	panels, err := h.service.GetPanels(*e.GuildID())
	if err != nil {
		h.logger.Error("failed to get panels", slog.Any("error", err))
		return
	}
	_ = e.UpdateMessage(view.BuilderListPanelUpdate(panels, len(panels)))
}

func (h *BuilderHandler) handleCreatePrompt(e *events.ComponentInteractionCreate) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.BuilderModuleID + ":create_modal",
		Title:    "パネル作成",
		Components: []discord.LayoutComponent{
			discord.NewLabel("パネル名",
				discord.NewShortTextInput(model.BuilderModuleID+":panel_name").
					WithRequired(true).WithPlaceholder("ルール・ウェルカム等"),
			),
		},
	})
}

func (h *BuilderHandler) handleSelect(e *events.ComponentInteractionCreate, extra string) {
	panelID := extra
	// If from StringSelect, extract from data
	if panelID == "" {
		data, ok := e.Data.(discord.StringSelectMenuInteractionData)
		if !ok || len(data.Values) == 0 {
			return
		}
		panelID = data.Values[0]
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	panel, err := h.service.GetPanel(id, *e.GuildID())
	if err != nil {
		h.logger.Error("failed to get panel", slog.Any("error", err))
		return
	}

	_ = e.UpdateMessage(view.BuilderEditPanel(panel))
}

// --- Add Components ---

func (h *BuilderHandler) handleAddTextPrompt(e *events.ComponentInteractionCreate, panelID string) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.BuilderModuleID + ":text_modal:" + panelID,
		Title:    "テキスト追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("テキスト (Markdown対応)",
				discord.NewParagraphTextInput(model.BuilderModuleID+":text_content").
					WithRequired(true).WithPlaceholder("## タイトル\n本文テキスト..."),
			),
		},
	})
}

func (h *BuilderHandler) handleAddSectionPrompt(e *events.ComponentInteractionCreate, panelID string) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.BuilderModuleID + ":section_modal:" + panelID,
		Title:    "セクション追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("テキスト1",
				discord.NewShortTextInput(model.BuilderModuleID+":section_text1").
					WithRequired(true),
			),
			discord.NewLabel("テキスト2 (任意)",
				discord.NewShortTextInput(model.BuilderModuleID+":section_text2").
					WithRequired(false),
			),
			discord.NewLabel("テキスト3 (任意)",
				discord.NewShortTextInput(model.BuilderModuleID+":section_text3").
					WithRequired(false),
			),
			discord.NewLabel("サムネイルURL (任意)",
				discord.NewShortTextInput(model.BuilderModuleID+":section_thumb").
					WithRequired(false).WithPlaceholder("https://..."),
			),
		},
	})
}

func (h *BuilderHandler) handleAddSeparator(e *events.ComponentInteractionCreate, panelID string) {
	options := []discord.StringSelectMenuOption{
		{Label: "小さいセパレータ", Value: "small"},
		{Label: "大きいセパレータ", Value: "large"},
	}
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay("セパレータの種類を選択:"),
			discord.NewActionRow(
				discord.NewStringSelectMenu(model.BuilderModuleID+":sep_select:"+panelID, "種類を選択...", options...),
			),
		),
	}))
}

func (h *BuilderHandler) handleSepSelect(e *events.ComponentInteractionCreate, panelID string) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	comp := model.PanelComponent{
		Type:    model.PanelComponentSeparator,
		Spacing: data.Values[0],
	}

	panel, err := h.service.AddComponent(id, *e.GuildID(), comp)
	if err != nil {
		h.logger.Error("failed to add separator", slog.Any("error", err))
		return
	}

	_ = e.UpdateMessage(view.BuilderEditPanel(panel))
}

func (h *BuilderHandler) handleAddMediaPrompt(e *events.ComponentInteractionCreate, panelID string) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.BuilderModuleID + ":media_modal:" + panelID,
		Title:    "画像追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("画像URL 1",
				discord.NewShortTextInput(model.BuilderModuleID+":media_url1").
					WithRequired(true).WithPlaceholder("https://..."),
			),
			discord.NewLabel("説明 1 (任意)",
				discord.NewShortTextInput(model.BuilderModuleID+":media_desc1").
					WithRequired(false),
			),
			discord.NewLabel("画像URL 2 (任意)",
				discord.NewShortTextInput(model.BuilderModuleID+":media_url2").
					WithRequired(false).WithPlaceholder("https://..."),
			),
			discord.NewLabel("説明 2 (任意)",
				discord.NewShortTextInput(model.BuilderModuleID+":media_desc2").
					WithRequired(false),
			),
		},
	})
}

func (h *BuilderHandler) handleAddLinksPrompt(e *events.ComponentInteractionCreate, panelID string) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.BuilderModuleID + ":links_modal:" + panelID,
		Title:    "リンクボタン追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("ボタン (1行1個: ラベル|URL)",
				discord.NewParagraphTextInput(model.BuilderModuleID+":links_input").
					WithRequired(true).WithPlaceholder("公式サイト|https://example.com\nDiscord|https://discord.gg/..."),
			),
		},
	})
}

// --- Manage ---

func (h *BuilderHandler) handleManage(e *events.ComponentInteractionCreate, panelID string) {
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	panel, err := h.service.GetPanel(id, *e.GuildID())
	if err != nil {
		h.logger.Error("failed to get panel", slog.Any("error", err))
		return
	}

	_ = e.CreateMessage(view.BuilderManagePanel(panel))
}

func (h *BuilderHandler) handleManageSelect(e *events.ComponentInteractionCreate, panelID string) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	index, err := strconv.Atoi(data.Values[0])
	if err != nil {
		return
	}

	panel, err := h.service.GetPanel(id, *e.GuildID())
	if err != nil || index >= len(panel.Components) {
		return
	}

	_ = e.UpdateMessage(view.BuilderComponentDetail(panel, index))
}

func (h *BuilderHandler) handleDeleteComponent(e *events.ComponentInteractionCreate, extra string) {
	panelID, indexStr, _ := strings.Cut(extra, ":")
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return
	}

	panel, err := h.service.RemoveComponent(id, *e.GuildID(), index)
	if err != nil {
		h.logger.Error("failed to remove component", slog.Any("error", err))
		return
	}

	if len(panel.Components) == 0 {
		_ = e.UpdateMessage(view.BuilderEditPanel(panel))
		return
	}

	msg := view.BuilderManagePanel(panel)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2(msg.Components))
}

func (h *BuilderHandler) handleMoveUp(e *events.ComponentInteractionCreate, extra string) {
	h.handleMove(e, extra, -1)
}

func (h *BuilderHandler) handleMoveDown(e *events.ComponentInteractionCreate, extra string) {
	h.handleMove(e, extra, 1)
}

func (h *BuilderHandler) handleMove(e *events.ComponentInteractionCreate, extra string, direction int) {
	panelID, indexStr, _ := strings.Cut(extra, ":")
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return
	}

	panel, err := h.service.MoveComponent(id, *e.GuildID(), index, index+direction)
	if err != nil {
		h.logger.Error("failed to move component", slog.Any("error", err))
		return
	}

	_ = e.UpdateMessage(view.BuilderComponentDetail(panel, index+direction))
}

// --- Preview & Deploy ---

func (h *BuilderHandler) handlePreview(e *events.ComponentInteractionCreate, panelID string) {
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	panel, err := h.service.GetPanel(id, *e.GuildID())
	if err != nil {
		h.logger.Error("failed to get panel", slog.Any("error", err))
		return
	}

	_ = e.CreateMessage(h.service.PreviewPanel(panel))
}

func (h *BuilderHandler) handleDeployPrompt(e *events.ComponentInteractionCreate, panelID string) {
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}
	_ = e.CreateMessage(view.BuilderDeployPrompt(id))
}

func (h *BuilderHandler) handleDeployChannel(e *events.ComponentInteractionCreate, panelID string) {
	data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	_ = e.UpdateMessage(view.BuilderDeployConfirm(id, data.Values[0]))
}

func (h *BuilderHandler) handleDeployConfirm(e *events.ComponentInteractionCreate, extra string) {
	panelID, channelIDStr, _ := strings.Cut(extra, ":")
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}
	channelID, err := snowflake.Parse(channelIDStr)
	if err != nil {
		return
	}

	panel, err := h.service.GetPanel(id, *e.GuildID())
	if err != nil {
		h.logger.Error("failed to get panel", slog.Any("error", err))
		return
	}

	if err := h.service.DeployPanel(panel, channelID); err != nil {
		h.logger.Error("failed to deploy panel", slog.Any("error", err))
		_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
			view.BuilderErrorContainer(fmt.Sprintf("配信に失敗しました: %s", err.Error())),
		}))
		return
	}

	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("**%s** を <#%d> に配信しました。", panel.Name, channelID)),
		),
	}))
}

// --- Delete & Rename ---

func (h *BuilderHandler) handleDeletePanel(e *events.ComponentInteractionCreate, panelID string) {
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	panel, err := h.service.GetPanel(id, *e.GuildID())
	if err != nil {
		return
	}

	_ = e.UpdateMessage(view.BuilderDeleteConfirm(id, panel.Name))
}

func (h *BuilderHandler) handleDeleteConfirm(e *events.ComponentInteractionCreate, panelID string) {
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	if err := h.service.DeletePanel(id, *e.GuildID()); err != nil {
		h.logger.Error("failed to delete panel", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("パネルの削除に失敗しました。"))
		return
	}

	h.handleList(e)
}

func (h *BuilderHandler) handleRenamePrompt(e *events.ComponentInteractionCreate, panelID string) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.BuilderModuleID + ":rename_modal:" + panelID,
		Title:    "パネル名変更",
		Components: []discord.LayoutComponent{
			discord.NewLabel("新しいパネル名",
				discord.NewShortTextInput(model.BuilderModuleID+":new_name").
					WithRequired(true),
			),
		},
	})
}
