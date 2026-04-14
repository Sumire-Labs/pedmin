// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/module"
	"github.com/Sumire-Labs/pedmin/internal/ui"
)

// ModeratorHandler implements module.Module for moderation commands.
type ModeratorHandler struct {
	logger *slog.Logger
}

// NewModeratorHandler creates a new ModeratorHandler.
func NewModeratorHandler(logger *slog.Logger) *ModeratorHandler {
	return &ModeratorHandler{logger: logger}
}

func (h *ModeratorHandler) Info() module.Info {
	return module.Info{
		ID:          model.ModeratorModuleID,
		Name:        "モデレーション",
		Description: "キック / BAN / タイムアウトコマンド",
		AlwaysOn:    true,
	}
}

func (h *ModeratorHandler) Commands() []discord.ApplicationCommandCreate {
	kickPerm := discord.PermissionKickMembers
	banPerm := discord.PermissionBanMembers
	modPerm := discord.PermissionModerateMembers
	minDelete, maxDelete := 0, 7

	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:                     "kick",
			Description:              "ユーザーをサーバーからキックする",
			DefaultMemberPermissions: omit.New(&kickPerm),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "対象ユーザー",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "reason",
					Description: "理由（監査ログに記録）",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "ban",
			Description:              "ユーザーをサーバーからBANする",
			DefaultMemberPermissions: omit.New(&banPerm),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "対象ユーザー",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "delete_days",
					Description: "削除するメッセージ履歴（0〜7日）",
					Required:    false,
					MinValue:    &minDelete,
					MaxValue:    &maxDelete,
				},
				discord.ApplicationCommandOptionString{
					Name:        "reason",
					Description: "理由（監査ログに記録）",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "to",
			Description:              "ユーザーをタイムアウトする",
			DefaultMemberPermissions: omit.New(&modPerm),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "対象ユーザー",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "duration",
					Description: "タイムアウト期間",
					Required:    true,
					Choices: []discord.ApplicationCommandOptionChoiceInt{
						{Name: "60秒", Value: 60},
						{Name: "5分", Value: 300},
						{Name: "10分", Value: 600},
						{Name: "1時間", Value: 3600},
						{Name: "1日", Value: 86400},
						{Name: "1週間", Value: 604800},
					},
				},
				discord.ApplicationCommandOptionString{
					Name:        "reason",
					Description: "理由（監査ログに記録）",
					Required:    false,
				},
			},
		},
	}
}

func (h *ModeratorHandler) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	switch e.Data.CommandName() {
	case "kick":
		h.handleKick(e)
	case "ban":
		h.handleBan(e)
	case "to":
		h.handleTo(e)
	}
}

func (h *ModeratorHandler) handleKick(e *events.ApplicationCommandInteractionCreate) {
	guildID, target, reason, ok := h.extractCommon(e)
	if !ok {
		return
	}

	opts := buildReasonOpts(reason)
	if err := e.Client().Rest.RemoveMember(guildID, target.ID, opts...); err != nil {
		h.logger.Warn("failed to kick member",
			slog.String("user_id", target.ID.String()),
			slog.Any("error", err),
		)
		_ = e.CreateMessage(ui.EphemeralError(fmt.Sprintf("キックに失敗しました: %v", err)))
		return
	}

	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("✅ **%s** をキックしました。%s",
				target.Username, formatReason(reason))),
		),
	))
}

func (h *ModeratorHandler) handleBan(e *events.ApplicationCommandInteractionCreate) {
	guildID, target, reason, ok := h.extractCommon(e)
	if !ok {
		return
	}

	data := e.SlashCommandInteractionData()
	deleteDays := 0
	if v, ok := data.OptInt("delete_days"); ok {
		deleteDays = v
	}
	deleteDuration := time.Duration(deleteDays) * 24 * time.Hour

	opts := buildReasonOpts(reason)
	if err := e.Client().Rest.AddBan(guildID, target.ID, deleteDuration, opts...); err != nil {
		h.logger.Warn("failed to ban member",
			slog.String("user_id", target.ID.String()),
			slog.Any("error", err),
		)
		_ = e.CreateMessage(ui.EphemeralError(fmt.Sprintf("BANに失敗しました: %v", err)))
		return
	}

	msg := fmt.Sprintf("✅ **%s** をBANしました。", target.Username)
	if deleteDays > 0 {
		msg += fmt.Sprintf(" %d日分のメッセージを削除。", deleteDays)
	}
	msg += formatReason(reason)

	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(discord.NewTextDisplay(msg)),
	))
}

func (h *ModeratorHandler) handleTo(e *events.ApplicationCommandInteractionCreate) {
	guildID, target, reason, ok := h.extractCommon(e)
	if !ok {
		return
	}

	data := e.SlashCommandInteractionData()
	durationSec, ok := data.OptInt("duration")
	if !ok {
		_ = e.CreateMessage(ui.EphemeralError("期間が指定されていません。"))
		return
	}
	until := time.Now().Add(time.Duration(durationSec) * time.Second)

	update := discord.MemberUpdate{
		CommunicationDisabledUntil: omit.New(&until),
	}
	opts := buildReasonOpts(reason)
	if _, err := e.Client().Rest.UpdateMember(guildID, target.ID, update, opts...); err != nil {
		h.logger.Warn("failed to timeout member",
			slog.String("user_id", target.ID.String()),
			slog.Any("error", err),
		)
		_ = e.CreateMessage(ui.EphemeralError(fmt.Sprintf("タイムアウトに失敗しました: %v", err)))
		return
	}

	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("✅ **%s** をタイムアウトしました（<t:%d:R>まで）。%s",
				target.Username, until.Unix(), formatReason(reason))),
		),
	))
}

// extractCommon pulls guildID, target user, and reason from the interaction,
// and runs the shared safety checks. Returns ok=false if the command should abort.
func (h *ModeratorHandler) extractCommon(e *events.ApplicationCommandInteractionCreate) (snowflake.ID, discord.User, string, bool) {
	guildID := e.GuildID()
	if guildID == nil {
		_ = e.CreateMessage(ui.ErrorMessage("サーバー内でのみ使用できます。"))
		return 0, discord.User{}, "", false
	}

	data := e.SlashCommandInteractionData()
	target, ok := data.OptUser("user")
	if !ok {
		_ = e.CreateMessage(ui.EphemeralError("対象ユーザーが指定されていません。"))
		return 0, discord.User{}, "", false
	}

	if target.ID == e.User().ID {
		_ = e.CreateMessage(ui.EphemeralError("自分自身は対象にできません。"))
		return 0, discord.User{}, "", false
	}
	if target.ID == e.Client().ID() {
		_ = e.CreateMessage(ui.EphemeralError("Bot自身は対象にできません。"))
		return 0, discord.User{}, "", false
	}

	reason, _ := data.OptString("reason")
	return *guildID, target, reason, true
}

func buildReasonOpts(reason string) []rest.RequestOpt {
	if reason == "" {
		return nil
	}
	return []rest.RequestOpt{rest.WithReason(reason)}
}

func formatReason(reason string) string {
	if reason == "" {
		return ""
	}
	return fmt.Sprintf("\n**理由:** %s", reason)
}

func (h *ModeratorHandler) HandleComponent(_ *events.ComponentInteractionCreate) {}
func (h *ModeratorHandler) HandleModal(_ *events.ModalSubmitInteractionCreate)   {}
func (h *ModeratorHandler) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}
