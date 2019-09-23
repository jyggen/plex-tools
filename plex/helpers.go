package plex

import (
	"fmt"
	"time"
)

func humanizeDuration(dur time.Duration) string {
	dur = dur.Round(time.Second)
	hour := dur / time.Hour
	dur -= hour * time.Hour
	min := dur / time.Minute
	dur -= min * time.Minute
	sec := dur / time.Second

	return fmt.Sprintf("%02dh %02dm %02ds", hour, min, sec)
}
