package docs

import "net/http"

var (
	unauthenticated = Error{Status: 401, Code: "unauthenticated", Description: "Missing, expired or revoked token."}
	forbidden       = Error{Status: 403, Code: "permission_denied", Description: "The token is valid but not admin-scoped."}
	notFound        = Error{Status: 404, Code: "not_found", Description: "No such record."}
	badRequest      = Error{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body or query."}
	internalError   = Error{Status: 500, Code: "internal", Description: "Unexpected server error."}
	managed         = Error{Status: 409, Code: "conflict", Description: "Registre manages this; the message links to it."}
	notMember       = Error{Status: 403, Code: "permission_denied", Description: "Not a member of the space."}
	tooLarge        = Error{Status: 413, Code: "resource_exhausted", Description: "Content exceeds 1 MiB."}
)

var (
	anyToken   = []Error{unauthenticated, internalError}
	adminToken = []Error{unauthenticated, forbidden, internalError}
)

var nameParam = Field{Name: "name", Type: "string", Description: "Record name, as it appears on disk."}
var claimProjectParam = Field{Name: "project", Type: "string", Description: "Project name, matched case-insensitively."}
var claimMachineParam = Field{Name: "machine", Type: "string", Description: "Machine that holds the claim."}
var claimAgentParam = Field{Name: "agent", Type: "string", Description: "Agent that holds the claim."}
var spaceParam = Field{Name: "id", Type: "string", Description: "Space ID."}
var artifactIdParam = Field{Name: "id", Type: "string", Description: "Artifact identifier."}
var pendingParam = Field{Name: "id", Type: "string", Description: "Pending action id, from the confirm link."}
var memberParam = Field{Name: "email", Type: "string", Description: "Member's email address."}

var Registry = Response{Modules: []Module{
	{
		Name:        "auth",
		Description: "Sign-in, session identity, and the OIDC handshake with Porte.",
		Routes: []Route{
			{Method: "GET", Path: "/auth/config", Summary: "Report which sign-in methods this deployment offers", ResponseBody: ConfigResponse{}, Errors: []Error{internalError}},
			{Method: "POST", Path: "/auth/login", Summary: "Sign in with a password", Description: "Rate-limited. Absent entirely when SSO_ONLY is set.", RequestBody: LoginRequest{}, ResponseBody: AuthResponse{}, Errors: []Error{badRequest, unauthenticated, {Status: 429, Code: "resource_exhausted", Description: "Too many failed attempts."}, internalError}},
			{Method: "GET", Path: "/auth/oidc", Summary: "Begin the OIDC handshake", Description: "Redirects to Porte. A CLI adds flow=cli, an optional port=<1024-65535> and an optional cli_state nonce to be sent back to a loopback listener instead of the web app. A CLI flow with no usable port ends on a page showing a code to paste.", Errors: []Error{badRequest, internalError}},
			{Method: "GET", Path: "/auth/oidc/callback", Summary: "Complete the OIDC handshake", Description: "Redirects back with the session token in the URL fragment, never the query string. A CLI flow is redirected to http://127.0.0.1:<port>/ with a one-time code and the nonce it sent, or, with no port, answers 200 and a page showing that code. Every failure redirects to /login?error=sso, never JSON: this endpoint is only ever reached by a browser.", Errors: []Error{badRequest, unauthenticated, internalError}},
			{Method: "POST", Path: "/auth/oidc/exchange", Summary: "Exchange a CLI login code for a token", Description: "The code is single-use and valid for sixty seconds.", RequestBody: ExchangeRequest{}, ResponseBody: AuthResponse{}, Errors: []Error{badRequest, unauthenticated, internalError}},
			{Method: "GET", Path: "/auth/me", Summary: "Return the caller's identity", Auth: "bearer", ResponseBody: MeResponse{}, Errors: anyToken},
			{Method: "POST", Path: "/auth/logout", Summary: "Revoke the caller's session", Auth: "bearer", Status: http.StatusNoContent, Errors: anyToken},
		},
	},
	{
		Name:        "device",
		Description: "Device-code flow: how a CLI on a new machine gets a token without a browser redirect.",
		Routes: []Route{
			{Method: "POST", Path: "/auth/device/start", Summary: "Start a device authorization", RequestBody: DeviceStartRequest{}, ResponseBody: DeviceStartResponse{}, Errors: []Error{badRequest, internalError}},
			{Method: "POST", Path: "/auth/device/poll", Summary: "Poll for approval", Description: "Returns the token once a signed-in admin approves the request.", RequestBody: DevicePollRequest{}, ResponseBody: AuthResponse{}, Errors: []Error{badRequest, {Status: 428, Code: "failed_precondition", Description: "Still pending approval."}, internalError}},
			{Method: "GET", Path: "/auth/device/info", Summary: "Describe a pending device request", Auth: "bearer (admin)", QueryParams: []Field{{Name: "code", Type: "string", Description: "Device user code."}}, ResponseBody: DeviceInfoResponse{}, Errors: append(adminToken, notFound)},
			{Method: "POST", Path: "/auth/device/approve", Summary: "Approve a device request", Auth: "bearer (admin)", RequestBody: DeviceApproveRequest{}, ResponseBody: DeviceApproveResponse{}, Errors: append(adminToken, notFound)},
			{Method: "POST", Path: "/auth/device/deny", Summary: "Deny a device request", Auth: "bearer (admin)", RequestBody: DeviceDenyRequest{}, Status: http.StatusNoContent, Errors: append(adminToken, notFound)},
		},
	},
	{
		Name:        "status",
		Description: "What this Mycelium instance holds.",
		Routes: []Route{
			{Method: "GET", Path: "/status", Summary: "Summarize the instance", Description: "Counts of memory pages, rules and skills, plus the sync generation.", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: StatusResponse{}, Errors: anyToken},
		},
	},
	{
		Name:        "memory",
		Description: "Search and routing over the wiki, without transferring it.",
		Routes: []Route{
			{Method: "GET", Path: "/memory/search", Summary: "Search memory pages", Description: "Full-text search over the wiki; `q` is the query.", Auth: "bearer", QueryParams: []Field{{Name: "q", Type: "string", Description: "Search query."}, {Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: []SearchResult{}, Errors: append(anyToken, badRequest)},
			{Method: "POST", Path: "/memory/search", Summary: "Search memory pages, lexically and semantically", Description: "Fuses the BM25 ranking with a vector ranking using Reciprocal Rank Fusion, and returns ranked chunks rather than lines. `space_id` scopes the search to a space the caller belongs to. `degraded` is true when the embedding backend or its index was absent or unreachable and the answer is lexical-only — never an error.", Auth: "bearer", RequestBody: MemorySearchRequest{}, ResponseBody: MemorySearchResponse{}, Errors: append(anyToken, badRequest, forbidden)},
			{Method: "GET", Path: "/memory/index", Summary: "Return the index", Description: "The one-line-per-page router that decides what to read next.", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: IndexResponse{}, Errors: anyToken},
			{Method: "GET", Path: "/memory/index/status", Summary: "Report embedding index progress", Description: "Progress of the vector index behind semantic search, for a progress bar: `enabled` is false with every other field zeroed when no embedding backend is configured. `total_chunks` is what the wiki holds, `indexed_chunks` what the store holds, and `chunks_per_second`/`eta_seconds` are measured over the pass in flight and are zero when idle. `last_error` carries the most recent embedding failure, so a model outage is visible rather than a stalled bar. `space_id` scopes the counts to a space the caller belongs to.", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: MemoryIndexStatusResponse{}, Errors: append(anyToken, forbidden)},
		},
	},
	{
		Name:        "sessions",
		Description: "Agent session telemetry across every machine.",
		Routes: []Route{
			{Method: "GET", Path: "/sessions/stats", Summary: "Aggregate session statistics", Auth: "bearer", QueryParams: []Field{{Name: "since", Type: "string", Description: "Cutoff duration: 7d, 30d, 12h, or all."}, {Name: "by", Type: "string", Description: "Group by: project, machine, agent, branch, or model."}, {Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: StatsResponse{}, Errors: anyToken},
			{Method: "GET", Path: "/sessions/recent", Summary: "List recent sessions", Auth: "bearer", QueryParams: []Field{{Name: "limit", Type: "integer", Description: "Maximum number of recent sessions to return."}, {Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: []Session{}, Errors: anyToken},
			{Method: "GET", Path: "/sessions/live", Summary: "List sessions currently active", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: []LiveEntry{}, Errors: anyToken},
			{Method: "GET", Path: "/sessions/timeline", Summary: "Bucket session activity over time", Description: "Gap-filled buckets for charts. `since` (7d, 30d, 12h, all — default 30d), `bucket` (day or month), `by` (project, machine, agent, branch, model, total). Series are capped, with the remainder folded into `Other`.", Auth: "bearer", QueryParams: []Field{{Name: "since", Type: "string", Description: "Cutoff duration (default 30d)."}, {Name: "bucket", Type: "string", Description: "Bucket resolution: day or month."}, {Name: "by", Type: "string", Description: "Grouping: total, project, machine, agent, branch, or model."}, {Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: TimelineResponse{}, Errors: append(anyToken, badRequest)},
		},
	},
	{
		Name:        "claims",
		Description: "In-flight task leases, so a second agent can see or take over work already claimed on a repo.",
		Routes: []Route{
			{Method: "GET", Path: "/claims", Summary: "List every active claim", Description: "Liveness is resolved against the clock at read time, not stored.", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: []ClaimEntry{}, Errors: anyToken},
			{Method: "DELETE", Path: "/claims/{project}/{machine}/{agent}", Summary: "Release a claim", Description: "Releasing an already-absent claim is not an error.", Auth: "bearer", PathParams: []Field{claimProjectParam, claimMachineParam, claimAgentParam}, QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, Status: http.StatusNoContent, Errors: anyToken},
		},
	},
	{
		Name:        "usage",
		Description: "Subscription rate-limit windows, as reported by Claude Code on each machine.",
		Routes: []Route{
			{Method: "GET", Path: "/usage", Summary: "Read each machine's latest usage snapshot", Description: "Empty array when nothing has been recorded yet.", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: []UsageSnapshot{}, Errors: anyToken},
			{Method: "GET", Path: "/usage/history", Summary: "Read recorded usage samples over time", Description: "`since` defaults to 7d; `machine` filters to one machine. Samples are irregular, so labels are the sample instants and a missing window carries null.", Auth: "bearer", QueryParams: []Field{{Name: "since", Type: "string", Description: "Cutoff duration (default 7d)."}, {Name: "machine", Type: "string", Description: "Filter to one machine."}, {Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: UsageHistory{}, Errors: append(anyToken, badRequest)},
		},
	},
	{
		Name:        "settings",
		Description: "Instance-wide configuration.",
		Routes: []Route{
			{Method: "GET", Path: "/settings", Summary: "Read the settings", Auth: "bearer (admin)", ResponseBody: SettingsResponse{}, Errors: adminToken},
			{Method: "PUT", Path: "/settings", Summary: "Replace the settings", Auth: "bearer (admin)", RequestBody: Settings{}, ResponseBody: SettingsResponse{}, Errors: append(adminToken, badRequest)},
		},
	},
	{
		Name:        "rules",
		Description: "The agent rules Mycelium renders into each tool's native config.",
		Routes: []Route{
			{Method: "GET", Path: "/rules", Summary: "List rules", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: []string{}, Errors: anyToken},
			{Method: "GET", Path: "/rules/{name}", Summary: "Read one rule", Auth: "bearer", PathParams: []Field{nameParam}, QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: Rule{}, Errors: append(anyToken, notFound)},
			{Method: "PUT", Path: "/rules/{name}", Summary: "Create or replace a rule", Auth: "bearer", PathParams: []Field{nameParam}, QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, RequestBody: Rule{}, Status: http.StatusNoContent, Errors: append(anyToken, badRequest)},
			{Method: "DELETE", Path: "/rules/{name}", Summary: "Delete a rule", Auth: "bearer", PathParams: []Field{nameParam}, QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, Status: http.StatusNoContent, Errors: append(anyToken, notFound)},
		},
	},
	{
		Name:        "skills",
		Description: "The skills Mycelium installs into each agent.",
		Routes: []Route{
			{Method: "GET", Path: "/skills", Summary: "List skills", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: []string{}, Errors: anyToken},
			{Method: "GET", Path: "/skills/{name}", Summary: "Read one skill", Auth: "bearer", PathParams: []Field{nameParam}, QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: Skill{}, Errors: append(anyToken, notFound)},
			{Method: "PUT", Path: "/skills/{name}", Summary: "Create or replace a skill", Auth: "bearer", PathParams: []Field{nameParam}, QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, RequestBody: Skill{}, Status: http.StatusNoContent, Errors: append(anyToken, badRequest)},
			{Method: "DELETE", Path: "/skills/{name}", Summary: "Delete a skill", Auth: "bearer", PathParams: []Field{nameParam}, QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, Status: http.StatusNoContent, Errors: append(anyToken, notFound)},
		},
	},
	{
		Name:        "flows",
		Description: "Recorded shell procedures. Trust pins and run history are per-machine and never reach the server; this is the definition only.",
		Routes: []Route{
			{Method: "GET", Path: "/flows", Summary: "List flows", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: []string{}, Errors: anyToken},
			{Method: "GET", Path: "/flows/{name}", Summary: "Read one flow", Auth: "bearer", PathParams: []Field{nameParam}, QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: FlowDetail{}, Errors: append(anyToken, notFound)},
		},
	},
	{
		Name:        "models",
		Description: "Typed step extensions under extensions/models. The server has no bun and no shell, so a model is served as raw source, never executed. Reading one is a wildcard route (GET /models/{path}, where path may contain slashes) and is not listed here for the same reason /sync/files/* is not.",
		Routes: []Route{
			{Method: "GET", Path: "/models", Summary: "List models", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: []ModelInfo{}, Errors: anyToken},
		},
	},
	{
		Name:        "artifacts",
		Description: "Markdown artifacts synced across machines. HTML is not supported. Artifacts are permanent: nothing deletes one except DELETE.",
		Routes: []Route{
			{Method: "GET", Path: "/artifacts", Summary: "List artifacts", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: []ArtifactSummary{}, Errors: anyToken},
			{Method: "POST", Path: "/artifacts", Summary: "Create a markdown artifact", Description: "Markdown only: content that looks like HTML is refused, and the body is stored byte for byte. Content is capped at 1 MiB and the title at 200 characters; an empty title is taken from the first heading.", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, RequestBody: ArtifactWriteRequest{}, ResponseBody: ArtifactDetail{}, Status: http.StatusCreated, Errors: append(anyToken, badRequest, notMember, tooLarge)},
			{Method: "GET", Path: "/artifacts/{id}", Summary: "Read one artifact", Auth: "bearer", PathParams: []Field{artifactIdParam}, QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: ArtifactDetail{}, Errors: append(anyToken, notFound)},
			{Method: "PUT", Path: "/artifacts/{id}", Summary: "Rewrite a markdown artifact in place", Description: "Keeps the id and created time, replaces the body, renames when a title is given, and sets updated. Same limits as create.", Auth: "bearer", PathParams: []Field{artifactIdParam}, QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, RequestBody: ArtifactWriteRequest{}, ResponseBody: ArtifactDetail{}, Errors: append(anyToken, badRequest, notMember, tooLarge, notFound)},
			{Method: "DELETE", Path: "/artifacts/{id}", Summary: "Delete an artifact", Auth: "bearer", PathParams: []Field{artifactIdParam}, QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, Status: http.StatusNoContent, Errors: append(anyToken, notFound)},
		},
	},
	{
		Name:        "users",
		Description: "Who can sign in to this instance.",
		Routes: []Route{
			{Method: "GET", Path: "/users", Summary: "List users", Auth: "bearer", ResponseBody: []User{}, Errors: anyToken},
			{Method: "PUT", Path: "/users/{email}", Summary: "Update a user account", Auth: "bearer (admin)", PathParams: []Field{memberParam}, RequestBody: UpdateUserRequest{}, ResponseBody: User{}, Errors: append(adminToken, notFound)},
		},
	},
	{
		Name:        "spaces",
		Description: "Shared memory spaces and their membership.",
		Routes: []Route{
			{Method: "GET", Path: "/spaces", Summary: "List spaces the caller belongs to", Auth: "bearer", ResponseBody: SpacesListResponse{}, Errors: anyToken},
			{Method: "POST", Path: "/spaces", Summary: "Create a space", Auth: "bearer", RequestBody: CreateSpaceRequest{}, ResponseBody: SpaceResponse{}, Status: http.StatusCreated, Errors: append(anyToken, badRequest)},
			{Method: "PUT", Path: "/spaces/{id}", Summary: "Update a space", Description: "Refused for a space linked to Registre.", Auth: "bearer", PathParams: []Field{spaceParam}, RequestBody: UpdateSpaceRequest{}, ResponseBody: SpaceResponse{}, Errors: append(anyToken, badRequest, notFound, managed)},
			{Method: "DELETE", Path: "/spaces/{id}", Summary: "Delete a space", Auth: "bearer", PathParams: []Field{spaceParam}, Status: http.StatusNoContent, Errors: append(anyToken, forbidden, notFound, managed)},
			{Method: "GET", Path: "/spaces/{id}/members", Summary: "List a space's members", Auth: "bearer", PathParams: []Field{spaceParam}, ResponseBody: SpaceMembersResponse{}, Errors: append(anyToken, notFound)},
			{Method: "POST", Path: "/spaces/{id}/members", Summary: "Add a member", Auth: "bearer", PathParams: []Field{spaceParam}, RequestBody: AddMemberRequest{}, ResponseBody: MemberResponse{}, Status: http.StatusCreated, Errors: append(anyToken, badRequest, forbidden, notFound, managed)},
			{Method: "PUT", Path: "/spaces/{id}/members/{email}", Summary: "Change a member's role", Auth: "bearer", PathParams: []Field{spaceParam, memberParam}, RequestBody: UpdateMemberRequest{}, ResponseBody: MemberResponse{}, Errors: append(anyToken, badRequest, forbidden, notFound, managed)},
			{Method: "DELETE", Path: "/spaces/{id}/members/{email}", Summary: "Remove a member", Auth: "bearer", PathParams: []Field{spaceParam, memberParam}, Status: http.StatusNoContent, Errors: append(anyToken, forbidden, notFound, managed)},
			{Method: "POST", Path: "/spaces/{id}/leave", Summary: "Leave a space", Auth: "bearer", PathParams: []Field{spaceParam}, Status: http.StatusNoContent, Errors: append(anyToken, notFound, managed)},
		},
	},
	{
		Name:        "tokens",
		Description: "Per-machine sync tokens. Stored sha256-hashed, so a created token is shown once and never again.",
		Routes: []Route{
			{Method: "GET", Path: "/tokens", Summary: "List tokens", Description: "Metadata only — the token values are not recoverable.", Auth: "bearer (admin)", ResponseBody: []TokenInfo{}, Errors: adminToken},
			{Method: "POST", Path: "/tokens", Summary: "Issue a token", Description: "Returns the value once. Scope is admin or sync.", Auth: "bearer (admin)", RequestBody: CreateTokenRequest{}, ResponseBody: CreatedToken{}, Status: http.StatusOK, Errors: append(adminToken, badRequest)},
			{Method: "DELETE", Path: "/tokens/{name}", Summary: "Revoke a token", Auth: "bearer (admin)", PathParams: []Field{nameParam}, Status: http.StatusNoContent, Errors: append(adminToken, notFound)},
		},
	},
	{
		Name:        "sync",
		Description: "The file transfer the CLI runs. Paths are wiki-relative and may contain slashes, so these routes take a trailing wildcard rather than a named parameter.",
		Routes: []Route{
			{Method: "GET", Path: "/sync/tree", Summary: "List every file and its hash", Description: "The client diffs this against its own tree to decide what to move.", Auth: "bearer", QueryParams: []Field{{Name: "space_id", Type: "string", Description: "Space ID."}}, ResponseBody: []FileEntry{}, Errors: anyToken},
		},
	},
	{
		Name:        "registre",
		Description: "Registre owns shared spaces. This instance pulls its snapshot at boot, every five minutes, and shortly after a nudge.",
		Routes: []Route{
			{Method: "POST", Path: "/registre/sync", Summary: "Nudge a pull of Registre's space snapshot", Description: "Unauthenticated and body-less: it only schedules the authenticated pull sooner.", Status: http.StatusAccepted},
			{Method: "POST", Path: "/registre/deletions", Summary: "Receive Registre's signed user.deleted notice", Description: "Unauthenticated by session: the body is a signed security event token that porte verifies against the issuer's keys. Erasing a person is idempotent.", Status: http.StatusAccepted, Errors: []Error{badRequest}},
		},
	},
	{
		Name:        "mcp",
		Description: "The confirm page of the remote MCP server. An agent's request to delete a wiki page waits here until the person confirms or dismisses it in their own browser session; CLI, API and delegated tokens are refused.",
		Routes: []Route{
			{Method: "GET", Path: "/mcp/pending", Summary: "List the actions waiting for you", Description: "Newest first, each with the page it would delete and whether that page still matches what you were asked about.", Auth: "bearer (browser session)", ResponseBody: PendingList{}, Errors: append(anyToken, forbidden, notFound)},
			{Method: "POST", Path: "/mcp/pending/{id}/confirm", Summary: "Confirm a waiting action", Description: "Runs it once. A page edited since the question closes the action instead.", Auth: "bearer (browser session)", PathParams: []Field{pendingParam}, RequestBody: EmptyRequest{}, ResponseBody: PendingResult{}, Errors: append(anyToken, forbidden, notFound)},
			{Method: "POST", Path: "/mcp/pending/{id}/dismiss", Summary: "Dismiss a waiting action", Description: "Closes it without running it.", Auth: "bearer (browser session)", PathParams: []Field{pendingParam}, RequestBody: EmptyRequest{}, ResponseBody: PendingResult{}, Errors: append(anyToken, forbidden, notFound)},
		},
	},
}}
