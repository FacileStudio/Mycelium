# Changelog

All notable changes to this project are documented here. The format is
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project
follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html). While on
`0.x`, a breaking change bumps the minor.

Every entry below was reconstructed from git history on 2026-08-24, so it
records what shipped rather than what was written down at the time.

## [0.47.3] - 2026-11-12

### Fixed

- MCP: `mycelium_publish_artifact` records the client acting for the person as the artifact's machine. A delegated token from Bulle's token exchange (`act.sub` = `bulle`) records `bulle`, so Bulle's artifacts badge as Bulle instead of `mcp`. A direct token, or a delegated one whose client id is an opaque generated id, still records `mcp`. Replacing an artifact by title keeps the machine it was first recorded from.

## [0.47.2] - 2026-11-11

### Fixed

- Brand assets come from the muse v0.19.0 generator: favicon.svg, favicon.ico, the touch and manifest icons, og-image.png and site.webmanifest now match the rest of the suite. `src/app.html` carries the generated head block, so the title, description, theme colours and every Open Graph and Twitter tag appear once, with the catalog's wording. The landing page title and description follow the catalog in every locale. `scripts/icons.ts` is gone: it regenerated the old set.

## [0.47.1] - 2026-11-11

### Fixed

- The English-only rule now covers the common tree alone, the owner's wiki whose search model reads English. A page, finding or delete reason in a space takes any language: `mycelium_write_memory_page`, `mycelium_add_finding` and `mycelium_delete_memory_page` used to refuse a French member writing their own space's memory. The CLI follows the same split: `mycelium memory add`, the sync warning and `mycelium doctor`'s `wiki language` check apply only on a machine on the common tree, not one bound to a space. The REST sync path never checked the language and still does not.
- MCP: `mycelium_list_artifacts` and `mycelium_read_artifact` answer the `machine` that recorded the artifact when it has one, so Bulle's Artifacts page shows its machine badge again. Additive: no other field changed.

## [0.47.0] - 2026-11-11

### Security

- build(deps): `github.com/go-jose/go-jose/v4` v4.1.4 to v4.1.6, `vite` 7.3.5 to 7.3.8, and client overrides for the transitive `devalue` (5.9.4), `nanoid` (3.3.20), `postcss` (8.5.32) and `source-map-js` (1.2.2), which a Sonde scan reported as critical or high. `douane scan` reports only low findings after.

### Added

- MCP: `mycelium_write_memory_page` creates or replaces a whole wiki page with the wiki's rules (frontmatter with `title:` and a descriptive `type:`, `updated:` and `created:` stamped, index line, `create-page` or `lint` log line, English only, 256 KiB cap). `mycelium_delete_memory_page` deletes a page only after the person confirms, through MRTR `input_required` or a pending action and the new `/pending` confirm page (`GET /api/mcp/pending`, `POST /api/mcp/pending/{id}/confirm|dismiss`, browser sessions only). `mycelium_list_memory_pages` lists a tree's pages, and `mycelium_create_space` creates a space the caller owns, reusing one of the same name. Both page tools refuse `index.md`, `log.md` and standards with a tool error. Nothing forgets or ratifies.
- MCP: `mycelium_update_artifact` rewrites an artifact by exact id (optional rename, encre validation and the 1 MiB cap as on `PUT /api/artifacts/{id}`), and `mycelium_delete_artifact` deletes one only after the person confirms, like a page delete.
- A Registre erasure also drops the person's pending MCP actions.

### Changed

