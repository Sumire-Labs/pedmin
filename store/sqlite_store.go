package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disgoorg/snowflake/v2"
	_ "modernc.org/sqlite"
)

var migrations = []struct {
	version int
	sql     string
}{
	{
		version: 1,
		sql: `
			CREATE TABLE IF NOT EXISTS guild_modules (
				guild_id   INTEGER NOT NULL,
				module_id  TEXT    NOT NULL,
				enabled    BOOLEAN NOT NULL DEFAULT 0,
				PRIMARY KEY (guild_id, module_id)
			);

			CREATE TABLE IF NOT EXISTS guild_module_settings (
				guild_id   INTEGER NOT NULL,
				module_id  TEXT    NOT NULL,
				settings   TEXT    NOT NULL DEFAULT '{}',
				PRIMARY KEY (guild_id, module_id)
			);
		`,
	},
	{
		version: 2,
		sql: `
			CREATE TABLE IF NOT EXISTS tickets (
				guild_id   INTEGER NOT NULL,
				number     INTEGER NOT NULL,
				channel_id INTEGER NOT NULL,
				user_id    INTEGER NOT NULL,
				subject    TEXT    NOT NULL DEFAULT '',
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				closed_at  TIMESTAMP,
				closed_by  INTEGER,
				PRIMARY KEY (guild_id, number)
			);
			CREATE INDEX IF NOT EXISTS idx_tickets_channel ON tickets (channel_id);
		`,
	},
	{
		version: 3,
		sql: `
			CREATE TABLE IF NOT EXISTS rss_feeds (
				id         INTEGER PRIMARY KEY AUTOINCREMENT,
				guild_id   INTEGER NOT NULL,
				url        TEXT    NOT NULL,
				channel_id INTEGER NOT NULL,
				title      TEXT    NOT NULL DEFAULT '',
				added_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(guild_id, url)
			);
			CREATE INDEX IF NOT EXISTS idx_rss_feeds_guild ON rss_feeds(guild_id);

			CREATE TABLE IF NOT EXISTS rss_seen_items (
				feed_id    INTEGER NOT NULL REFERENCES rss_feeds(id) ON DELETE CASCADE,
				item_hash  TEXT    NOT NULL,
				seen_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				PRIMARY KEY (feed_id, item_hash)
			);
		`,
	},
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(wal)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	s := &SQLiteStore{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return s, nil
}

func (s *SQLiteStore) migrate() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version    INTEGER PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	for _, m := range migrations {
		var exists int
		err := s.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", m.version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration %d: %w", m.version, err)
		}
		if exists > 0 {
			continue
		}

		tx, err := s.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", m.version, err)
		}

		if _, err := tx.Exec(m.sql); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %w", m.version, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", m.version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", m.version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", m.version, err)
		}
	}
	return nil
}

func (s *SQLiteStore) Get(guildID snowflake.ID) (*GuildSettings, error) {
	settings := &GuildSettings{
		GuildID:        guildID,
		EnabledModules: make(map[string]bool),
		ModuleSettings: make(map[string]any),
	}

	// Load enabled modules
	rows, err := s.db.Query("SELECT module_id, enabled FROM guild_modules WHERE guild_id = ?", int64(guildID))
	if err != nil {
		return nil, fmt.Errorf("failed to query guild modules: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var moduleID string
		var enabled bool
		if err := rows.Scan(&moduleID, &enabled); err != nil {
			return nil, fmt.Errorf("failed to scan guild module: %w", err)
		}
		settings.EnabledModules[moduleID] = enabled
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate guild modules: %w", err)
	}

	// Load module settings
	settingsRows, err := s.db.Query("SELECT module_id, settings FROM guild_module_settings WHERE guild_id = ?", int64(guildID))
	if err != nil {
		return nil, fmt.Errorf("failed to query guild module settings: %w", err)
	}
	defer settingsRows.Close()

	for settingsRows.Next() {
		var moduleID string
		var settingsJSON string
		if err := settingsRows.Scan(&moduleID, &settingsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan guild module settings: %w", err)
		}
		var parsed any
		if err := json.Unmarshal([]byte(settingsJSON), &parsed); err != nil {
			return nil, fmt.Errorf("failed to parse module settings for %s: %w", moduleID, err)
		}
		settings.ModuleSettings[moduleID] = parsed
	}
	if err := settingsRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate guild module settings: %w", err)
	}

	return settings, nil
}

