package domain

import (
	"strings"
	"time"
)

const UTCFormat = "2006-01-02T15:04:05.00Z"

// Timestamp is a time.Time wrapper that has consistent JSON marshalling and unmarshalling format.
// Format is as defined by format.formatTime and format.parseTime.
type Timestamp struct {
	time.Time
}

// MarshalJSON converts t.Time to UTC using format.formatTime.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	s := `"` + formatTime(t.Time) + `"`
	return []byte(s), nil
}

// Unmarshal converts time string in UTC to Timestamp using format.parseTime.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)

	ts, err := parseTime(s)
	if err != nil {
		return err
	}
	t.Time = ts
	return nil
}

// FormatTimestamp formats time with format "2006-01-02T15:04:05.000Z".
// Note the timezone is always UTC with no offsets ('Z').
func formatTime(time time.Time) string {
	return time.UTC().Format(UTCFormat)
}

// parseTime parses str with format time.RFC3339Nano.
func parseTime(str string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, str)
}