- encre [v0.9.4](https://github.com/FacileStudio/encre/releases/tag/v0.9.4) validates and renders artifacts (`@facile/encre`, `ENCRE_VERSION`, `github.com/FacileStudio/encre/go`). An oversized input now comes back as an invalid verdict, which the adapter already reads as a refusal; Mycelium's own 1 MiB cap still answers first. Pasted tables are normalized in the editor and a failed wasm load is no longer cached.
- The web app keeps the page you opened through sign-in: a signed-out visit to `/pending?action=...` comes back there.

## [0.46.0] - 2026-11-11

### Added

- Remote MCP server at `/mcp` (set `MCP_RESOURCE`), on `tronc/mcp` v0.22.1 and `porte/mcpauth` v0.2.1, for Registre access tokens bound to it with `mycelium.read` or `mycelium.write`. Seven tools: `mycelium_read_memory_index`, `mycelium_search_memory`, `mycelium_read_memory_page`, `mycelium_list_artifacts`, `mycelium_read_artifact`, `mycelium_add_finding` and `mycelium_publish_artifact` (create, or replace by title, with the 1 MiB cap and encre's refusal placed at a line and column). A caller reaches exactly the trees the same person reaches in the web app, and the first call creates a non-admin account. Nothing deletes, forgets or ratifies; flows stay on the local stdio `mycelium mcp`.

### Changed

- build(deps): `tronc` v0.17.1 to v0.22.1, `porte` v0.12.1 to v0.18.4, `porte/directory` v0.3.0 to v0.4.0, `modelcontextprotocol/go-sdk` v1.7.0 to v1.8.0.

## [0.45.4] - 2026-11-11

### Security

- build: Go 1.26.11 (`go` directive, `golang:1.26.11-alpine` builder, mise pin and lock). It fixes GO-2026-6848 to GO-2026-6859 in the standard library and toolchain, which douane reported as unscored gaps and which failed the gate.
- build(deps): `golang.org/x/net` v0.56.0 to v0.62.0 (GO-2026-6612, critical; GO-2026-6860, GO-2026-6866), which requires `golang.org/x/crypto` v0.58.0 (from v0.54.0, clearing GO-2026-6303); `x/sys` and `x/term` follow.
- build(client): `dompurify` 3.4.14 to 3.4.17 (CVE-2026-116804, critical).
- `douane.yml` permits GO-2026-5932 on `golang.org/x/crypto` until 2027-02-09: its fix exists only as an unreleased x/crypto commit, and Mycelium links `ssh`, `argon2`, `blake2b`, `hkdf`, `sha3` and `cast5`, never `openpgp`.

### Changed

- refactor: `mycelium stats` capitalises its column header without the deprecated `strings.Title` (SA1019), and the unused `pageRankOf` test helper is gone (U1000). Both failed `filet gate`.

## [0.45.3] - 2026-11-11

### Changed

- encre [v0.9.2](https://github.com/FacileStudio/encre/releases/tag/v0.9.2) validates and renders artifacts. Its new `EncreDocumentError` only concerns serializing, which the viewer never does.

## [0.45.2] - 2026-11-11

### Changed

- encre [v0.9.1](https://github.com/FacileStudio/encre/releases/tag/v0.9.1) validates and renders artifacts.

### Fixed

- The viewer's notice for an artifact encre refuses (one stored through the unvalidated sync path) now places the refusal as the API does: `line L, column C: <message>`, from encre's `validateMarkdown` on the stored body, for a document encre cannot read as for one that reads but does not validate. Line numbers count from the first line under the frontmatter, as on the server: the viewer now drops the blank lines under the frontmatter too.

## [0.45.1] - 2026-11-11

### Changed

- encre [v0.9.0](https://github.com/FacileStudio/encre/releases/tag/v0.9.0) validates and renders artifacts. Every markdown path of encre is now linear on hostile input, and encre refuses markdown that nests more than 64 levels or opens more than 8 `callout`, `columns` or `column` blocks on one path, with the line and column of the fence that crosses the limit: `400 invalid document: line L, column C: custom blocks nest more than 8 levels deep` on the API, the same text on the CLI and `publish_artifact`. The viewer shows such a document in degraded mode (the notice with encre's reason, then the source). A validation message no longer starts with the child's path (`content[0]: `); the `path` field says where the problem is. The 1 MiB cap stays.

## [0.45.0] - 2026-11-10

### Added

- Artifacts are validated and rendered by [encre](https://github.com/FacileStudio/encre) v0.8.0, whose markdown now carries diagrams (` ```mermaid `, ` ```d2 `, ` ```dot `), `columns` (and the `compare`, `pros-cons` and `vs` aliases), `[badge:tone text]`, `[status:state]`, `==highlight==`, per-line tones for `diff` code blocks and `stat.sub`. The web viewer draws `mermaid` diagrams with the app's mermaid, loaded only when a document holds one, in the app's theme; `d2` and `dot` show their source. Code blocks are highlighted with highlight.js through DOMPurify, and the viewer's copy follows the app's locale.
- A document encre refuses is reported with its place: `invalid document: line L, column C: <reason>` on the API, the CLI and the `publish_artifact` tool.
- Every write path (`POST` and `PUT /api/artifacts`, `mycelium artifact add`, `publish_artifact`) refuses a body over 1 MiB before validation runs: `413` on the API, a refusal on the CLI and the tool.

### Changed

- The server and the CLI call encre through its Go module (`github.com/FacileStudio/encre/go`) instead of their own subprocess code. `ENCRE_BIN` and the PATH lookup are unchanged.
- The pre-encre markdown renderer (`MarkdownMuse`, its mermaid block and modal), which nothing rendered since 0.44.0, is removed. Legacy `chart:bar`, `chart:metrics` and `> [!NOTE]`-style alerts are still mapped when the page renders.

## [0.44.0] - 2026-11-10

### Added

- Artifacts are validated by [encre](https://github.com/FacileStudio/encre) v0.7.0 on every write. `POST` and `PUT /api/artifacts`, `mycelium artifact add` and the `publish_artifact` MCP tool run `encre validate --from markdown` on the body under the frontmatter and refuse a document encre refuses with its reason (`400 invalid document: ...` on the API). The image downloads encre's static `encre` binary for the release `ENCRE_VERSION` in the Dockerfile and checks it against the release's `checksums.txt`; `ENCRE_BIN` points the server or the CLI at another binary. A machine without encre records artifacts with a warning, and the server logs an error at boot.
- The web viewer renders artifacts with encre's `Document` component. Mycelium's older fences are mapped when the page renders, never on disk: `chart:bar` (and a non-JSON `chart`) become an encre `chart`, each `chart:metrics` card a `stat`, and `> [!NOTE]`-style alerts a `callout`. `mermaid`, `diff`, `compare`, badges and `==highlight==` have no encre block yet and show as plain code or text. An artifact encre cannot validate still opens, with a notice and its source as the fallback.

### Fixed

- `mycelium artifact add` of a file that carried its own frontmatter kept that frontmatter's closing `---` line at the top of the stored body.

## [0.43.0] - 2026-11-10

### Added

- Mycelium is licensed under AGPL-3.0-only (`LICENSE`). A commercial licence is available from Facile Studio.

### Fixed

- Registre's `user.deleted` notice now erases an account that was provisioned by a delegated suite call (Bulle) or by a Registre space snapshot. Those paths stored no subject, so the notice found no account and erased nothing while the receiver recorded it as handled. The subject is now stored on every path that learns it, and an account that only carries the subject's picture URL (`<issuer>/avatars/<subject>`) is matched through it. The erasure removes the `.users.json` record, the person's tokens and sessions, their memberships, and every space, artifacts included, that nobody else uses and Registre does not own.
- `scripts/erasure-check.sh <subject> <email>` counts, read-only, what every suite database on this host still holds for an erased person, with no row contents.

## [0.42.0] - 2026-11-10

### Removed

- HTML artifacts. Artifacts are markdown only: `mycelium artifact add` of an `.html` file or an HTML document, the `publish_artifact` MCP tool and `POST`/`PUT /api/artifacts` refuse HTML with "artifacts are markdown only". `.html` files already in a tree stay on disk untouched but are no longer listed, found, served or rendered. The `format` field is gone from the artifact API.

### Changed

- The web viewer renders raw HTML inside markdown as literal text, never as markup, and shows a `chart` fence holding JSON (encre's chart block) as plain code instead of an empty bar chart. Unknown fences stay plain code.
- An artifact body is stored byte for byte under the frontmatter, so fenced custom blocks survive untouched.

## [0.41.0] - 2026-11-10

### Added

- `GET /api/artifacts?all=true` lists every artifact the caller can read, newest first: the common tree when they may use it, then each space they belong to. Each entry carries `space_id` and `space_name`.
- `GET /api/spaces` answers `common_artifacts`: whether the caller's personal artifacts live in the common tree.

### Changed

- A delegated suite token of an admin (Bulle acting for that admin) reaches the common tree on the artifact routes only, so Bulle lists the same artifacts as `mycelium artifact list` and what Bulle saves syncs to the admin's machines. Rules, memory, skills and sync stay closed to every delegated token, and a non-admin still never reaches the common tree.

## [0.40.4] - 2026-11-07

### Changed

- build(client): bump muse to v0.16.1 (from v0.14.4). The secondary buttons use the `secondary` variant in place of the deprecated `outline`.
- `GET /version` and `mycelium --version` in the container come from the `VERSION` file, which the release commit bumps. The `APP_VERSION` build arg and the Dokploy env no longer take part, so a stale env value cannot make the server report an old release.

## [0.40.3] - 2026-11-07

### Changed

- build(client): bump muse to v0.14.4 (from v0.14.2). Switching the language in the user menu re-renders the whole UI at once, with no reload. A single 100% DonutChart slice draws as a full ring.

## [0.40.2] — 2026-11-07

### Fixed

- Docker: the image version comes from the tag on the checked-out commit first, then the `APP_VERSION` build arg, then `git describe`. A stale `APP_VERSION` in the Dokploy env no longer makes `mycelium --version` report an old release.

### Tests

- Sync never deletes an artifact because of an old `expires` line, and mirrors a deletion made on the server.

## [0.40.1] — 2026-11-07

### Changed

- Web: the API tokens create drawer no longer shows an expiry field and the list no longer shows an Expires column, since Mycelium tokens do not expire.
- Libraries: `@facile/muse` v0.14.1 to v0.14.2 (the user menu stays inside the viewport).

## [0.40.0] — 2026-11-07

### Changed

- Web: the sidebar card and the phone avatar open the shared muse user menu (Facile account, settings, language, theme, help, sign out).
- Web: settings use the muse settings layout: Profile, Preferences, Connect, Bus, API tokens, Advanced. The API tokens list, create drawer and revoke dialog are the muse tokens section, and the same API calls. Appearance is now Preferences, and `/settings/api` and `/settings/appearance` redirect.
- Web: the theme choice applies before the first paint.
- Web: removed the hand-rolled card, badge, table, callout and toggle borders; surfaces use fills.
- Libraries: `@facile/muse` v0.13.0 to v0.14.1.

### Removed

- Web: the sign-out button inside settings (it lives in the user menu) and the local profile card.

 — 2026-11-07

### Changed

- Libraries: porte v0.12.0 to v0.12.1 and `@facile/muse` v0.12.1 to v0.13.0 (borderless
  controls with the light `fc-field` fill, `Select` as a custom combobox, badges that never
  wrap). The server keeps JSON storage and uses no `porte/pg`, so the porte fix does not touch
  it. No app code change: the page-create and space-member selects keep their labels.

## [0.39.0] — 2026-11-07

### Changed

- **Sign-in goes through Registre only.** With `OIDC_ISSUER` set, `POST /api/auth/login` is not
  mounted and the dashboard offers only "Sign in with Facile" (muse `AuthScreen`). `SSO_ONLY` is
  gone from the code, `.env.example` and docs: delete it from the deployment env after the deploy.
  `GET /api/auth/config` still answers `sso_only`, true whenever an issuer is set, so old CLIs
  keep choosing the browser or device flow. API tokens, machine tokens and sessions already issued
  keep working. Clients that signed in with the shared password must switch to a token or to
  `mycelium login`.
- **OIDC accounts link by subject, then by verified email.** `.users.json` rows gain `issuer`,
  `subject` and `locale`. A row that holds the subject signs in as itself. Otherwise the identity
  links by email only when the id_token says `email_verified`, and never onto a row that already
  holds another subject of the same issuer (a recycled address cannot take over a live account).
  An unverified email or a second subject ends on `/login?error=sso`. Existing rows link on their
  next sign-in. This is porte v0.10.0's rule, ported to Mycelium's own user store.
- porte v0.7.0 to v0.12.0, `porte/directory` v0.1.0 to v0.3.0, tronc v0.16.0 to v0.17.1,
  `@facile/muse` v0.11.0 to v0.12.1. No porte table is involved: Mycelium keeps its state in
  files, so there is no migration.
- tronc security headers are on. The CSP is report-only (framing is blocked and enforced); the
  inline script hashes of the built client and the iconify and Journal origins are allowed.
  The artifact preview iframe is a `srcdoc` and keeps working.
- Per-IP limits (password login, device grants) read the address the chassis resolved from
  `TRUSTED_PROXIES` (`private,cloudflare` behind Traefik) instead of the raw `X-Forwarded-For`,
  which any client could spoof. The configured `TRUSTED_PROXIES` now reaches the router.
- `GET /api/auth/me` returns the person's `locale`. The login page forwards `ui_locales` to the
  issuer.
- Every dashboard string is available in fr, en, es and de (es and de need native review), with
  dates, numbers and durations through `Intl` in the session locale. The landing page uses muse
  `LandingHeader` and `LandingFooter`.

### Added

- GDPR deletion. Mycelium polls Registre's deletion list and accepts its signed push on
  `POST /api/registre/deletions` (outside the bearer middleware). A deleted person loses their
  account, tokens and memberships; spaces only they used and Registre does not own are removed;
  shared spaces keep their data and an owner is promoted when none is left. Replays are no-ops.
  The cursor lives in `.deletions.json` in the data directory.
- `GET /version`, from the `APP_VERSION` build arg (falls back to `git describe`).
- CI builds the server and the client and checks the message catalogs.

## [0.38.1] — 2026-11-07

### Security

- The dashboard sanitises rendered artifact markdown with DOMPurify and renders mermaid
  diagrams in strict mode. Before, raw HTML in an artifact, such as an image with an
  `onerror` handler written by an agent, ran in the page.

## [0.38.0] — 2026-11-06

### Added

- `POST /api/artifacts` creates a markdown artifact and `PUT /api/artifacts/{id}` rewrites
  one in place, keeping its id and created time. Both honour `space_id`, accept markdown
  only, and cap content at 1 MB and the title at 200 characters.
- Artifacts record an `updated` time, shown in the CLI listing and the dashboard.
- Delegated suite tokens (Bulle) can list, read, create, update and delete artifacts in
  their own spaces, never in the common tree.

### Changed

- **Artifacts are permanent.** Nothing expires or sweeps them. Old files that carry
  `expires:` load unchanged and are never deleted. The API shapes lose `expires` and
  `expired`, and `publish_artifact` loses its `keep` input.
- A lookup by exact artifact id wins over a fuzzy slug match.

### Removed

- `mycelium artifact sweep`, and the sweep `artifact list` ran first. `--expires` on
  `artifact add` is still accepted, hidden, and does nothing, so old scripts keep working.

## [0.37.0] — 2026-11-06

### Added

- The landing page shows Facile's plans under "Tarifs", from Ardoise's public catalog. The
  section is absent while the catalog is empty or unreachable.

### Changed

- `@facile/muse` v0.7.4 to v0.11.0.

## [0.36.0] — 2026-11-03

### Added

- **Registre owns shared spaces.** The server imports every space with more
  than one member into Registre once, retrying with backoff, then applies
  Registre's snapshot at boot, every five minutes, and after a nudge on
  `POST /api/registre/sync`. Members are resolved by email; unknown people are
  provisioned as non-admin users. Renaming or deleting a linked space, and any
  membership change, now answers 409 with the Registre link, and the dashboard
  shows a "Manage in Registre" button.
- **Delegated suite tokens.** Setting `OIDC_MACHINE_AUDIENCE` turns on a second
  bearer path: a Registre token-exchange JWT, verified offline and mapped to the
  user by email. A delegated identity reaches only an explicit route allowlist,
  needs `mycelium.read` for GET and `mycelium.write` otherwise, and is refused
  the common tree. Unset, auth is unchanged.
- **Profile pictures for every person.** The server stores the OIDC `picture`
  claim at sign-in and exposes it as `avatar_url` on `/auth/me`, `/users` and
  space members. The sidebar, mobile nav, profile card and members table render
  it, falling back to the initial.

### Changed

- **Config, data and state move to XDG directories.** `~/.mycelium.yml` becomes
  `~/.config/mycelium/config.yml`, `~/.mycelium` becomes
  `~/.local/share/mycelium`, and logs, stamps and offsets move to
  `~/.local/state/mycelium`. The first path lookup migrates an existing install
  in place, idempotently, and stands down when `DATA_DIR` names an explicit
  store. The legacy config is still read when the new one is absent, so an
  unmigrated machine does not log out.
- **Spaces use the share-circle icon** across the dashboard.
- **The client markdown renderer and command wiring are split up.**
  `MarkdownMuse` is broken into smaller components and helpers, and formatting
  helpers are consolidated in `format.ts`. The client bumps muse to v0.7.3.

### Removed

- **`publish_report` MCP tool alias removed.** The compatibility alias left
  behind when `publish_artifact` became the primary tool name is no longer
  registered. Clients must call `publish_artifact`. Tool descriptions no longer
  refer to artifacts as "reports".

### Fixed

- **Every CLI command is registered again.** Builds from source shipped 13 of
  26 commands: `doctor`, `memory`, `status`, `init`, `flow`, `rules`,
  `sessions`, `daemon`, `mcp`, `claim`, `logout`, `diff`, `recap` and `usage`
  answered "unknown command". All `newXxxCmd` constructors are now called in one
  place.
- **The fire-and-forget artifact sync in the MCP server is explicit.** A failed
  sync after recording an artifact no longer looks like an ignored error.
- **`antenne-client` points at a revision that still exists**, so a clean
  module download no longer fails.
- **Tests no longer touch the real config.** Tests that moved `HOME` left
  `XDG_CONFIG_HOME` alone, so on any machine that sets it, one test wrote its
  fixture into the real config and another read it. Tests now isolate all XDG
  bases.

## [0.35.0] — 2026-09-10

### Fixed

- **A deploy no longer kills an in-flight `mycelium login`.** Pending device
  authorizations and one-time login codes were held only in process memory, so
  the twenty-second redeploy any push to main triggers invalidated a code the
  user was mid-approve or mid-paste. Both stores now persist to
  `devices.json` and `login-codes.json` under the data dir and reload at boot.
  Device codes are stored keyed by hash (matching `tokens.json`), the minted
  token is never written to disk, and a consumed code is removed at consume-time
  so a restart can't replay it. Writes are atomic (temp file + rename), so a
  torn file on crash cannot silently empty the store.

### Changed

- **Server state files are fenced from sync.** `devices.json` and
  `login-codes.json` join `tokens.json` in `syncSkip`, so a machine can't fetch
  them through the file sync that moves wiki content.

### Removed

- **The nacelle transcript collector.** `collectNacelle` tailed
  `~/.nacelle/sessions/` — a privacy-narrowed transcript — while Nacelle usage
  already arrives through the canonical `events/` feed. The redundant,
  transcript-reading collector is gone; no nacelle usage is lost.

## [0.34.0] — 2026-09-10

### Added

- **Every recorded flow is now its own MCP tool.** `~/.mycelium/flows/*.yml` each
  generate a dedicated `run_flow_<name>` tool carrying the flow's description, so a
  host calling Mycelium's MCP server gets a named entry per procedure instead of one
  generic `run_flow`. The generic `run_flow` and `list_flows` remain as fallbacks.

- **`mycelium spaces` matches nuage.** The CLI now renders the suite standard — a
  leading `*` selection marker, an explicit `common` row, 8-character short ids — and
  resolves a space by id prefix as well as full id and case-insensitive name, with an
  explicit `--json` output for script use.

### Changed

- **The compose file is swarm-safe and comment-light.** `docker-compose.yml` was
  reformatted for zero-downtime swarm deploys and sheds its block explanatory comments
  (the env-driven bus and search provenance is now documented in `docs/deployment.md`
  and `AGENTS.md` instead).
- **The API reference is generated from the scalar OpenAPI spec.** `internal/documentation`
  now derives the endpoint reference from OpenAPI rather than a hand-maintained list,
  so the docs cannot drift from the routes.
- **`filet.yml` harmonised to the strict suite standard**, and the release config drops
  the deprecated Homebrew cask `url.verified` key.

### Fixed

- **The container healthcheck probes the binary, not `nc`.** The base image carries no
  `nc`, so the compose healthcheck could never pass; it now runs `/mycelium healthcheck`
  at a tighter 10s interval.

## [0.33.0] — 2026-08-30

### Added

- **`mycelium login --no-listener`, for the terminal whose browser is on another
  machine.** The server has shown a paste code since 0.32.3 when a CLI login
  arrives with no loopback port, but the CLI always opened a listener and always
  put `port=` in the URL it printed, so there was no way to reach that page. A
  URL copied to a laptop sent the code to the laptop's own loopback, and the
  terminal that started the login waited out its timeout for a callback that was
  never addressed to it. With the flag the login URL carries no port, the server
  renders the code on the hand-off page, and the command reads it from stdin.

### Changed

- **The SSO loopback listener is now `porte/loopback`.** `cmd/login_sso.go`
  hand-rolled its own, which is the architecture gap the suite CLI standard
  already records. The shared listener keeps the ephemeral port, the exact
  comparison of the echoed nonce, the two second shutdown grace and the three
  minute timeout, and it renders the suite's hand-off page on `127.0.0.1` rather
  than one line of `text/plain` in the browser's default serif. Refusals render
  that same page instead of `http.Error`.
- A callback carrying the wrong nonce no longer ends the login. It is refused,
  the page says so, and the listener keeps waiting. Ending it let any page the
  user has open close a login it did not start, which is the same class of
  problem as a browser's unprompted `/favicon.ico` request.
- `porte/loopback` is standard library only, so the module adds nothing to the
  binary beyond itself and the page it renders.

## [0.32.3] — 2026-08-29

### Fixed

- **A failed `mycelium login` no longer dead-ends the browser on a JSON error object.** The OIDC
  callback read the `flow=cli` marker as "the caller is a program" and answered a refusal with
  the API envelope, so a browser sitting on a top-level navigation rendered
  `{"error":{"code":"unauthenticated","message":"unauthorized"}}` as a web page. The callback is
  a browser navigation whichever flow started it, so every failure now redirects to
  `/login?error=sso`, which is the rule the browser half was already fixed to on 2026-08-26.
  JSON stays on `POST /api/auth/oidc/exchange`, where the caller really is a program, and the
  `cli` flag `oidcIdentity` took for this is gone rather than left dead.

- **The approve button on `/authorize` is never missing without a reason.** It only existed once
  a lookup had succeeded, so a code that could not be found left the entry form on screen with
  no button and no explanation, which reads exactly as "the approve button is disabled". The
  page now runs on an explicit state and says which of the four things happened: nothing looked
  up yet, looking, refused and why, or found and here is the machine. A refusal separates an
  account that may not approve machines (403) from a code this server is not holding, and an
  expired dashboard session (401) returns to the login page carrying the code.

- **The device code field shows what it will send.** Entry is canonicalised on every keystroke
  and paste, to upper case, alphabet characters only, eight at most, with a dash after the
  fourth. Characters outside the alphabet are named rather than silently swallowed. An unfinished code
  is answered in a sentence instead of a greyed-out button. The alphabet and the shape live in
  `apps/client/src/lib/deviceCode.ts` alone, checked against the server's by
  `deviceCode.check.ts` in the repository gate.

- **The artifact cards no longer push the page sideways on a phone, and no longer carry a
  border.** Each card was a hand-rolled `<a>` with `border border-fc-border`, which CHARTE §5
  forbids on a container surface: a card separates itself from the page with its fill. Nothing
  in the card set `min-w-0`, so a 70-character artifact id in a grid item that defaults to
  `min-width: auto` widened its own track and the whole grid with it, and the `truncate` meant to
  prevent exactly that could never take effect. The card is now the `EntityCard` the other six
  list pages already use, with the badges in its `trailing` snippet and the machine and date in a
  new optional `footer`.

- **Two states on the artifact page rendered nothing at all.** `hover:bg-fc-surface-hover`,
  `hover:shadow-fc-sm` and `bg-fc-primary text-fc-primary-fg` name tokens that do not exist, so
  Tailwind emitted no CSS for them: cards had no hover step and the Preview/Source toggle had no
  visible selection. The selected tab is now the inverted `bg-fc-accent text-fc-accent-fg` the
  charte specifies for an active state, and the card hover comes from `Card` itself. Twenty-three
  more of these dead utilities remain in `MarkdownMuse.svelte`, `MermaidBlock.svelte` and
  `MermaidModal.svelte`; replacing `fc-primary` there is a palette decision, not a rename.

- **An artifact's own page fits a phone now.** The metadata line under the title was a `flex` row
  with no wrap holding a full artifact path, so it ran off the side of a phone. It wraps now,
  the path breaks rather than overflowing, and the machine name truncates.

### Added

- **`flow=cli` with no loopback port ends on a page showing a code to paste.** It used to be
  refused twice, at `/api/auth/oidc` and again at the callback, which left `mycelium login` from
  a machine whose browser lives elsewhere with no route through at all. The code is stored
  exactly as the loopback case stores it, same TTL and same single use, and the page is the
  markup every porte-based app in the suite renders at the end of a CLI login.

## [0.32.2] — 2026-08-29

### Fixed

- **`go test ./...` published a real artifact to the production server.** The MCP publish test
  pointed `MYCELIUM_URL` at `mycelium.facile.studio`, and `publish_artifact` syncs what it
  records. `DATA_DIR` was a temp dir but the credential was not: `LoadMyceliumConfig` reads
  `~/.mycelium.yml`, so the test authenticated as the developer's own machine and pushed a
  document titled "Inline Artifact" into the real store, which then synced everywhere. Three of
  them landed on 2026-08-29 before anyone noticed. The test now names a port nothing listens on,
  which exercises the same code and reaches nothing.

- **`artifact open` no longer puts prose on stdout.** The "no browser on this machine" note moved
  to stderr, so `mycelium artifact open <id> | pbcopy` copies the link and nothing else, as
  CLI-STANDARD §7.3 requires. A terminal still shows both lines.

- **`artifact add` syncs before it hands you the link.** It printed the server URL and opened a
  browser at it while the file was still on its way up, so a fast machine could race the upload
  and land on a 404. With no server configured it now prints the local path, which is the only
  thing that resolves there.

- **`cmd /c start <url>` reads a lone quoted argument as the window title.** Windows is not a
  release target today, but the opener now passes the empty title `start "" <url>` the way the
  documented invocation does.

## [0.32.1] — 2026-08-29

### Fixed

- **A machine with no browser gets the link, not an error.** `mycelium artifact open` shelled out
  to `xdg-open` whenever `DISPLAY` was set, so a container that inherits a display from its host
  without inheriting an opener died with a raw `exec: "xdg-open": executable file not found`, a
  usage dump, and a non-zero exit. The artifact was recorded and the link was never printed. Both
  `open` and `add` now print the link first and treat the browser as the optional part: a machine
  that has none says so on one dim line and exits `0`.

- **"Can this machine show a page" is one predicate, not two.** `internal/browser` answers it for
  the artifact commands and for `login`, which each carried their own version and each got a
  different half right. It now requires a display *and* an opener on `PATH`, so an SSH session
  into a Mac and a container without `xdg-open` both fall back the way a headless Linux box
  already did: `login` goes to the device flow rather than waiting three minutes for a browser
  that was never going to open.

## [0.32.0] — 2026-08-29

### Added

- **`mycelium artifact add` reads from stdin, and `open` goes to the hosted URL.** A page
  generated in a pipeline never has to touch a temp file, and an artifact recorded against a
  configured server opens at `https://<server>/artifacts/<id>` rather than at a local path
  nobody else can reach. `publish_artifact` returns the rendered content inline, so the model
  that recorded it can read back what it filed.

- **Journal logging in the backend and the dashboard.** The server ships its logs to Journal
  when `JOURNAL_URL` and `JOURNAL_TOKEN` are set, and the web UI connects through
  `@facile/journal` with credentials fetched from `GET /api/auth/config`, so a browser session
  never carries the service token.

- **An admin can promote another user to admin.** Bootstrapping a second administrator
  previously meant editing the database by hand.

- **Zoom and fullscreen for Mermaid diagrams in the artifact viewer.** A flowchart wider than
  the column was unreadable and had no way out of it. Fullscreen opens over a blurred backdrop
  and fits the diagram to the viewport on entry.

### Changed

- **A space is always explicit.** The unscoped common tree is gone from the client, and every
  request carries the space it belongs to. The switcher no longer offers a bucket whose
  contents depended on who was asking.

- **Artifacts wear `solar:file-smile-linear`** in the sidebar and on the empty state, instead
  of the generic card glyph shared with every other list.

- **muse pinned at v0.7.0, with `WordReveal` vendored into the client.** muse dropped the
  component as its only GSAP-plugin consumer and names Mycelium as the sole user, so it moves
  here with `gsap` as a direct dependency.

### Fixed

- **Space access is restricted to members, and the admin bypass is removed.** An admin could
  read every space on the instance regardless of membership. Membership is now the only thing
  the server checks.

- **The markdown palette maps onto muse tokens.** muse v0.7.0 sets `--color-*: initial`, which
  unsets Tailwind's stock palette, so 57 utilities in `MarkdownMuse.svelte` emitted no CSS at
  all: callouts, diff badges, status pills and highlights rendered as unstyled text. Nothing
  errored and the build stayed green, which is how it nearly shipped.

- **An admin can select their personal space again.** The switcher suppressed the row by
  passing `undefined`, which Svelte resolves to the prop's default, so the row rendered for a
  non-admin and `pickSpace` then refused the click. `null` is the suppression signal.

- **`publish_artifact` syncs the artifact to the server.** It wrote the file locally and left
  it there until the next sync tick, so the URL it returned 404ed for the first minute.

- **Alerts parse full markdown, and a Mermaid SVG is no longer capped.** Alert bodies rendered
  their markup as literal text, and the diagram was clamped to a width narrower than its own
  viewBox.

- **Null lists no longer break space switching in the CLI.** An empty response decoded to a nil
  slice and every caller assumed a list.

## [0.31.0] — 2026-08-27

### Added

- **Native Mermaid diagram rendering in web markdown viewer.** Integrated client-side Mermaid flowchart rendering with syntax error fallbacks and a source/diagram toggle in `MarkdownMuse.svelte`.
- **`publish_artifact` MCP tool.** Registered as the primary MCP tool for recording artifacts, retaining `publish_report` as a compatibility alias.

### Changed

- **Renamed reports to artifacts across the full stack.** Reorganized `internal/reports` into `internal/artifacts`, mounted `/api/artifacts` HTTP endpoints with `/api/reports` compatibility aliases, and updated web routes to `/artifacts`.

### Added

- **Web dashboard reports viewer and API endpoints.** Added `/api/reports` and `/api/reports/{id}` in the Go server, plus a full Reports view (`/reports` and `/reports/[id]`) in the SvelteKit dashboard with `@facile/muse` components, preview/source toggle, and deletion support.

## [0.29.0] — 2026-08-27

### Added

- **`people/` memory directory and person classification.** `mycelium init` scaffolds `memory/people/` alongside `standards/`, `classifyDir` routes person-related findings to `people/`, and the web dashboard displays `people` and `standards` folders.

## [0.28.0] — 2026-08-26

### Added

- **`mycelium report` records a rendered page and carries it to the machine you are sitting
  at.** An agent working over SSH on a headless box can produce an HTML report and has
  nowhere to put it. `mycelium report add <file>` writes it into `~/.mycelium/reports/`,
  which syncs like the rest of the tree, so it opens on the machine with a screen. It is
  stored and never hosted: a file opened from disk gets an opaque origin and cannot read the
  bearer token the web UI keeps in `localStorage`, which is the whole reason there is no URL
  to hand anybody. The identifier comes from the document's `<title>`, so recording the same
  page twice replaces it rather than piling up copies, and a report expires after thirty days
  unless `--expires never` pins it. Metadata rides in the file as meta tags rather than in a
  sidecar, because sync reconciles one file at a time and a shared manifest would conflict
  every time two machines recorded a report in the same minute. `add` names any relative
  `src` or `href` it finds, since a page opened from disk cannot fetch its siblings and a
  missing stylesheet is a failure only the reader ever sees.

- **A report never enters the wiki.** Every indexer, embedder and language scan is already
  scoped to `memory/**.md`, so `reports/` crosses machines without joining the corpus a
  search answers from. The wiki holds the fact as text and keeps it; a report is the picture
  of one and expires.

- **`publish_report`, a fourth MCP tool.** Its description carries the trigger as well as the
  behaviour, because that string is the only text guaranteed to be in front of a model at the
  moment it decides. It says what a report is not, twice: not a link, and not a finding.

### Changed

- **`/api/sync/files/{path}` answers `application/octet-stream` with `nosniff`.**
  `http.ServeFile` types a response from the file's extension, so the endpoint would have
  answered a synced `.html` as `text/html` on the origin whose `localStorage` holds the token
  that mints API tokens and writes `rules/`. Auth there is header-only, so a browser
  navigation gets a 401 and the path was never reachable that way. Presetting the type is
  what keeps that true now the tree carries agent-authored HTML, and it leaves ServeFile's
  Range and If-Modified-Since handling intact because ServeFile only sniffs when the header
  is absent.

## [0.27.0] — 2026-08-26

### Added

- **Agents reach memory and flows as tools, not as instructions.** `mycelium mcp` speaks
  MCP over stdio and serves three of them: `search_memory`, `list_flows` and `run_flow`.
  `run_flow` takes a flow name and nothing else, because a fixed set of procedures a human
  read and pinned is the whole defence against the injection this kind of server usually
  ships with, and an unpinned flow comes back with the command that pins it rather than an
  exit code. Stdio is not a default here but the only transport that can work: the steps run
  `sh -c` on the calling machine and the trust pin is per machine.

- **`mycelium install` declares that server in each assistant's own config**, merging into
  `~/.claude.json`, `~/.gemini/settings.json` and `~/.config/opencode/opencode.json` rather
  than overwriting them, and replacing the entry mycelium owns whatever name it carries so
  the pre-rename `jardin` entry stops pointing at a deleted binary. Codex is skipped on
  purpose: its MCP servers live in TOML. A rule now renders for one audience or the other,
  so an assistant is told to call `search_memory` or to run `mycelium memory search`, never
  both.

- **`mycelium memory add` files a finding and everything that goes with it in one call**: the
  finding on its page, that page's `updated:` stamp, the pointer in `index.md` and the line
  in `log.md`. All four or none. Done by hand the one that gets skipped is always the index,
  because it is the only edit that is not where the writing happened. `--body-stdin` takes
  prose that would otherwise be a fight with shell quoting, and the log line is an argument
  rather than something derived from the diff, because it records what was wrong before and
  a diff does not carry that.

- **`mycelium doctor` reports the MCP declaration.** An assistant whose config will not parse
  is refused rather than overwritten, its rules quietly render for the CLI, and the result
  looks exactly like an assistant that was never meant to have the tools. The check names
  the file and fails, and it separates that from an assistant install has simply not reached
  yet.

- **A failed background sync announces itself on Linux.** The generated systemd unit carries
  `OnFailure=`, where before a nonzero exit only marked the unit failed in systemd's own
  state and nothing read it until somebody typed `journalctl`. launchd already redirected to
  `~/.mycelium/daemon.log` and needs no equivalent.

- **One mark, generated rather than hand-edited.** `scripts/icons.ts` writes all six brand
  assets from a single glyph. They had drifted into three different marks, and the card that
  renders wherever the site is shared still advertised a hostname that answers 404.

### Changed

- **`mycelium sync --force` and `mycelium pull` only run from a terminal.** Both hand a
  machine the power to empty its own wiki in one call, and on 2026-08-25 a reconcile deleted
  102 files here because the guard that stopped it printed the flag that waived it. An agent
  whose goal is a finished sync reads that sentence as the next step. The refusal a
  non-interactive caller gets now carries no flag at all, and a refused bulk delete states
  how many files each side holds, which is the line that separates a cleanup somebody meant
  from a server that lost its data volume.

- **`mycelium claim start` asks the server before taking a claim.** Local claim files arrive
  on a daemon tick, so two agents on two machines could both read "no claims" for up to a
  minute and both take the repo. An unreachable server downgrades the verdict instead of
  ending the command: the claim is taken and reported as unverified, because a lock that
  fails closed would stop the work every time a laptop was on a train.

- **Sync staleness is measured in daemon ticks, not in a round number of hours.** The daemon
  runs every 60 seconds and the old threshold was 24 hours, a gap of 1440x, so sync failed
  for twenty hours on 2026-08-25 while `doctor` printed a green tick beside
  `last sync: 19h43m36s ago`. It is thirty ticks now when the service is installed and still
  a day when syncing is something a person does by hand, and `mycelium recap` prints the
  same line at the start of a session, only when it is stale.

- **A browser session from a password login expires and evicts the one it replaces**, which
  is what every SSO session already did. Until now that path minted an admin credential that
  never expired, which is also why it slipped past the fix below.

### Fixed

- **Login sessions are no longer listed as API tokens.** The settings page compared a whole
  token name against `session`, which matches neither `session:<email>` nor `cli:<email>`, so
  every browser and CLI login anyone had ever made appeared as an API token with a revoke
  button beside it. The filter moved to the server, where every client gets it right at once,
  and it tests the expiry rather than the name.

- **A page being written can no longer reach the server half-landed.** Every wiki write is
  staged beside its target and renamed into place, but the staged file was named so that sync
  did not recognise it as temporary. A crash between the two would have published a whole
  second copy of the page to the server and to every other machine.

- **Two agents filing a finding at the same moment no longer overwrite each other.**
  `index.md` and `log.md` are read, appended to and written back, and an agent files a
  finding when its task ends, which is exactly when several finish at once. The read and the
  write now happen under one lock. Sync deliberately does not take it: serialising a
  reconcile against a write would put a network round trip between an agent and its own
  finding.

- **An index pointer is named after the page it points at.** It carried the page's slug where
  every hand-written line in the corpus carries a title, so a filed entry stopped being
  scannable among the others.

- **`mycelium install claude` cannot leave the account record truncated.** `~/.claude.json`
  holds project history and the account record, is not regenerable, and is rewritten by any
  running session. It is now replaced through a rename rather than opened with `O_TRUNC`,
  which closes the window where a crash lost all of it.

## [0.26.0] — 2026-08-24

### Changed

- **Breaking.** The name is now Mycelium, and every user-visible surface moved with it. The
  binary is `mycelium`, the data directory is `~/.mycelium/`, the credentials file is
  `~/.mycelium.yml`, the environment variables are `MYCELIUM_URL`, `MYCELIUM_TOKEN`,
  `MYCELIUM_SERVER_URL`, `MYCELIUM_USAGE_TOKEN` and `MYCELIUM_VECTOR_SEARCH`, the Go module
  path is `github.com/FacileStudio/Mycelium`, and the server answers on
  `mycelium.facile.studio`. There is no compatibility shim: a machine set up before this
  release has to reinstall and log in again, and a self-hosted deployment has to move its
  data volume and update its OIDC application slug and redirect URL.

## [0.25.0] — 2026-08-24

### Added

- The daemon consolidates episodes into memory. A new stage reads recent records from
  `events/<agent>/*.jsonl`, proposes candidate findings with deterministic heuristics,
  asks a local Ollama model whether each is still useful in 30 days (failing open to a
  heuristic-only verdict when it is unreachable), applies the storage gate, dedupes against
  the wiki via embeddings with a lexical fallback — superseding an existing claim only when
  similarity, judge agreement and a newer date all agree, failing closed to NOOP otherwise —
  and writes the survivors directly into `memory/` using the prose conventions: findings
  appended, old claims struck through rather than deleted. A watermark cursor makes runs
  idempotent; the stage runs at most once per hour and works fully offline.

- Nacelle transcripts count toward usage. `sessions scan` tails
  `~/.nacelle/sessions/*.jsonl` alongside the Claude ones, reading only `turn` records
  because `done` repeats the run total and would double every figure.

- Episodes can stay on one machine. Anything under `local/` is excluded from sync, so a
  harness can log raw conversation text to `local/events/<agent>/` and have consolidation
  read it without the text ever reaching the server.

### Fixed

- Consolidation no longer overwrites a wiki page. A new finding whose title kebabbed onto
  an existing filename replaced that page whole, including anything a human had written in
  it — the deduper only reports a match above its weak-similarity floor, so a page can
  exist that the decision knows nothing about. A collision appends now, like every other
  write in the stage. Superseding also strikes each paragraph separately, because one `~~`
  pair around a two-paragraph claim renders as half retracted and half still standing.

- A candidate carrying a credential no longer reaches the judge. The gate ran after the
  round trip, and `OLLAMA_URL` can name a host that is not this machine, so the secret rule
  fired only once the text had already left the process.

- The same event line always produces the same text. Object keys were walked in map order,
  so a record carrying two of `message`/`content`/`text` assembled differently on every run
  — and the cursor hash, the similarity score and the create-or-noop decision all read it.

- Half-written pages stay off the wire. A crash between writing `<page>.md.tmp` and
  renaming it left the fragment in `memory/` forever, and sync would have published it to
  every machine and indexed it as a whole page.

- A Nacelle record written without its trailing newline is counted. The offset stayed
  before it and the file was skipped for good once its size stopped changing; a torn write
  still waits for the rest, since it cannot parse as whole JSON.

- A failed consolidation run stamps its timestamp, so a broken Ollama or an unwritable
  memory directory no longer re-runs the whole pipeline on every 60-second daemon tick.
  `--force` still bypasses the wait.

## [0.24.0] — 2026-08-24

### Added

- Normative pages are ratified per machine. A page carrying `type: standard` is
  pinned in `~/.mycelium/.memory-ratified.json` at the checksum a human accepted,
  with the machine and the day. `mycelium memory ratify` accepts one, `forget`
  closes out a deletion, `mycelium doctor` fails on a changed or missing page, and
  `mycelium memory search` marks a result from a changed page. The pin never syncs
  by design: a travelling pin would let one machine accepting a wrong edit clear
  the flag everywhere, which is the propagation the check exists to catch.

### Fixed

- A page name quoted in backticks or inside a fenced block is no longer read as
  a wiki link, so a page documenting the syntax stops crediting everything it
  quotes. The eval helper had always dropped fences; the parser had not.

## [0.23.0] — 2026-08-24

### Added

- A wiki link now carries a query-conditioned credit into memory ranking, so a
  page linked from a strong hit is pulled up with it.
- `mycelium doctor` reports a golden set the retrieval eval cannot run, instead of
  the eval quietly grading nothing.
- The eval fixture corpus doubled, with hard cases and a link case set, and it
  grades both search paths.

### Fixed

- `mycelium flow` compares the work dir through `EvalSymlinks`, so the gate passes
  on macOS, where `/tmp` is a symlink.
- An eval-set name pointing outside the wiki counts as missing rather than
  passing.
- The eval's corpus floor is tied to the doctor models, so the two cannot drift
  apart.

### Removed

- The cross-language eval set. The corpus is English by decision, so grading
  French queries against it measured nothing.

## [0.22.0] — 2026-08-23

### Added

- Memory has a history it can be restored from: a journal of what changed, per
  finding.
- A search result shows how old its claim is, and ranking weighs how current a
  claim is rather than how well it reads.
- `mycelium doctor` reports a history that has stopped recording.

### Fixed

- Strikethrough inside backticks is no longer stripped, and `node_modules` never
  syncs.

### Changed

- The container image carries the version stamp, without `--dirty`.

## [0.21.0] — 2026-08-23

### Fixed

- Conflict copies are kept out of the pages an agent reads, so a reconcile no
  longer feeds both sides of a conflict into a context window.

### Changed

- Sync prunes excluded directories instead of walking into them.

## [0.20.0] — 2026-08-23

### Added

- A sync refuses a reconcile that would destroy more than ten files. Bulk
  deletion is now a decision, not an accident.
- Memory parses a per-finding metadata block under each heading.

### Fixed

- `mycelium doctor` fails the last-sync check once it goes stale, rather than
  passing on any timestamp at all.

### Changed

- Git hooks moved from tracked scripts to lefthook.

## [0.19.0] — 2026-08-22

### Changed

- Every agent that follows the AGENTS.md standard is served from one generated
  tree, instead of one adapter per agent.
- Gemini skills are written as files rather than inlined into `GEMINI.md`.

## [0.18.0] — 2026-08-20

### Changed

- Auth, the OIDC callback, the alert scan, the timeline, the usage read and the
  sync reconcile were untangled, each behind tests written first. No behaviour
  change was intended.

## [0.17.0] — 2026-08-20

### Added

- Retrieval has a gate that runs everywhere, graded against a committed corpus.
- Tests cover the flows and models routes, which shipped without any.

### Changed

- The oversized files were split: `internal/server` went from seven files to
  twenty-three, and the client and the rest of the Go tree followed.

## [0.16.1] — 2026-08-19

### Fixed

- A model's closure covers `require()` and backticked specifiers, not only
  static `import` strings.

## [0.16.0] — 2026-08-19

### Changed

- Pinning a model pins its imports, not just the file it names. A trusted model
  can no longer change by editing something it pulls in.

### Fixed

- The retrieval eval is skipped when its corpus is not on this machine, instead
  of failing the gate.

## [0.15.4] — 2026-08-19

### Fixed

- A model that takes no arguments receives an empty object rather than nothing.

## [0.15.3] — 2026-08-19

### Changed

- Rules and Skills are one Instructions entry in the dashboard rail.

## [0.15.2] — 2026-08-19

### Changed

- Flows and Models are one Automation entry in the dashboard rail.

## [0.15.1] — 2026-08-19

### Changed

- Flow and model handlers moved out of `server.go`.

## [0.15.0] — 2026-08-19

### Added

- A step can declare which of its values are secret, and they are redacted from
  the run artifact.
- `defineModel()` gives a model file the contract, so it stops being written by
  hand.
- The dashboard shows flows and models, read-only.

## [0.14.1] — 2026-08-19

### Fixed

- An ephemeral value stays ephemeral once something consumes it. It used to leak
  into the consumer's artifact.

## [0.14.0] — 2026-08-19

### Added

- The flow roadmap is finished: typed steps that run a model extension, a
  dependency graph via `depends_on`, and run history.
- The CLI-standard pass is complete, and the document it was measured against
  was corrected where the two disagreed.

## [0.13.2] — 2026-08-19

### Fixed

- A sync with nothing else to do prunes the stale backups, so they stop
  accumulating on an idle machine.

## [0.13.1] — 2026-08-19

### Fixed

- A machine no longer conflicts with itself over its own telemetry.

## [0.13.0] — 2026-08-19

### Added

- `needs` binds an environment variable to an earlier step's `stdout`, `stderr`
  or `exit_code`. Backward references only, and the value is data, never spliced
  into the command string.

## [0.12.0] — 2026-08-19

The largest release in the history: flows, memory ranking, semantic search and
the French-language check all landed here.

### Added

- `mycelium flow`: recorded procedures, pinned before they execute, with `flow
  add` to scaffold one, descriptions in `flow list`, and this machine's flows in
  the session recap. The trust gate is reviewable, and one bad flow no longer
  hides the rest.
- Memory search ranks with BM25 against a golden set, retrieves finding blocks
  rather than whole pages, and stops caring about word order. Ranking is
  reproducible and a pin can be dropped.
- Semantic search: embeddings via ollama, vectors stored flat or in qdrant, the
  server's lexical half ranked by chunk, and search from the dashboard. The
  index heals itself on a timer and indexes what is already there at worker
  start.
- A sync warns when it carries French, and never blocks on it. The check has one
  implementation, with a measured floor and a statement of what it cannot catch.
- Accent folding, so a French wiki answers an unaccented query.
- A claims coordination layer: task leases with scratchpads, injected in the
  recap, with a claims API, dashboard UI, an opencode adapter and a codex recap
  hook.

### Changed

- Install targets the agents this machine has, not every adapter.
- Landing and auth buttons harmonized on rounded corners and an SSO-first
  layout.

### Fixed

- `doctor` detects opencode, which it could not see before.
- Provenance is stripped from server excerpts too, not only local ones.

## [0.11.1] — 2026-08-11

### Added

- `mycelium doctor`, `mycelium skills validate`, and source frontmatter on skills.
  Feature work shipped under a patch bump.

## [0.11.0] — 2026-08-11

### Added

- The sessions dashboard gains tabs and cost tracking, plus a filet CI job.
- Monochrome brand logos in the landing adapters grid, and pi listed among the
  adapters.

### Fixed

- Grouping by model no longer shows `(none)` rows; cost model matching, stat
  card order and grid responsiveness corrected.

## [0.10.1] — 2026-08-11

### Added

- A canonical event store: the pi extension writes, Mycelium collects. Feature
  work shipped under a patch bump.
- The bus is configurable from the environment, and the variables are
  documented.

### Fixed

- A restarting server no longer produces "you are not an admin".
- `TRUSTED_PROXIES` is passed into the router instead of being inert, and
  Cloudflare is trusted as a proxy, so rate limits key on the visitor again.
  Picks up tronc v0.10.1, which stops trusting any `X-Forwarded-For`.

### Changed

- Favicons, icons, OG images and metadata harmonized across the suite.

## [0.10.0] — 2026-08-10

### Added

- The CLI signs in through the browser over a loopback SSO flow, with `logout`
  and environment overrides.

### Changed

- Installation delegates to the `facile` CLI, bootstrapped from
  `get.facile.studio`.

## [0.9.1] — 2026-08-08

### Security

- Admin scope is re-derived live instead of being trusted from a frozen session
  token, so revoking admin takes effect without waiting for the token to expire.

### Fixed

- Usage alerts key on the account, not the machine.

## [0.9.0] — 2026-08-08

### Added

- An alert fires on the Antenne when a usage window crosses its threshold.

## [0.8.0] — 2026-08-08

### Changed

- Breaking: Nook is renamed Antenne throughout.
- The server runs on tronc v0.8.0 and serves the client from one container.
- The client is rebuilt on the muse design system, with muse owning the anchors,
  the empty state and the sidebar collapse button.

### Added

- The dashboard charts subscription usage limits and tracks them, with the
  timeline and usage surface documented.
- The API reference is served at `/docs` through `tronc/apiref`.
- A quality gate runs before every push, and the curl installer gains
  `--no-color`.

### Fixed

- The pool emitter stops listening to clients it has already retired.
- The OIDC callback parses the error envelope and validates the token before
  redirecting.
- The HTTP server recovers from handler panics instead of dropping the process.

## [0.7.0] — 2026-08-04

### Added

- Live sessions: see what every machine is working on right now.

## [0.6.1] — 2026-08-04

### Security

- The common tree is owner-private rather than world-readable.

### Fixed

- Machine tokens are tied to users, so a machine can sync spaces.

## [0.6.0] — 2026-08-04

### Added

- Spaces, and Authentik SSO.

## [0.5.2] — 2026-08-04

### Added

- Canonical project identity, and emit hygiene on the event path.

## [0.5.1] — 2026-08-04

### Fixed

- Session scans take a `flock`, and shard reads are deduped by block id.

## [0.5.0] — 2026-08-04

### Added

- Session tracking: collect what an agent did, recap the context, publish it to
  Nook.

## [0.4.1] — 2026-08-03

### Fixed

- The daemon writes a stable binary path into its service unit, so an upgrade
  does not orphan the unit.

### Changed

- Leaf logo and favicon; the beehive metaphor is gone.

## [0.4.0] — 2026-08-03

### Changed

- Breaking: Ruche is renamed Mycelium.
- The README documents self-update and dashboard memory editing.

### Fixed

- Agent config distribution uses dir-based detection, logs a skipped sync, and
  writes to the Claude skills dir.

## [0.3.0] — 2026-06-23

### Added

- A memory editor in the dashboard, and `ruche update` for self-update.

## [0.2.0] — 2026-06-23

### Added

- `login` runs a browser device-authorization flow.

## [0.1.0] — 2026-06-23

First release, built and renamed in a single day: Hive became Ruche partway
through, and the on-disk vocabulary settled on `memory`.

### Added

- The CLI, a shared agent memory that syncs over HTTP, with a three-way
  reconcile so a concurrent edit does not silently lose data.
- A SvelteKit dashboard: machines with live connectivity, a memory file-tree
  browser and viewer, rules and skills as cards with detail pages, a landing
  page, and a copy-paste agent integration prompt.
- A background sync service, auto-enabled on login, that detects agents itself
  and tracks a token's last-seen time.
- A codex adapter that emits `~/.codex/skills/<name>/SKILL.md` instead of
  inlining the skill.
- YAML config at `~/.ruche.yml`, persisted tokens, and `ruche login`.
- GoReleaser v2 releases with a Homebrew tap.

### Security

- Tokens are hashed at rest and scoped, login is rate-limited, and the browser
  session token is hidden from the API tokens list and from the machines page.
  Session tokens are named per machine and rotate on re-login.
- Adapter output refuses to write through a symlink.

### Changed

- Cells are gone; there is one flat data dir. `sync_url` and `sync_token` are
  `url` and `token`, and the `RUCHE_` env prefix is dropped.

### Fixed

- `rule_order` is a hint, not an allowlist.
- A 401 clears the stored token, which breaks a dashboard redirect loop.
- Traefik gives API routes priority over the client, and the Docker build uses
  `golang:alpine`.

[Unreleased]: https://github.com/FacileStudio/Mycelium/compare/v0.33.0...HEAD
[0.33.0]: https://github.com/FacileStudio/Mycelium/compare/v0.32.3...v0.33.0
[0.32.3]: https://github.com/FacileStudio/Mycelium/compare/v0.32.2...v0.32.3
[0.32.2]: https://github.com/FacileStudio/Mycelium/compare/v0.32.1...v0.32.2
[0.32.1]: https://github.com/FacileStudio/Mycelium/compare/v0.32.0...v0.32.1
[0.32.0]: https://github.com/FacileStudio/Mycelium/compare/v0.31.0...v0.32.0
[0.31.0]: https://github.com/FacileStudio/Mycelium/compare/v0.29.0...v0.31.0
[0.29.0]: https://github.com/FacileStudio/Mycelium/compare/v0.28.0...v0.29.0
[0.28.0]: https://github.com/FacileStudio/Mycelium/compare/v0.27.0...v0.28.0
[0.27.0]: https://github.com/FacileStudio/Mycelium/compare/v0.26.0...v0.27.0
[0.23.0]: https://github.com/FacileStudio/Mycelium/compare/v0.22.0...v0.23.0
[0.22.0]: https://github.com/FacileStudio/Mycelium/compare/v0.21.0...v0.22.0
[0.21.0]: https://github.com/FacileStudio/Mycelium/compare/v0.20.0...v0.21.0
[0.20.0]: https://github.com/FacileStudio/Mycelium/compare/v0.19.0...v0.20.0
[0.19.0]: https://github.com/FacileStudio/Mycelium/compare/v0.18.0...v0.19.0
[0.18.0]: https://github.com/FacileStudio/Mycelium/compare/v0.17.0...v0.18.0
[0.17.0]: https://github.com/FacileStudio/Mycelium/compare/v0.16.1...v0.17.0
[0.16.1]: https://github.com/FacileStudio/Mycelium/compare/v0.16.0...v0.16.1
[0.16.0]: https://github.com/FacileStudio/Mycelium/compare/v0.15.4...v0.16.0
[0.15.4]: https://github.com/FacileStudio/Mycelium/compare/v0.15.3...v0.15.4
[0.15.3]: https://github.com/FacileStudio/Mycelium/compare/v0.15.2...v0.15.3
[0.15.2]: https://github.com/FacileStudio/Mycelium/compare/v0.15.1...v0.15.2
[0.15.1]: https://github.com/FacileStudio/Mycelium/compare/v0.15.0...v0.15.1
[0.15.0]: https://github.com/FacileStudio/Mycelium/compare/v0.14.1...v0.15.0
[0.14.1]: https://github.com/FacileStudio/Mycelium/compare/v0.14.0...v0.14.1
[0.14.0]: https://github.com/FacileStudio/Mycelium/compare/v0.13.2...v0.14.0
[0.13.2]: https://github.com/FacileStudio/Mycelium/compare/v0.13.1...v0.13.2
[0.13.1]: https://github.com/FacileStudio/Mycelium/compare/v0.13.0...v0.13.1
[0.13.0]: https://github.com/FacileStudio/Mycelium/compare/v0.12.0...v0.13.0
[0.12.0]: https://github.com/FacileStudio/Mycelium/compare/v0.11.1...v0.12.0
[0.11.1]: https://github.com/FacileStudio/Mycelium/compare/v0.11.0...v0.11.1
[0.11.0]: https://github.com/FacileStudio/Mycelium/compare/v0.10.1...v0.11.0
[0.10.1]: https://github.com/FacileStudio/Mycelium/compare/v0.10.0...v0.10.1
[0.10.0]: https://github.com/FacileStudio/Mycelium/compare/v0.9.1...v0.10.0
[0.9.1]: https://github.com/FacileStudio/Mycelium/compare/v0.9.0...v0.9.1
[0.9.0]: https://github.com/FacileStudio/Mycelium/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/FacileStudio/Mycelium/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/FacileStudio/Mycelium/compare/v0.6.1...v0.7.0
[0.6.1]: https://github.com/FacileStudio/Mycelium/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/FacileStudio/Mycelium/compare/v0.5.2...v0.6.0
[0.26.0]: https://github.com/FacileStudio/Mycelium/compare/v0.25.0...v0.26.0
[0.5.2]: https://github.com/FacileStudio/Mycelium/compare/v0.5.1...v0.5.2
[0.5.1]: https://github.com/FacileStudio/Mycelium/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/FacileStudio/Mycelium/compare/v0.4.1...v0.5.0
[0.4.1]: https://github.com/FacileStudio/Mycelium/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/FacileStudio/Mycelium/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/FacileStudio/Mycelium/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/FacileStudio/Mycelium/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/FacileStudio/Mycelium/releases/tag/v0.1.0
