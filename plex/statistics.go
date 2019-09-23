package plex

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"io"
	"math"
	"sort"
	"time"
)

type Statistics struct {
	probe         *Probe
	audioChannels map[int]int
	audioCodec    map[string]int
	bitrate       struct {
		minimum uint64
		mean    uint64
		median  uint64
		maximum uint64
		total   uint64
	}
	duration struct {
		minimum time.Duration
		mean    time.Duration
		median  time.Duration
		maximum time.Duration
		total   time.Duration
	}
	quality map[string]int
	rating  struct {
		minimum float64
		mean    float64
		median  float64
		maximum float64
		total   float64
	}
	size struct {
		minimum uint64
		mean    uint64
		median  uint64
		maximum uint64
		total   uint64
	}
	total      int
	videoCodec map[string]int
	year       struct {
		minimum int
		mean    int
		median  int
		maximum int
		total   int
	}
}

func (p *Probe) Statistics() *Statistics {
	s := &Statistics{}

	s.audioChannels = make(map[int]int, 0)
	s.audioCodec = make(map[string]int, 0)
	s.bitrate.minimum = math.MaxUint64
	s.duration.minimum = 24 * 365 * time.Hour
	s.quality = make(map[string]int, 0)
	s.rating.minimum = math.MaxFloat64
	s.size.minimum = math.MaxUint64
	s.total = len(p.media)
	s.videoCodec = make(map[string]int, 0)
	s.year.minimum = math.MaxUint16

	bitrates := make([]uint64, s.total)
	durations := make([]time.Duration, s.total)
	ratings := make([]float64, s.total)
	sizes := make([]uint64, s.total)
	years := make([]int, s.total)

	for k, v := range p.media {
		if _, ok := s.audioChannels[v.AudioChannels]; !ok {
			s.audioChannels[v.AudioChannels] = 0
		}

		s.audioChannels[v.AudioChannels] += 1

		if _, ok := s.audioCodec[v.AudioCodec]; !ok {
			s.audioCodec[v.AudioCodec] = 0
		}

		s.audioCodec[v.AudioCodec] += 1

		if _, ok := s.quality[v.Quality]; !ok {
			s.quality[v.Quality] = 0
		}

		s.quality[v.Quality] += 1

		if _, ok := s.videoCodec[v.VideoCodec]; !ok {
			s.videoCodec[v.VideoCodec] = 0
		}

		s.videoCodec[v.VideoCodec] += 1

		bitrates[k] = v.Bitrate
		durations[k] = v.Duration
		ratings[k] = v.Rating
		sizes[k] = v.Size
		years[k] = v.Year

		s.bitrate.total += v.Bitrate
		s.duration.total += v.Duration
		s.rating.total += v.Rating
		s.size.total += v.Size
		s.year.total += v.Year

		if v.Bitrate > s.bitrate.maximum {
			s.bitrate.maximum = v.Bitrate
		}

		if v.Bitrate < s.bitrate.minimum {
			s.bitrate.minimum = v.Bitrate
		}

		if v.Duration > s.duration.maximum {
			s.duration.maximum = v.Duration
		}

		if v.Duration < s.duration.minimum {
			s.duration.minimum = v.Duration
		}

		if v.Rating > s.rating.maximum {
			s.rating.maximum = v.Rating
		}

		if v.Rating < s.rating.minimum {
			s.rating.minimum = v.Rating
		}

		if v.Size > s.size.maximum {
			s.size.maximum = v.Size
		}

		if v.Size < s.size.minimum {
			s.size.minimum = v.Size
		}

		if v.Year > s.year.maximum {
			s.year.maximum = v.Year
		}

		if v.Year < s.year.minimum {
			s.year.minimum = v.Year
		}
	}

	s.bitrate.mean = uint64(math.Round(float64(s.bitrate.total) / float64(s.total)))
	s.duration.mean = time.Duration(math.Round(float64(s.duration.total) / float64(s.total)))
	s.rating.mean = s.rating.total / float64(s.total)
	s.size.mean = uint64(math.Round(float64(s.size.total) / float64(s.total)))
	s.year.mean = int(math.Round(float64(s.year.total) / float64(s.total)))

	sort.Slice(bitrates, func(i, j int) bool {
		return bitrates[i] > bitrates[j]
	})

	sort.Slice(durations, func(i, j int) bool {
		return durations[i] > durations[j]
	})

	sort.Slice(ratings, func(i, j int) bool {
		return ratings[i] > ratings[j]
	})

	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	sort.Slice(years, func(i, j int) bool {
		return years[i] > years[j]
	})

	h := s.total / 2

	if s.total%2 == 0 {
		s.bitrate.median = bitrates[h]
		s.duration.median = durations[h]
		s.rating.median = ratings[h]
		s.size.median = sizes[h]
		s.year.median = years[h]
	} else {
		s.bitrate.median = (bitrates[h-1] + bitrates[h]) / 2
		s.duration.median = (durations[h-1] + durations[h]) / 2
		s.rating.median = (ratings[h-1] + ratings[h]) / 2
		s.size.median = (sizes[h-1] + sizes[h]) / 2
		s.year.median = (years[h-1] + years[h]) / 2
	}

	return s
}

