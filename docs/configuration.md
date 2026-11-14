# Mycelium — Configuration

Every variable and file the CLI and the server actually read, and what happens when one is
missing.

## The three roots

Mycelium follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir/latest):
one config file, one data tree, one state directory. Configured data is portable, machine-local
state is not, so they never share a directory.

| Root | Default | Holds | Overrides it |
|---|---|---|---|
| Config | `~/.config/mycelium/config.yml` | `machine`, `url`, `token`, `agents`, `rule_order` | `XDG_CONFIG_HOME` |
| Data | `~/.local/share/mycelium/` | `memory/`, `rules/`, `skills/`, `sessions/`, `events/`, `flows/`, `machines/`, `usage/`, `.sync-base.json`, `.conflicts/`, `.git` | `DATA_DIR`, `XDG_DATA_HOME` |
| State | `~/.local/state/mycelium/` | `daemon.log`, `.last-install`, `.sessions-state.json`, `.consolidate-cursor.json`, `.consolidate-run.json`, `.flow-trust.json`, `runs/` | `MYCELIUM_STATE_DIR`, `XDG_STATE_HOME` |

When `DATA_DIR` points somewhere explicit, its state sits beside it at `<DATA_DIR>/state` instead,
so one placed store stays one self-contained tree.

Migration: the first run of a new binary moves a pre-XDG install into place. The data directory
`$HOME/.mycelium` becomes `~/.local/share/mycelium`, the config file `$HOME/.mycelium.yml` becomes
`~/.config/mycelium/config.yml`, and the state files are lifted out of the data tree into
`~/.local/state/mycelium`. The migration is idempotent, never moves the data tree when `DATA_DIR`
is set, and is reported by `mycelium doctor`.

## Server environment

`internal/env` reads the whole configuration once, at startup, through `tronc/env`. Any
error it returns is a reason not to start: `mycelium serve` logs it and exits **1**. That is
deliberate — a half-configured server is worse than one that refuses to boot.

