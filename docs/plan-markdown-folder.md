# Plan: Mycelium as a markdown folder, plus artifact search

Status: proposed, 2026-11-14. Branch `refactor/drop-vector-search`, which already holds commit
677db62 (the vector-search code removal, gates green). Nothing is merged, tagged or deployed.

## Goal

Make Mycelium a plain markdown folder with keyword search, add one artifact search over titles
and content, and give the Mycelium and Bulle artifact pages a card grid and a real search bar.

## Why

- The owner: "we don't need this fancy search shit anymore ... we don't need a vector db". The
  vector half never ran in prod: no ollama container answers, the boot log reads
  `lookup ollama ... no such host` (projects/mycelium.md, v0.46.0 finding), and the prod data
  volume holds no `.embeddings` directory (direct observation, 2026-11-14).
- Artifact search today matches titles only (`titleHolds` in `internal/server/mcp_artifacts.go`;
  the web page filters the loaded list client-side in
  `apps/client/src/routes/(app)/artifacts/+page.svelte`). Nobody can find an artifact by what
  it says.
- The Mycelium artifact grid is fixed at 1/2/3 columns (`sm:grid-cols-2 lg:grid-cols-3`). The
  Bulle artifact page is a month-grouped list with the `icons.layers` glyph.

## Research findings

What the field does in 2026, and what it changes in this design.