func (s *Statistics) Ascii(w io.Writer) {
	t := tablewriter.NewWriter(w)

	t.SetAlignment(tablewriter.ALIGN_CENTER)
	t.SetHeader([]string{"Type", "Min", "Mean", "Median", "Max", "Total"})

	t.Append([]string{
		"Bitrate",
		humanize.Bytes(s.bitrate.minimum) + "ps",
		humanize.Bytes(s.bitrate.mean) + "ps",
		humanize.Bytes(s.bitrate.median) + "ps",
		humanize.Bytes(s.bitrate.maximum) + "ps",
		"n/a",
	})

	t.Append([]string{
		"Rating",
		fmt.Sprintf("%.2f", s.rating.minimum),
		fmt.Sprintf("%.2f", s.rating.mean),
		fmt.Sprintf("%.2f", s.rating.median),
		fmt.Sprintf("%.2f", s.rating.maximum),
		"n/a",
	})

	t.Append([]string{
		"Duration",
		humanizeDuration(s.duration.minimum),
		humanizeDuration(s.duration.mean),
		humanizeDuration(s.duration.median),
		humanizeDuration(s.duration.maximum),
		humanizeDuration(s.duration.total),
	})

	t.Append([]string{
		"Size",
		humanize.Bytes(s.size.minimum),
		humanize.Bytes(s.size.mean),
		humanize.Bytes(s.size.median),
		humanize.Bytes(s.size.maximum),
		humanize.Bytes(s.size.total),
	})

	t.Append([]string{
		"Year",
		fmt.Sprintf("%d", s.year.minimum),
		fmt.Sprintf("%d", s.year.mean),
		fmt.Sprintf("%d", s.year.median),
		fmt.Sprintf("%d", s.year.maximum),
		"n/a",
	})

	t.Render()

	t = tablewriter.NewWriter(w)

	t.SetAlignment(tablewriter.ALIGN_CENTER)
	t.SetHeader([]string{"Audio Channels", "Count", "Percentage"})

	for k, v := range s.audioChannels {
		t.Append([]string{
			fmt.Sprintf("%d", k),
			fmt.Sprintf("%d", v),
			fmt.Sprintf("%.2f%%", float64(v)/float64(s.total)*100),
		})
	}

	t.Render()

	t = tablewriter.NewWriter(w)

	t.SetAlignment(tablewriter.ALIGN_CENTER)
	t.SetHeader([]string{"Audio Codec", "Count", "Percentage"})

	for k, v := range s.audioCodec {
		t.Append([]string{
			k,
			fmt.Sprintf("%d", v),
			fmt.Sprintf("%.2f%%", float64(v)/float64(s.total)*100),
		})
	}

	t.Render()

	t = tablewriter.NewWriter(w)

	t.SetAlignment(tablewriter.ALIGN_CENTER)
	t.SetHeader([]string{"Quality", "Count", "Percentage"})

	for k, v := range s.quality {
		t.Append([]string{
			k,
			fmt.Sprintf("%d", v),
			fmt.Sprintf("%.2f%%", float64(v)/float64(s.total)*100),
		})
	}

	t.Render()

	t = tablewriter.NewWriter(w)

	t.SetAlignment(tablewriter.ALIGN_CENTER)
	t.SetHeader([]string{"Video Codec", "Count", "Percentage"})

	for k, v := range s.videoCodec {
		t.Append([]string{
			k,
			fmt.Sprintf("%d", v),
			fmt.Sprintf("%.2f%%", float64(v)/float64(s.total)*100),
		})
	}

	t.Render()
}
