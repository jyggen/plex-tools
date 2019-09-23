package plex

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/jrudio/go-plex-client"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Media struct {
	AudioChannels int
	AudioCodec    string
	Bitrate       uint64
	Duration      time.Duration
	FrameRate     string
	Quality       string
	Rating        float64
	Size          uint64
	Title         string
	VideoCodec    string
	Year          int
}

var specialCharacters = regexp.MustCompile(`(\s|\.|,|_|-|=|'|\|)+`)
var nonWordCharacters = regexp.MustCompile(`[^\w\s]`)
var conjunctions = regexp.MustCompile(`(?i)\b(a|an|the|and|or|of)\b\s?`)
var superfluousWhitespace = regexp.MustCompile(`\s{2,}`)
var moreCharacters = regexp.MustCompile(`(&|:|\|/)+`)

func NewMediaSlice(v plex.Metadata) []*Media {
	media := make([]*Media, len(v.Media))

	for k, m := range v.Media {
		quality := strings.ToUpper(m.VideoResolution)

		if _, err := strconv.Atoi(m.VideoResolution); err == nil {
			quality += "p"
		}

		for _, p := range m.Part {
			title := v.Title

			if v.Type == "episode" {
				title = fmt.Sprintf("%s (S%02dE%02d): %s", v.GrandparentTitle, v.ParentIndex, v.Index, v.Title)
			}

			media[k] = &Media{
				AudioChannels: m.AudioChannels,
				AudioCodec:    m.AudioCodec,
				Bitrate:       uint64(m.Bitrate) * humanize.KByte,
				Duration:      time.Duration(p.Duration) * time.Millisecond,
				FrameRate:     m.VideoFrameRate,
				Quality:       quality,
				Rating:        v.Rating,
				Size:          uint64(p.Size),
				Title:         title,
				VideoCodec:    m.VideoCodec,
				Year:          v.Year,
			}
		}
	}

	return media
}

func (m *Media) DurationInNanoseconds() int64 {
	return m.Duration.Nanoseconds()
}

func (m *Media) HumanizeBitRate() string {
	return humanize.Bytes(m.Bitrate)
}

func (m *Media) HumanizeDuration() string {
	return humanizeDuration(m.Duration)
}

func (m *Media) HumanizeSize() string {
	return humanize.Bytes(m.Size)
}

func (m *Media) SortTitle() string {
	title := specialCharacters.ReplaceAllString(m.Title, " ")
	title = nonWordCharacters.ReplaceAllString(title, "")
	title = conjunctions.ReplaceAllString(title, "")
	title = superfluousWhitespace.ReplaceAllString(title, " ")
	title = moreCharacters.ReplaceAllString(title, "")

	return title
}
