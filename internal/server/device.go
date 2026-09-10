package server

import (
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"
)

const (
	deviceCodeTTL      = 10 * time.Minute
	devicePollInterval = 5
	maxPendingDevices  = 256
	// No vowels (avoids real words) and no 0/1/I/L/O/U (avoids confusion).
	userCodeAlphabet = "23456789BCDFGHJKMNPQRSTVWXYZ"
)

const (
	devicePending  deviceStatus = "pending"
	deviceApproved deviceStatus = "approved"
	deviceDenied   deviceStatus = "denied"
)

// ErrTooManyDevices is returned when too many device authorizations are pending.
var ErrTooManyDevices = errors.New("too many pending device authorizations")

type deviceStatus string

type deviceRequest struct {
	// DeviceCode is the bearer the terminal polls with, handed back exactly once
	// at create and never serialized.
	DeviceCode string       `json:"-"`
	UserCode   string       `json:"user_code"`
	Machine    string       `json:"machine"`
	IP         string       `json:"ip"`
	Status     deviceStatus `json:"status"`
	// Token is the minted credential delivered to a poll. It is the one
	// long-lived secret in the record, so it never reaches disk, and an approved
	// request is dropped from the persisted set entirely so a restart cannot
	// replay it.
	Token   string    `json:"-"`
	Expires time.Time `json:"expires"`
}

// deviceStore tracks pending and in-flight device authorizations. Entries
// persist to path (keyed by hash of the device code, matching tokens.json) so a
// restarted server — any push to main is a deploy — keeps a login a human
// started before the restart instead of telling them their code died. A nil or
// empty path keeps the store purely in memory, which is what the unit tests
// want.
type deviceStore struct {
	mu       sync.Mutex
	byDevice map[string]*deviceRequest
	byUser   map[string]*deviceRequest
	path     string
	log      *slog.Logger
}

func newDeviceStore(path string, log *slog.Logger) *deviceStore {
	d := &deviceStore{
		byDevice: make(map[string]*deviceRequest),
		byUser:   make(map[string]*deviceRequest),
		path:     path,
		log:      log,
	}
	d.load()
	return d
}

// normalizeUserCode upper-cases and strips anything that isn't part of the
// code (hyphens, spaces) so user entry is forgiving and matches the stored key.
func normalizeUserCode(code string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(code) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (d *deviceStore) sweep(now time.Time) {
	for code, req := range d.byDevice {
		if now.After(req.Expires) {
			delete(d.byDevice, code)
			delete(d.byUser, normalizeUserCode(req.UserCode))
		}
	}
}

func (d *deviceStore) create(machine, ip string, now time.Time) (deviceRequest, error) {
	deviceCode, err := generateToken()
	if err != nil {
		return deviceRequest{}, err
	}
	userCode, err := generateUserCode()
	if err != nil {
		return deviceRequest{}, err
	}
	req := &deviceRequest{
		DeviceCode: deviceCode,
		UserCode:   userCode,
		Machine:    machine,
		IP:         ip,
		Status:     devicePending,
		Expires:    now.Add(deviceCodeTTL),
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.sweep(now)
	if len(d.byDevice) >= maxPendingDevices {
		return deviceRequest{}, ErrTooManyDevices
	}
	d.byDevice[deviceKey(deviceCode)] = req
	d.byUser[normalizeUserCode(userCode)] = req
	d.persist()
	return *req, nil
}

func (d *deviceStore) info(userCode string) (deviceRequest, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.sweep(time.Now())
	req, ok := d.byUser[normalizeUserCode(userCode)]
	if !ok {
		return deviceRequest{}, false
	}
	return *req, true
}

func (d *deviceStore) approve(userCode, token string) (string, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.sweep(time.Now())
	req, ok := d.byUser[normalizeUserCode(userCode)]
	if !ok || req.Status != devicePending {
		return "", false
	}
	req.Status = deviceApproved
	req.Token = token
	d.persist()
	return req.Machine, true
}

func (d *deviceStore) deny(userCode string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.sweep(time.Now())
	req, ok := d.byUser[normalizeUserCode(userCode)]
	if !ok || req.Status != devicePending {
		return false
	}
	req.Status = deviceDenied
	d.persist()
	return true
}

// poll returns the request status; once approved it returns the token and
// consumes the request so a token can only be retrieved once.
func (d *deviceStore) poll(deviceCode string) (deviceStatus, string, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.sweep(time.Now())
	key := deviceKey(deviceCode)
	req, ok := d.byDevice[key]
	if !ok {
		return "", "", false
	}
	if req.Status == deviceApproved {
		token := req.Token
		delete(d.byDevice, key)
		delete(d.byUser, normalizeUserCode(req.UserCode))
		d.persist()
		return deviceApproved, token, true
	}
	return req.Status, "", true
}
