// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

import "github.com/disgoorg/snowflake/v2"

// AutoroleSettings holds per-guild autorole configuration (stored as JSON).
type AutoroleSettings struct {
	UserRoleID snowflake.ID `json:"user_role_id"`
	BotRoleID  snowflake.ID `json:"bot_role_id"`
}

// DefaultAutoroleSettings returns the default settings.
func DefaultAutoroleSettings() *AutoroleSettings {
	return &AutoroleSettings{}
}
