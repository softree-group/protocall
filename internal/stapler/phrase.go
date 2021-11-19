package stapler

import "time"

type Phrase struct {
	Time time.Time
	User string
	Text string
}