func (s *SQLiteStore) Save(settings *GuildSettings) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	gid := int64(settings.GuildID)

	// Upsert enabled modules
	for moduleID, enabled := range settings.EnabledModules {
		_, err := tx.Exec(
			"INSERT INTO guild_modules (guild_id, module_id, enabled) VALUES (?, ?, ?) ON CONFLICT(guild_id, module_id) DO UPDATE SET enabled = excluded.enabled",
			gid, moduleID, enabled,
		)
		if err != nil {
			return fmt.Errorf("failed to upsert module %s: %w", moduleID, err)
		}
	}

	// Upsert module settings
	for moduleID, modSettings := range settings.ModuleSettings {
		settingsJSON, err := json.Marshal(modSettings)
		if err != nil {
			return fmt.Errorf("failed to marshal settings for %s: %w", moduleID, err)
		}
		_, err = tx.Exec(
			"INSERT INTO guild_module_settings (guild_id, module_id, settings) VALUES (?, ?, ?) ON CONFLICT(guild_id, module_id) DO UPDATE SET settings = excluded.settings",
			gid, moduleID, string(settingsJSON),
		)
		if err != nil {
			return fmt.Errorf("failed to upsert settings for %s: %w", moduleID, err)
		}
	}

	return tx.Commit()
}

func (s *SQLiteStore) IsModuleEnabled(guildID snowflake.ID, moduleID string) (bool, error) {
	var enabled bool
	err := s.db.QueryRow(
		"SELECT enabled FROM guild_modules WHERE guild_id = ? AND module_id = ?",
		int64(guildID), moduleID,
	).Scan(&enabled)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check module enabled: %w", err)
	}
	return enabled, nil
}

