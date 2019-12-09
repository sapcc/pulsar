package util

import (
	"fmt"
	"time"
)

const timestampFormat = "15:04:05 01.02.2006 UTC"

// HumanizeTimestamp may be used to increase readability of a timestamp.
func HumanizeTimestamp(t time.Time) string {
	return t.Format(timestampFormat)
}

// HumanizeDuration may be used to increase readability of a timestamp.
func HumanizeDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

// StringToTimestamp converts a string to RFC 3339 timestamp.
func StringToTimestamp(theString string) time.Time {
	t, _ := time.Parse(time.RFC3339, theString)
	return t
}

// TimestampToString converts the given time to RFC 3339 format string.
func TimestampToString(ts time.Time) string {
	return ts.Format(time.RFC3339)
}

func TimeStartOfDay() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 00, 00, 00, 00, now.Location())
}

func TimeEndOfDay() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 00, 00, now.Location())
}