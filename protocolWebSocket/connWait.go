package protocolWebSocket

import (
	"time"
)

type connWait struct {
	min    time.Duration
	max    time.Duration
	next   time.Duration
	factor uint16
}

func newConnWait() connWait {
	return connWait{
		min:    1 * time.Second,
		max:    5 * time.Minute,
		factor: 5,
	}
}

func (c *connWait) nextDuration() time.Duration {
	switch {
	case c.next == 0:
		c.next = c.min
		break
	case c.next >= c.max:
		c.next = c.max
		break
	default:
		c.next = c.next * time.Duration(c.factor)
		if c.next >= c.max {
			c.next = c.max
		}
	}
	return c.next
}

func (c *connWait) reset() {
	c.next = 0
}
