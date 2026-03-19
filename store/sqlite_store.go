package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
