# Pedmin - Discord Bot

Pedmin (pepe + administrator) is a modular Discord bot built with Go 1.26.1 and disgo v0.19.2. It serves as a Probot replacement, featuring Components V2 UI, music playback via Lavalink, and a layered Feature Module architecture. Runs on Windows Docker Desktop.

## Tech Stack
- **Language**: Go 1.26.1
- **Discord Library**: disgo v0.19.2
- **Lavalink Client**: disgolink v3.1.0
- **Lavalink Server**: Lavalink 4 (Alpine)
- **Data Storage**: SQLite (`modernc.org/sqlite`, pure Go), behind `GuildStore` interface

## Commands
```bash
# Build
go build ./...

# Run tests
go test ./...

# Vet
go vet ./...

# Docker
docker compose up        # Start bot + Lavalink
docker compose up -d     # Detached mode
docker compose build     # Rebuild bot image
```

## Architecture: Layered Feature Module Pattern

Each feature is a self-contained module with internal layer separation (handler/service/view), all within the same Go package.

```
main.go                        # Entrypoint: DI wiring, graceful shutdown
config/config.go               # Env var loading
module/module.go               # Module interface definition
bot/
├── bot.go                     # Client init, module registry, lifecycle
├── commands.go                # Global command sync
├── router.go                  # Interaction → Module dispatch
├── ui.go                      # Shared UI helpers (errorMessage)
└── voice.go                   # VoiceState/VoiceServer → Lavalink relay
store/
├── store.go                   # GuildStore interface
└── sqlite_store.go            # SQLite implementation (WAL mode)
features/settings/
├── module.go                  # Info, Commands, empty stubs
├── handler.go                 # HandleCommand / HandleComponent
└── view.go                    # UI builders (mainPanel, modulePanel)
features/player/
├── module.go                  # Info, Commands, empty stubs
├── handler_command.go         # /player slash command
├── handler_component.go       # Button/select switch dispatch
├── handler_modal.go           # Add-to-queue modal
├── service.go                 # Playback logic (Discord API independent)
├── voice.go                   # VC connection helper
├── queue.go                   # Queue data structure
├── queue_manager.go           # Per-guild queue management
├── loop_mode.go               # LoopMode type + constants
├── lavalink.go                # Lavalink event listeners + node connection
├── view_player.go             # Player UI builder
├── view_queue.go              # Queue UI builder
└── view_helpers.go            # Progress bar, duration format, thumbnails
```

## Key Design Decisions

### 1 File = 1 Responsibility
Every `.go` file has a single, clear responsibility. No file mixes handler logic with UI building or service logic.

### Feature Module Pattern
Each feature (`features/player/`, `features/settings/`) is a self-contained Go package. Internal layers (handler → service → view) are separated by file, not by package. Same `package player` throughout — no circular import issues.

### Module Interface (`module.Module`)
All features implement: `Info()`, `Commands()`, `HandleCommand()`, `HandleComponent()`, `HandleModal()`, `SettingsPanel()`, `HandleSettingsComponent()`. Registered in `main.go` via `bot.Register()`.

### CustomID Convention
Component CustomIDs follow `{moduleID}:{action}:{extra}`. Router splits on the first colon to dispatch.

### Components V2
All UI uses `discord.NewMessageCreateV2()`. View files are pure functions: state in → components out.

### GuildStore Interface
`store.GuildStore` abstracts persistence. SQLite at `data/pedmin.db` with WAL mode. No JSON fallback.

## Documentation
- `docs/ARCHITECTURE.md` - System architecture, layers, data flow
- `docs/MODULE_GUIDE.md` - How to create new modules
- `docs/COMPONENTS_V2.md` - Components V2 reference for disgo
- `docs/LAVALINK.md` - Lavalink integration guide
- `docs/STORE.md` - Data persistence guide