func (s *SQLiteStore) SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error {
	_, err := s.db.Exec(
		"INSERT INTO guild_modules (guild_id, module_id, enabled) VALUES (?, ?, ?) ON CONFLICT(guild_id, module_id) DO UPDATE SET enabled = excluded.enabled",
		int64(guildID), moduleID, enabled,
	)
	if err != nil {
		return fmt.Errorf("failed to set module enabled: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetModuleSettings(guildID snowflake.ID, moduleID string) (string, error) {
	var settings string
	err := s.db.QueryRow(
		"SELECT settings FROM guild_module_settings WHERE guild_id = ? AND module_id = ?",
		int64(guildID), moduleID,
	).Scan(&settings)
	if err == sql.ErrNoRows {
		return "{}", nil
	}
	return settings, err
}

func (s *SQLiteStore) SetModuleSettings(guildID snowflake.ID, moduleID string, settings string) error {
	_, err := s.db.Exec(
		`INSERT INTO guild_module_settings (guild_id, module_id, settings) VALUES (?, ?, ?)
		 ON CONFLICT(guild_id, module_id) DO UPDATE SET settings = excluded.settings`,
		int64(guildID), moduleID, settings,
	)
	return err
}

func (s *SQLiteStore) CreateTicket(guildID snowflake.ID, number int, channelID, userID snowflake.ID, subject string) error {
	_, err := s.db.Exec(
		"INSERT INTO tickets (guild_id, number, channel_id, user_id, subject) VALUES (?, ?, ?, ?, ?)",
		int64(guildID), number, int64(channelID), int64(userID), subject,
	)
	return err
}

func (s *SQLiteStore) GetTicketByChannel(channelID snowflake.ID) (*Ticket, error) {
	var t Ticket
	var guildID, chID, userID int64
	var closedAt sql.NullTime
	var closedBy sql.NullInt64

	err := s.db.QueryRow(
		"SELECT guild_id, number, channel_id, user_id, subject, created_at, closed_at, closed_by FROM tickets WHERE channel_id = ?",
		int64(channelID),
	).Scan(&guildID, &t.Number, &chID, &userID, &t.Subject, &t.CreatedAt, &closedAt, &closedBy)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	t.GuildID = snowflake.ID(guildID)
	t.ChannelID = snowflake.ID(chID)
	t.UserID = snowflake.ID(userID)
	if closedAt.Valid {
		t.ClosedAt = &closedAt.Time
	}
	if closedBy.Valid {
		id := snowflake.ID(closedBy.Int64)
		t.ClosedBy = &id
	}
	return &t, nil
}

func (s *SQLiteStore) CloseTicket(channelID snowflake.ID, closedBy snowflake.ID) error {
	_, err := s.db.Exec(
		"UPDATE tickets SET closed_at = ?, closed_by = ? WHERE channel_id = ?",
		time.Now().UTC(), int64(closedBy), int64(channelID),
	)
	return err
}

func (s *SQLiteStore) DeleteTicket(channelID snowflake.ID) error {
	_, err := s.db.Exec("DELETE FROM tickets WHERE channel_id = ?", int64(channelID))
	return err
}

func (s *SQLiteStore) CreateRSSFeed(feed *RSSFeed) error {
	err := s.db.QueryRow(
		"INSERT INTO rss_feeds (guild_id, url, channel_id, title) VALUES (?, ?, ?, ?) RETURNING id",
		int64(feed.GuildID), feed.URL, int64(feed.ChannelID), feed.Title,
	).Scan(&feed.ID)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint") {
		return fmt.Errorf("%w: %s", ErrDuplicateFeed, feed.URL)
	}
	return err
}

func (s *SQLiteStore) DeleteRSSFeed(id int64, guildID snowflake.ID) error {
	// Delete seen items first (SQLite foreign key support varies)
	_, _ = s.db.Exec("DELETE FROM rss_seen_items WHERE feed_id = ?", id)
	_, err := s.db.Exec("DELETE FROM rss_feeds WHERE id = ? AND guild_id = ?", id, int64(guildID))
	return err
}

func (s *SQLiteStore) GetRSSFeeds(guildID snowflake.ID) ([]RSSFeed, error) {
	rows, err := s.db.Query(
		"SELECT id, guild_id, url, channel_id, title, added_at FROM rss_feeds WHERE guild_id = ? ORDER BY added_at",
		int64(guildID),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRSSFeeds(rows)
}

func (s *SQLiteStore) GetAllRSSFeeds() ([]RSSFeed, error) {
	rows, err := s.db.Query("SELECT id, guild_id, url, channel_id, title, added_at FROM rss_feeds ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRSSFeeds(rows)
}

func scanRSSFeeds(rows *sql.Rows) ([]RSSFeed, error) {
	var feeds []RSSFeed
	for rows.Next() {
		var f RSSFeed
		var gid, chid int64
		if err := rows.Scan(&f.ID, &gid, &f.URL, &chid, &f.Title, &f.AddedAt); err != nil {
			return nil, err
		}
		f.GuildID = snowflake.ID(gid)
		f.ChannelID = snowflake.ID(chid)
		feeds = append(feeds, f)
	}
	return feeds, rows.Err()
}

func (s *SQLiteStore) CountRSSFeeds(guildID snowflake.ID) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM rss_feeds WHERE guild_id = ?", int64(guildID)).Scan(&count)
	return count, err
}

func (s *SQLiteStore) IsItemSeen(feedID int64, itemHash string) (bool, error) {
	var exists int
	err := s.db.QueryRow(
		"SELECT 1 FROM rss_seen_items WHERE feed_id = ? AND item_hash = ?",
		feedID, itemHash,
	).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *SQLiteStore) MarkItemsSeen(feedID int64, itemHashes []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO rss_seen_items (feed_id, item_hash) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, h := range itemHashes {
		if _, err := stmt.Exec(feedID, h); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *SQLiteStore) PruneSeenItems(olderThan time.Time) error {
	_, err := s.db.Exec("DELETE FROM rss_seen_items WHERE seen_at < ?", olderThan)
	return err
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