| Variable | Required | Default | What it does |
|---|---|---|---|
| `APP_ENV` | no | `development` | `development`, `staging`, or `production`. Anything unrecognized parses as `development` |
| `PORT` | no | `8420` | HTTP listen port. Must be 1–65535 |
| `LOG_LEVEL` | no | `info` | Level for the tronc structured logger |
| `DATA_DIR` | no | `~/.local/share/mycelium` | Root of the markdown tree and the server's JSON state |
| `CLIENT_DIR` | no | `./client` | Directory holding the built SPA. The image sets `/client` |
| `ENCRE_BIN` | no | `encre` on `PATH` | The [encre](https://github.com/FacileStudio/encre) binary that validates artifact bodies on `POST` and `PUT /api/artifacts`. The image downloads the static binary of release `ENCRE_VERSION` (Dockerfile), checks it against the release's `checksums.txt` and sets `/usr/local/bin/encre`. Missing, the server logs an error at boot and stores bodies unvalidated. The CLI reads the same variable for `mycelium artifact add` |
| `PASSWORD` | conditional | — | Shared password for `mycelium login --password` |
| `CORS_ALLOWED_ORIGINS` | no | — | Comma-separated origins. Empty denies every cross-origin browser caller |
| `OIDC_ISSUER` | no | — | Issuer URL of the identity provider. Setting it turns SSO on and makes the three below required |
| `OIDC_CLIENT_ID` | conditional | — | Required once `OIDC_ISSUER` is set |
| `OIDC_CLIENT_SECRET` | conditional | — | Required once `OIDC_ISSUER` is set |
| `OIDC_REDIRECT_URL` | conditional | — | Required once `OIDC_ISSUER` is set. Points at `/api/auth/oidc/callback` |
| `OIDC_SUCCESS_URL` | no | `<request base>/auth/callback` | Where the callback sends the browser, with the session in the URL fragment |
| `OIDC_MACHINE_AUDIENCE` | no | — | The `aud` a delegated suite JWT must carry, `mycelium` in production. Setting it lets a Registre token-exchange JWT reach the delegated allowlist (see [architecture](architecture.md#authentication)). Unset, only Mycelium's own tokens authenticate. Requires `OIDC_ISSUER` |
| `MCP_RESOURCE` | no | — | Canonical URI of the remote MCP endpoint, `https://mycelium.facile.studio/mcp` in production. Setting it serves `POST /mcp` and its RFC 9728 metadata for Registre access tokens bound to that URI (see [API](api.md#remote-mcp)). Unset, `/mcp` is off. Requires `OIDC_ISSUER` |
| `JOURNAL_URL` | no | — | Journal ingest URL. Log shipping needs both this and the token |
| `JOURNAL_TOKEN` | no | — | Per-app Journal key |

`APP_ENV`, `PORT`, `LOG_LEVEL`, `CORS_ALLOWED_ORIGINS`, `JOURNAL_URL` and `JOURNAL_TOKEN`
are `troncenv.Core` fields, so they carry the same names here as in every other Facile API.
Mycelium fills `Core` field by field rather than calling `troncenv.LoadCore`, because
`LoadCore` requires `DATABASE_URL` and Mycelium has no database.

`CORS_ALLOWED_ORIGINS` is read through `troncenv.CORSOrigins`, which also accepts the older
suite spellings — `ALLOWED_ORIGINS`, `DOMAINS`, `DOMAIN`, `CORS_ORIGINS`,
`TRUSTED_ORIGINS`, `CLIENT_ORIGIN` — first one set wins.

### The refusals

`mycelium serve` exits 1 rather than start when:

- `OIDC_ISSUER` is set but any of `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`,
  `OIDC_REDIRECT_URL` is missing.
- `OIDC_MACHINE_AUDIENCE` or `MCP_RESOURCE` is set but `OIDC_ISSUER` is not.
- `MCP_RESOURCE` is not an absolute URI.
- `APP_ENV=production` with neither `PASSWORD` nor `OIDC_ISSUER` — otherwise every request
  is served as admin.

Outside production, running with neither `PASSWORD` nor `OIDC_ISSUER` is allowed and every
request is served as admin. The server logs a warning at startup saying exactly that.

### Search needs nothing outside the binary

Memory search is BM25 over the markdown files in `DATA_DIR`, ranked per `### finding` block.
There is no index to build and no model to run. Since v0.53.0 neither the server nor the
CLI reads `OLLAMA_URL`, `EMBED_MODEL`, `VECTOR_STORE` or `QDRANT_URL`; a stale value in the
deploy panel is harmless and can be removed.

### Flags that override the environment

```sh
mycelium serve --port 9000 --data /srv/mycelium
```

`--port` only takes effect when explicitly passed; `--data` overrides `DATA_DIR`.

## CLI configuration

The CLI keeps its state in two places.

`~/.config/mycelium/config.yml` is written by `mycelium login` and read by every command that talks to a
server:

```yaml
machine: lucy
url: https://mycelium.facile.studio
token: <bearer token>
space: <space uuid>
usage_token: <subscription oauth token>
rule_order:
  - 00-core
  - 10-style
agents:
  - claude
  - codex
```

| Key | What it does |
|---|---|
| `machine` | This machine's name — the token name, the shard directory, the machine block |
| `url` | Mycelium server base URL |
| `token` | Bearer token for that server |
| `space` | Space to sync; empty means the common tree |
| `usage_token` | Last-resort store for the subscription usage token, when no OS keychain exists |
| `rule_order` | Rules emitted first, in this order, before the rest alphabetically |
| `agents` | Agents the daemon refreshes; empty means autodetect |
| `consolidate` | The daemon's consolidation stage; both keys optional, see below |

The file is written mode `0600`, because `usage_token` may hold a credential.

`~/.local/share/mycelium` (or `$DATA_DIR`) is the data tree itself: `memory/`, `rules/`, `skills/`,
`machines/`, `sessions/`, `usage/`. One local-only file sits beside it and is never synced:
`.sync-base.json`, the reconcile manifest. `.sessions-state.json`, the per-transcript byte
offsets, lives in `~/.local/state/mycelium/`.

`DATA_DIR` is read by the CLI too, not just the server. `mycelium rules edit` shells out to
`$EDITOR`.

### Consolidation

The daemon's consolidation stage reads recent episodes from `~/.local/share/mycelium/events/`, proposes
candidate findings heuristically, runs them through the storage gate, and writes the ones that pass
and are not near-duplicates into `memory/`. No model is involved. One key tunes it:

```yaml
consolidate:
  enabled: true
```

| Key | What it does |
|---|---|
| `enabled` | Turns the stage on. Defaults to `true`; an absent key keeps the default, an explicit `false` turns the stage off |

A `judge_model` key left in an old config file is ignored.

### The server session

`mycelium login` writes `url` and `token` into `~/.config/mycelium/config.yml` at mode `0600`. Two environment
variables override that file, with precedence **flag > environment > config file**:

| Variable | Overrides |
|---|---|
| `MYCELIUM_TOKEN` | The session token |
| `MYCELIUM_URL` (or `MYCELIUM_SERVER_URL`) | The server to talk to |

A CI job cannot run an interactive login and must not commit a config file, so the
environment is the only credential channel it has. Nothing ever writes these back to disk.

`mycelium logout` revokes the session on the server when it can reach it and clears `token`
from `~/.config/mycelium/config.yml`, leaving `url`, `machine`, `space`, `rule_order`, `agents` and
`usage_token` alone. Running it while logged out is not an error.

### The usage token

`mycelium usage --live` needs a subscription OAuth token from `claude setup-token`. Everything
else about usage tracking works without one. `ResolveToken` walks four sources in order of
decreasing safety and takes the first non-empty one:

| Order | Source |
|---|---|
| 1 | `CLAUDE_CODE_OAUTH_TOKEN` in the environment |
| 2 | `MYCELIUM_USAGE_TOKEN` in the environment |
| 3 | The OS keychain, service `mycelium-usage-token` — `security` on macOS, `secret-tool` on Linux |
| 4 | `usage_token` in `~/.config/mycelium/config.yml`, plaintext, last resort |

`mycelium usage login` writes to the keychain and only falls through to `~/.config/mycelium/config.yml` when the
machine has no keychain backend Mycelium can drive, saying so out loud. It refuses outright if
`~/.config/mycelium/config.yml` resolves inside the synced data directory.

The token is only ever used to read subscription limits from Anthropic's usage endpoint. It is
never sent to the Mycelium server, never written into the synced data directory, and never used
to make model requests. Claude Code's own `Claude Code-credentials` keychain entry is
deliberately never read — rotating that token would sign you out of your own CLI.

## Server state files

These live under `DATA_DIR` and are excluded from the file sync:

| File | Holds |
|---|---|
| `tokens.json` | Token hashes, name, scope, user email, timestamps. Mode `0600` |
| `.users.json` | OIDC users keyed by email, first one flagged admin |
| `.spaces.json` | Spaces and their membership roles |
| `.settings.json` | Antenne settings, managed through `PUT /api/settings` |
| `.pool-ledger.json` | Keys already emitted to the Antenne: session block IDs, plus `usage:`-prefixed alert keys |

## Antenne settings

Not environment variables — these are edited from the dashboard's Settings page, through the
admin-scoped `PUT /api/settings`, and stored in `.settings.json`.

| Field | What it does |
|---|---|
| `enabled` | Turns the emitter on. Requires `instance` and `secret` |
| `instance` | Antenne instance URL. Must parse as `http` or `https` with a host |
| `secret` | Pool secret |
| `user_email` | Default `user_email` on emitted events |
| `machine_emails` | Per-machine override of `user_email` |
| `emit_since` | RFC3339 watermark. Defaults to now on first enable, so no backfill |
| `usage_alerts` | Publishes `usage_alert.created` when a subscription window crosses the threshold. Defaults `false` — opt-in, like the emitter itself |
| `usage_threshold` | Percent at which a window alerts. Defaults **80** |

An absent or `0` `usage_threshold` resolves to 80, never to "alert at 0%", and the value is
clamped to 1–100 with nonsense ignored rather than rejected. Every call site reads it through one
accessor, so no path can pick up the raw zero value.

`machine_emails` does more than label usage alerts: their dedupe key is keyed on the resolved
email, because a subscription limit is per Anthropic account. Two machines mapped to the same
email therefore produce **one** alert per window instance, not two.

`usage_alerts` publishes an event and nothing else — no email, no push, no webhook. Antenne is the
alert aggregator and owns delivery.

Back to the [documentation index](README.md).
