package utils

import "time"

// TimeFormat format time
func TimeFormat(t time.Time) string {
	timeString := t.Format("2006-01-02 15:04:05")
	return timeString
}
