package guard

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Signed time (AUTH-DESIGN §4.3): a Thing never decides with NTP's time. Its
// clock is the last heartbeat auth signed for its nonce, plus the monotonic
// time elapsed since. The first heartbeat after a start must be later than
// the floor, written to a file at most hourly; the next ones later than the
// last and within maxDrift of the clock.

// maxDrift is the furthest a heartbeat may be from the clock; floorEvery,
// how often the floor is written.
const maxDrift, floorEvery = time.Hour, time.Hour

type clock struct {
	mu    sync.Mutex
	file  string    // "": no floor kept
	floor int64     // unix seconds
	base  time.Time // time of the last heartbeat accepted
	at    time.Time // local (monotonic) reading when it came; zero: none yet
	saved time.Time // local reading of the last floor written
}

func newClock(file string) *clock {
	c := &clock{file: file}
	if b, err := os.ReadFile(file); err == nil {
		c.floor, _ = strconv.ParseInt(strings.TrimSpace(string(b)), 10, 64)
	}
	return c
}

// now is the signed time; false before the first heartbeat.
func (c *clock) now() (time.Time, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.at.IsZero() {
		return time.Time{}, false
	}
	return c.base.Add(time.Since(c.at)), true
}

// accept takes the heartbeat issued at ms, its nonce already checked.
func (c *clock) accept(ms int64) error {
	ts := time.UnixMilli(ms)
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.at.IsZero() {
		if ts.Unix() <= c.floor {
			return errors.New("heartbeat before the time floor")
		}
	} else if now := c.base.Add(time.Since(c.at)); !ts.After(c.base) || ts.Sub(now).Abs() > maxDrift {
		return errors.New("heartbeat not later than the last, or too far from the clock")
	}
	c.base, c.at = ts, time.Now()
	if c.file != "" && (c.saved.IsZero() || time.Since(c.saved) >= floorEvery) {
		tmp := c.file + ".tmp"
		if os.WriteFile(tmp, []byte(strconv.FormatInt(ts.Unix(), 10)), 0o600) == nil && os.Rename(tmp, c.file) == nil {
			c.floor, c.saved = ts.Unix(), time.Now()
		}
	}
	return nil
}
