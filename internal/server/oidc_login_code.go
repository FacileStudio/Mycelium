package server

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"
	"time"
)

// loginCode is a bearer credential for sixty seconds, so it is kept hashed and
// carries the identity rather than an already-minted session: a code nobody
// exchanges must not leave a live token behind.
type loginCode struct {
	Email   string    `json:"email"`
	Scope   string    `json:"scope"`
	Expires time.Time `json:"expires"`
}

// loginCodeStore tracks one-time codes handed to a CLI login. It persists to
// path keyed by the code's own hash — the store never holds a plaintext code —
// so a restarted server (any push to main is a deploy) keeps a code a person is
// mid-paste instead of telling them it vanished. A nil or empty path keeps the
// store purely in memory, which is what the unit tests want.
type loginCodeStore struct {
	mu    sync.Mutex
	codes map[string]loginCode
	path  string
	log   *slog.Logger
}

func newLoginCodeStore(path string, log *slog.Logger) *loginCodeStore {
	l := &loginCodeStore{
		codes: make(map[string]loginCode),
		path:  path,
		log:   log,
	}
	l.load()
	return l
}

func (l *loginCodeStore) load() {
	if l.path == "" {
		return
	}
	data, err := os.ReadFile(l.path)
	if err != nil {
		return
	}
	var raw map[string]loginCode
	if err := json.Unmarshal(data, &raw); err != nil {
		if l.log != nil {
			l.log.Error("login-codes: failed to parse", slog.String("path", l.path), slog.Any("error", err))
		}
		return
	}
	now := time.Now()
	for hash, code := range raw {
		if now.After(code.Expires) {
			continue
		}
		l.codes[hash] = code
	}
}

func (l *loginCodeStore) persist() {
	if l.path == "" {
		return
	}
	data, err := json.MarshalIndent(l.codes, "", "  ")
	if err != nil {
		if l.log != nil {
			l.log.Error("login-codes: marshal failed", slog.Any("error", err))
		}
		return
	}
	if err := atomicWriteFile(l.path, data); err != nil {
		if l.log != nil {
			l.log.Error("login-codes: write failed", slog.String("path", l.path), slog.Any("error", err))
		}
	}
}

func (l *loginCodeStore) sweep(now time.Time) {
	for hash, code := range l.codes {
		if now.After(code.Expires) {
			delete(l.codes, hash)
		}
	}
}

func (l *loginCodeStore) create(hash, email, scope string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sweep(now)
	if len(l.codes) >= maxPendingLoginCodes {
		return false
	}
	l.codes[hash] = loginCode{Email: email, Scope: scope, Expires: now.Add(loginCodeTTL)}
	l.persist()
	return true
}

// consume returns the identity behind a code and removes it in the same
// critical section, which is what makes the code single-use under concurrency.
// The second return distinguishes an expired code from one that never existed
// so the caller can log a replay as the incident it is.
func (l *loginCodeStore) consume(hash string, now time.Time) (loginCode, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	code, ok := l.codes[hash]
	delete(l.codes, hash)
	l.sweep(now)
	l.persist()
	if !ok || now.After(code.Expires) {
		return loginCode{}, false
	}
	return code, true
}
