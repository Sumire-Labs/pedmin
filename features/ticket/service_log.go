package ticket

import (
	"log/slog"

	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/store"
)

func (t *Ticket) sendTicketLog(guildID snowflake.ID, ticket *store.Ticket) {
	settings, err := LoadSettings(t.store, guildID)
	if err != nil || settings.LogChannelID == 0 {
		return
	}

	// Reload to get closed_at/closed_by
	fresh, err := t.store.GetTicketByChannel(ticket.ChannelID)
	if err != nil || fresh == nil {
		fresh = ticket
	}

	msg := BuildTicketLog(fresh)
	if _, err := t.client.Rest.CreateMessage(settings.LogChannelID, msg); err != nil {
		t.logger.Error("failed to send ticket log", slog.Any("error", err))
	}
}

func (t *Ticket) sendTranscriptLog(guildID snowflake.ID, ticket *store.Ticket) {
	settings, err := LoadSettings(t.store, guildID)
	if err != nil || settings.LogChannelID == 0 {
		return
	}

	// Reload to get closed_at/closed_by
	fresh, err := t.store.GetTicketByChannel(ticket.ChannelID)
	if err != nil || fresh == nil {
		fresh = ticket
	}

	file, err := t.generateTranscript(guildID, fresh)
	if err != nil {
		t.logger.Error("failed to generate transcript", slog.Any("error", err))
		return
	}

	msg := BuildTicketLog(fresh).AddFiles(file)
	if _, err := t.client.Rest.CreateMessage(settings.LogChannelID, msg); err != nil {
		t.logger.Error("failed to send transcript log", slog.Any("error", err))
	}
}
