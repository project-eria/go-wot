package guard

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// Signed time (AUTH-DESIGN §4.3): a Thing never decides with NTP's time. Its
// clock is the last heartbeat auth signed for its nonce, plus the monotonic
// time elapsed since; after a start and before any heartbeat, the floor kept
// in a file plus the monotonic time. The first heartbeat after a start must
// be later than the floor; the next ones later than the last and within
// maxDrift of the clock. The file also keeps the time of the last heartbeat:
// isolated since = clock − last heartbeat (§4.4).

// maxDrift is the furthest a heartbeat may be from the clock; floorEvery,
// how often the file is written.
const maxDrift, floorEvery = time.Hour, time.Hour

type clock struct {
	mu       sync.Mutex
	file     string    // "": nothing kept
	floor    int64     // unix seconds
	lastBeat int64     // unix seconds of the last heartbeat accepted, ever
	base     time.Time // time at `at`: the last heartbeat, or the floor
	at       time.Time // local (monotonic) reading; zero: no time at all
	beatAt   time.Time // local reading of the last heartbeat; zero: none since the start
	saved    time.Time // local reading of the last write
}

// newClock reads the file: "<floor> <last heartbeat>".
func newClock(file string) *clock {
	c := &clock{file: file}
	if b, err := os.ReadFile(file); err == nil {
		fmt.Sscan(string(b), &c.floor, &c.lastBeat)
	}
	if c.floor > 0 {
		c.base, c.at = time.Unix(c.floor, 0), time.Now()
	}
	return c
}

// now is the clock; false without floor nor heartbeat.
func (c *clock) now() (time.Time, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.read()
}

func (c *clock) read() (time.Time, bool) {
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
	if c.beatAt.IsZero() {
		if ts.Unix() <= c.floor {
			return errors.New("heartbeat before the time floor")
		}
	} else if now, _ := c.read(); !ts.After(c.base) || ts.Sub(now).Abs() > maxDrift {
		return errors.New("heartbeat not later than the last, or too far from the clock")
	}
	c.base, c.at, c.beatAt, c.lastBeat = ts, time.Now(), time.Now(), ts.Unix()
	c.write(false)
	return nil
}

// missed tells whether no heartbeat came for `after` (or since the start).
func (c *clock) missed(after time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.beatAt.IsZero() || time.Since(c.beatAt) > after
}

// isolated is how long since the last heartbeat accepted, ever: 0 when none
// ever came (never activated) or without a clock.
func (c *clock) isolated() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	now, ok := c.read()
	if !ok || c.lastBeat == 0 {
		return 0
	}
	return now.Sub(time.Unix(c.lastBeat, 0))
}

// save writes the floor (the clock now) and the last heartbeat: at most
// every floorEvery, or at once when forced (a switch of mode).
func (c *clock) save(force bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.write(force)
}

func (c *clock) write(force bool) {
	now, ok := c.read()
	if c.file == "" || !ok || !force && !c.saved.IsZero() && time.Since(c.saved) < floorEvery {
		return
	}
	tmp := c.file + ".tmp"
	if os.WriteFile(tmp, fmt.Appendf(nil, "%d %d", now.Unix(), c.lastBeat), 0o600) == nil && os.Rename(tmp, c.file) == nil {
		c.floor, c.saved = now.Unix(), time.Now()
	}
}