- **OKF is the Open Knowledge Format**, Google Cloud's open spec (v0.2) "that formalizes the
  LLM-wiki pattern into a portable, interoperable format" (blog, 2026-06-12:
  https://cloud.google.com/blog/products/data-analytics/how-the-open-knowledge-format-can-improve-data-sharing).
  A bundle is "a directory of markdown files with YAML frontmatter. There is no schema registry,
  no central authority, and no required tooling"; storage, serving and query infrastructure are
  explicit non-goals (https://raw.githubusercontent.com/GoogleCloudPlatform/open-knowledge-format/main/SPEC.md).
  `type` is the only required key; `title`, `description`, `tags`, `sources`, `status`
  (`draft | stable | deprecated`), `stale_after` and `verified` are optional. `index.md` and
  `log.md` are reserved names: `index.md` lists a directory's concepts as
  `* [Title](path) - description` for progressive disclosure; `log.md` is date-grouped
  (`## YYYY-MM-DD`), newest first. Links are plain relative markdown links.
- **Karpathy's LLM wiki** (https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f):
  an agent-maintained folder of interlinked markdown with `index.md` (content catalog, read
  first on every query) and `log.md` (append-only, `## [YYYY-MM-DD] op | title` lines). The index
  "works surprisingly well at moderate scale (~100 sources, ~hundreds of pages) and avoids the
  need for embedding-based RAG infrastructure". Search tooling is optional: "at small scale the
  index file is enough"; qmd (BM25 plus vectors) is suggested only "as the wiki grows".
- **Claude Code memory** (https://code.claude.com/docs/en/memory): CLAUDE.md files plus auto
  memory, plain markdown under `~/.claude/projects/<project>/memory/`, `MEMORY.md` loaded at
  start (first 200 lines or 25KB), topic files read on demand. No index service.
- **AGENTS.md** (https://agents.md/): "a README for agents", plain markdown, nested files where
  "the closest AGENTS.md to the edited file wins".
- **Manus** (https://manus.im/blog/Context-Engineering-for-AI-Agents-Lessons-from-Building-Manus):
  "we treat the file system as the ultimate context ... unlimited in size, persistent by nature,
  and directly operable by the agent itself"; compression stays restorable by keeping paths.
- **Anthropic, Claude Managed Agents memory** (2026-04-23,
  https://claude.com/resources/articles/claude-managed-agents-memory): memory "mounts directly
  onto a filesystem", memories are files with an audit log and rollback.
- **Letta** (https://www.letta.com/blog/benchmarking-ai-agent-memory): a file-backed agent with
  grep and search tools scored 74.0% on LoCoMo against Mem0's reported 68.5%; "memory is more
  about how agents manage context than the exact retrieval mechanism". Caveat: Letta's files were
  also embedded for a `search_files` tool, so this is not a pure grep-only result.

Design consequences:

1. Keep the markdown tree canonical and keep BM25 keyword search (already in
   `internal/memory/search_chunks.go`). Drop vectors.
2. Keep `index.md` and `log.md`: Mycelium's shape already matches Karpathy and OKF. The `log.md`
   order differs (Mycelium appends oldest first, OKF lists newest first); leave it, it is
   append-only by rule and a reorder buys nothing.
3. Mycelium frontmatter (`title`, `type`, `sources`, `related`, `confidence`, `created`,
   `updated`) already covers OKF's required `type` and most optional keys. Do not rename keys;
   an OKF export can map `related` to links later if anyone asks (YAGNI).
4. Dead ends: none of the sources gives a benchmark of BM25-only against hybrid on a wiki of
   this size. The repo's own fixture eval (`internal/memory/fixture_eval_test.go`, recall@5 1.000
   on 70 pages) is the evidence that lexical search suffices here.

## Checked against

conventions/facile-suite-rules (owner rules: no borders, MCP free, vitrine and Field off
limits), projects/mycelium (deploy: compose KFh8c0GZjIhwhnASe3CDs, VERSION file, backups,
webhook misses), conventions/mycelium-agent-surface, syntheses/agent-memory-practice-2026, the
repo CLAUDE.md, `scripts/check.sh` header, `filet.yml`, `mycelium flow list`.

## Approach

Lane A finishes the removal already on the branch, then adds `internal/artifacts/search.go`, one
pure function over the artifact files that both a REST route and a new MCP tool call, scoped by
the same guards the list routes use today (`artifactRoot`, `artifactCommonAllowed`,
`memberSpaces`, `mcpScope`). Lane B rebuilds the Mycelium artifacts page on that endpoint. Lane
C adds a Bulle proxy route over the new MCP tool and rebuilds Bulle's artifacts page the same
way. The contract below is fixed up front so B and C can build against it while A is in review.

## Contract (fixed before the lanes start) [contract]

REST, Mycelium: `GET /api/artifacts/search`

- Query: `q` (required, trimmed, 1 to 200 characters), `space_id` (one tree, through the
  `artifactRoot` guard) or `all=true` (common tree when `artifactCommonAllowed`, plus every
  member space), `cursor` (opaque, from `next_cursor`), `limit` (1 to 50, default 20).
- Match: `q` is lowercased and split on whitespace; an artifact matches when every word occurs
  in its title or its body (frontmatter stripped with `stripFrontmatter`). Case-insensitive
  substring, no accents folding (same rule as `titleHolds` today).
- Order: artifacts whose title holds every word first, then content-only matches; inside each
  group, newest `updated` first, then id.
- Answer: `{"results":[{id,title,machine,created,updated,space_id,space_name,match,snippet}],
  "total":N,"next_cursor":"..."}`. `match` is `title` or `content`. `snippet` is about 160
  characters of body around the first word found, whitespace collapsed, `…` at a cut; for a
  title-only match, the body's first 160 characters.
- Errors: 400 for an empty or oversize `q` and a bad cursor, 403 from the existing guards.

MCP, Mycelium: `mycelium_search_artifacts`

- Input: `query` (required), `space_id` (empty searches every tree in `mcpScope`), `cursor`.
- Output: `{results: [mcpArtifact + match + snippet], total, next_cursor}`, pages of 20.
- Annotations: read-only, idempotent, closed world, like `mycelium_list_artifacts`.

REST, Bulle: `GET /api/artifacts/search?q=&cursor=` answering
`{artifacts: Artifact[] (plus snippet, match), next_cursor, available}`, the same envelope as
`GET /api/artifacts`.

## Lanes

Lanes A, B and C run in parallel; each owns its files. Releases run in order A, B, C.

### Lane A: Mycelium backend (one agent, branch `refactor/drop-vector-search`)

Owns: `internal/`, `cmd/`, `docs/` (except this file), `README.md`, `SPEC.md`, `ROADMAP.md`,
`CHANGELOG.md`, `VERSION`, `docker-compose.yml`, `.env.example`. Never touches `apps/client/`
after 677db62.

Done in 677db62 (verified: `go build ./...`, `go vet ./...`, `go test ./...`, `bun run check`,
`bun run test`, `filet check` exit 0): deleted `internal/memory/{embed,ollama,store_flat,store_qdrant,client}.go`
and tests, `embed_eval_test.go`, `internal/server/embed_{queue,reconcile,status,worker}.go`,
`memory_index_api.go` and tests, `scripts/ai-sidecars.sh`,
`apps/client/src/lib/components/IndexStatus.svelte`; renamed `hybrid.go` to `search_chunks.go`
and dropped `FuseRRF`; `memory_api.go` and `mcp_memory.go` rank with `SearchChunks` only;
`degraded` left the REST, MCP and client types; `cmd/serve.go` and `internal/env` no longer read
`OLLAMA_URL`, `EMBED_MODEL`, `VECTOR_STORE`, `QDRANT_URL`; `mycelium memory search` lost
`--local`/`--verbose` and the remote path; `vector_search` and `MYCELIUM_VECTOR_SEARCH` are
gone; consolidate dedupe lost its embedding backend (the Ollama durability judge stays: it is an
LLM call, not search); compose lost the four env keys and the `facile-ai` network. 4623 lines
removed, 149 added, `go mod tidy` changed nothing (no vector dependency was in go.mod).

Steps left:

1. Docs: remove the vector, ollama, qdrant and embedding passages from `docs/deployment.md`
   (30 hits), `docs/configuration.md` (15), `docs/architecture.md` (8), `SPEC.md` (8),
   `ROADMAP.md` (7), `docs/api.md`, `docs/usage.md`, `docs/flow-composition.md`; add a one-line
   superseded note at the top of `docs/plan-consolidation.md` (a historical plan, not rewritten).
   Search reads as "BM25 over the markdown folder". Verify:
   `grep -n -i -E 'qdrant|ollama|embedding|vector|hybrid|minilm' README.md SPEC.md ROADMAP.md docs/*.md | grep -v plan-` -> only the superseded note and the consolidate judge's Ollama.
2. `CHANGELOG.md`: under `[Unreleased]`, `### Removed` (vector search, index status route,
   `degraded`, `--local`, `--verbose`, `vector_search`, env keys, `facile-ai` network) and a
   `BREAKING` line. Verify: `sh scripts/check.sh --go-only` -> exit 0.
3. `internal/artifacts/search.go` (new) and `search_test.go` (new): `Search(trees []Tree, query
   string) []Hit` per the contract; `Tree` is `{Root, SpaceID, SpaceName}`. Reuse
   `List`, `stripFrontmatter`. Test: title match outranks content match, AND across words,
   snippet window, frontmatter never matched. Verify: `go test ./internal/artifacts/` -> ok.
4. `internal/server/artifacts_search.go` (new), route in `internal/server/server.go` next to
   `r.Get("/artifacts", ...)`: `r.Get("/artifacts/search", s.auth(false, s.artifactsSearch))`.
   Scope with `artifactRoot` for `space_id` and the `artifactsListAll` logic for `all=true`
   (extract the tree list into one helper both call, do not copy it). Cursor with
   `mcpkit.EncodeCursor`. [auth: no new token type, no delegated allowlist change; Bulle comes in
   through MCP, so `delegatedRoutes` stays as is]. Test in `artifacts_search_test.go` (new):
   403 on a foreign `space_id`, `all=true` returns member spaces only, pagination. Verify:
   `go test ./internal/server/ -run ArtifactsSearch` -> ok.
5. `internal/server/mcp_artifact_search.go` (new), tool registered in `mcp_tools.go`
   (`toolSearchArtifacts = "mycelium_search_artifacts"`, mentioned in the server instructions
   string), input and output types in `mcp_types.go`. Test with the existing MCP fixture
   (`mcp_fixture_test.go`). Verify: `go test ./internal/server/ -run MCPSearchArtifacts` -> ok.
6. `internal/documentation/registry.go` and `types.go`: document the route. Verify:
   `go build ./...`.
7. Gate: `sh scripts/check.sh` then `filet check`, both exit 0.
8. Release v0.53.0 (breaking, minor on 0.x): bump `VERSION`, move `[Unreleased]` to
   `## [0.53.0] - <date>`; `scripts/check-version.sh` enforces the match. PR, merge, then check
   `git log -1` before `git tag v0.53.0` (the v0.47.4 tag slip). Push the tag before main.
   CI may not run (billing); record that and rely on the local gates. The pre-push hook runs
   the full check.sh: wait for the remote ref, do not re-push.
9. Deploy and verify [deploy]: back up the volume first
   (`docker run --rm --platform linux/amd64 -v compose-input-optical-alarm-5soe99_mycelium-data:/data:ro -v /home/yann/backups:/out alpine:3.20 tar czf /out/mycelium-prod-<stamp>.tar.gz -C /data .`),
   then `COMPOSE_ID=KFh8c0GZjIhwhnASe3CDs COMMIT=<release sha> VERIFY_URLS=https://mycelium.facile.studio/health mycelium flow run dokploy-deploy-wait`.
   Judge by container StartedAt, `/version` = 0.53.0 and the served bundle. Live checks:
   `POST /api/memory/search` returns results without `degraded`; `GET /api/artifacts/search?q=encre&all=true`
   returns hits with snippets; `mcp-live-smoke` with `MCP_TOOL=mycelium_search_artifacts`.
10. Local CLI on ruche: `facile update` (or the repo's `install.sh`) to `~/.local/bin/mycelium`,
    then `mycelium --version`, `mycelium memory search "dokploy compose"` and `mycelium doctor`
    all succeed.
11. Report to the owner (do not remove anything): Dokploy env keys the server no longer reads,
    `OLLAMA_URL`, `EMBED_MODEL`, `VECTOR_STORE`, `QDRANT_URL` (all four set on the running
    container today; `dokploy-env-safe` cannot drop keys, so a human removes them in the panel);
    unused docker resources: network `facile-ai` (only the Mycelium container joins it), volumes
    `ollama` and `qdrant_storage` (from `scripts/ai-sidecars.sh`); no `ollama` or `qdrant`
    container exists. Volumes `grimoire_qdrant_data` and `marcel-qdrant-ws2aqf_qdrant_data`
    belong to other projects: leave them out of the report as Mycelium's.
12. Wiki (end gate): in `projects/mycelium.md` mark the vector claims
    `[SUPERSEDED by: v0.53.0, <date>]` and add a v0.53.0 finding (search is BM25 only, the
    artifact search contract); update `syntheses/agent-memory-practice-2026.md` with the OKF and
    Karpathy findings above; append `log.md` lines.

### Lane B: Mycelium UI (one agent, worktree in `~/.cache/facile-wt/mycelium-ui`, branch from main)

Owns: `apps/client/` only. Load `/muse` first. [ui: Svelte 5, muse components and tokens, no
`border-*` or `divide-*`, focus is an outline ring, no hard-coded color].

1. `apps/client/package.json`: `@facile/encre` from `#v0.10.0` to `#v0.12.1` once the owner tags
   it (nowrap tables with horizontal scroll); `bun install`. Verify:
   `cd apps/client && bun run test` -> `scripts/check-encre.ts` passes.
2. `apps/client/src/lib/backend.ts` and `backendTypes.ts`: `artifactsAll()` (`GET /artifacts?all=true`),
   `artifactsSearch(q, cursor?)` (`GET /artifacts/search?all=true&q=`), types per the contract.
3. `apps/client/src/routes/(app)/artifacts/+page.svelte`: grid
   `grid grid-cols-[repeat(auto-fill,minmax(min(18rem,100%),1fr))] gap-4` (one column on a
   phone, about 18rem cards above). Each card: `solar:file-smile-linear` icon, title, date, space
   name (common tree shows the personal label). Reuse `$lib/components/EntityCard.svelte`.
   Search bar: muse `Input type="search"`, debounced 250 ms, value mirrored in `?q=` with
   `goto(..., { replaceState: true, keepFocus: true })`, read back on load; `/` focuses it unless
   focus is already in a field; Escape clears. Empty query shows the full list; a query calls
   `artifactsSearch` and shows the snippet under the title. States: loading `Spinner`, error
   `Alert`, no artifacts `EmptyState`, no match `EmptyState` naming the query. Load more with
   `next_cursor`.
4. `apps/client/src/routes/(app)/artifacts/[id]/+page.svelte` and `backend.ts`: a card links
   `/artifacts/<id>?space=<space_id>`; the page passes that space to `artifactGet`, so a card from
   another space opens. Verify: open a space artifact from the grid while the common tree is
   active.
5. `apps/client/src/lib/i18n/messages/automation.ts`: new `artifacts.*` keys in the 4 locales.
   Verify: `bun run test` -> "placeholders agree".
6. Gate: `sh scripts/check.sh` and `filet check` exit 0; `bun run check` 0 errors. Check at 360,
   390, 768, 1024 and 1440 px with headless chromium: no horizontal scroll, one column at 360.
7. Release after lane A is deployed: rebase on main, v0.54.0, same deploy and verify steps as
   A9, plus load `/artifacts` and one artifact holding a table on prod.

Depends on: the contract (to build), lane A deployed (to verify), encre v0.12.1 tag (step 1).

### Lane C: Bulle (one agent, worktree in `~/.cache/facile-wt/bulle-artifacts`)

Owns: `~/Code/Facile/Bulle` `apps/api/internal/suite/mycelium_artifacts*.go`,
`apps/api/internal/suite/myceliumtest/myceliumtest.go`, `apps/api/modules/artifacts/`,
`apps/client/src/lib/artifacts.ts` and its test, `apps/client/src/lib/suite-icons.ts`,
`apps/client/src/routes/(app)/artifacts/+page.svelte`,
`apps/client/src/lib/components/ArtifactCard.svelte`, the `/artifacts` nav entry in
`apps/client/src/routes/(app)/+layout.svelte`, `apps/client/package.json`. Follow Bulle's own
CLAUDE.md for its gates and release.

1. `suite/mycelium_artifacts.go`: `SearchArtifacts(ctx, userID, query, cursor)` calling
   `mycelium_search_artifacts`; add the tool to the `myceliumtest` fake. Verify: Bulle
   `sh scripts/check.sh` -> exit 0.
2. `modules/artifacts/router.go`, `handler.go`, `types.go`, `documentation.go`:
   `GET /artifacts/search` with `requireAuth`, the `disconnected` -> `available: false` path like
   `listArtifacts`. Test in `handler_test.go`.
3. `apps/client/src/lib/suite-icons.ts`: register the `solar:file-smile-linear` body by hand (the
   body copied from the Iconify Solar set; Mycelium's client already renders it through muse);
   switch the artifacts nav item, the list page and `ArtifactCard.svelte` from `icons.layers` to
   it. Leave the other `icons.layers` uses (memory tab, ToolGroups) alone.
4. Artifacts page: the same grid, card content and search bar as lane B step 3, backed by the
   Bulle route; drop `groupByMonth` from the page if nothing else uses it (grep first), keep the
   space `Select`. Remove `class="border-0"` overrides only if muse already ships borderless
   inputs.
5. `apps/client/package.json`: `@facile/encre` to `#v0.12.1`, `bun install`, `bun run test`.
6. Gate: Bulle `sh scripts/check.sh`, `filet check`, client `bun run check` and `bun run test`.
7. Release after Mycelium v0.53.0 is live (the MCP tool must exist in prod first), then deploy
   per Bulle's CLAUDE.md and verify the page loads, a search returns a content-only match, and an
   artifact table scrolls instead of wrapping.

Depends on: Mycelium v0.53.0 deployed (live), encre v0.12.1.

## Rule text to propose to the owner (do not edit the rules)

`~/.mycelium/rules/20-memory.md` line 50 and `~/.claude/CLAUDE.md` line 118 read:

> Wiki pages and queries are English. The embedding model (`all-minilm`) is English-only, so
> neither half of hybrid search crosses the language boundary.

Proposed replacement:

> Wiki pages and queries are English. Search is keyword matching (BM25) over the markdown files,
> and a keyword does not cross languages: a French page is invisible to an English query.

The same files also say "neither half of hybrid search": the replacement covers it. The `mycelium
doctor` `wiki language` check stays: it is a word-list heuristic in `internal/memory/language.go`,
not an embedding, and the keyword search needs one language just as much.

## Boundaries

- Ask first: the env keys and the `facile-ai` network leaving compose (lane A, already on the
  branch), the new REST route and MCP tool (contract), any change to `delegatedRoutes`.
- Never: raise a filet threshold, add to `disabled:`, bypass the pre-push hook, delete Dokploy
  env keys, containers or volumes, touch vitrine or Field, merge lane C before lane A is live.

## Done when

`sh scripts/check.sh` and `filet check` exit 0 in Mycelium and Bulle; prod `/version` is the
released version; `mycelium memory search` and `mycelium doctor` work on ruche; a content-only
word finds an artifact in both web apps; the grids show one column at 360 px.

## Risks

- Removing `degraded` and `source` from answers: Bulle decodes `degraded` with `omitempty`
  (`apps/api/internal/suite/mycelium.go`), so a missing key is false; no other reader found.
- Scripts passing `mycelium memory search --local` now fail on an unknown flag; a grep of
  `~/.config`, `~/.claude`, `~/.mycelium` and nacelle found none.
- Content search reads every artifact file per query (164 in the common tree on 2026-11-10);
  fine at this size, revisit past a few thousand.
- Dokploy webhooks miss pushes on this host; `dokploy-deploy-wait` reports a missing row as
  queued. Never deploy by hand after a push.
- GitHub Actions may be blocked by billing; local gates are then the only proof.

## Skip

- An OKF exporter or frontmatter rename: nobody consumes OKF bundles from Mycelium today.
- Removing the consolidate judge's Ollama call: it is a local LLM durability check, not search;
  propose separately if the owner wants Ollama gone everywhere.
- Reordering `log.md` newest first to match OKF.
