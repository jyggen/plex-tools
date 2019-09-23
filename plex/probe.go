package plex

import (
	"fmt"
	"github.com/jrudio/go-plex-client"
	"github.com/olekukonko/tablewriter"
	"html/template"
	"io"
	"strconv"
)

const probeTemplate = `<!doctype html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
		<title>Probe: {{.Probe.Library}} @ {{.Probe.Server.Name}}</title>
	</head>
	<body>
		<table id="probe" class="table table-striped table-sm">
			<thead class="thead-dark">
				<tr>
					<th scope="col">Title</th>
					<th scope="col">Year</th>
					<th scope="col">Duration</th>
					<th scope="col">Rating</th>
					<th scope="col">Size</th>
					<th scope="col">Quality</th>
					<th scope="col">Bitrate</th>
					<th scope="col">Video</th>
					<th scope="col">Frame Rate</th>
					<th scope="col">Audio</th>
					<th scope="col">Channels</th>
				</tr>
			</thead>
			<tbody>
				{{range .Probe.Media}}<tr>
					<td data-sort="{{.SortTitle}}">{{.Title}}</td>
					<td>{{.Year}}</td>
					<td data-sort="{{.DurationInNanoseconds}}">{{.HumanizeDuration}}</td>
					<td>{{.Rating}}</td>
					<td data-sort="{{.Size}}">{{.HumanizeSize}}</td>
					<td>{{.Quality}}</td>
					<td data-sort="{{.Bitrate}}">{{.HumanizeBitRate}}</td>
					<td>{{.VideoCodec}}</td>
					<td>{{.FrameRate}}</td>
					<td>{{.AudioCodec}}</td>
					<td>{{.AudioChannels}}</td>
				</tr>{{end}}
			</tbody>
		</table>
		<script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
		<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/tablesort/5.1.0/tablesort.min.js" integrity="sha256-p3wukcf2d2jxbVnlqPDO9t4AAjnl42D2aIzrK4S0X6w=" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/tablesort/5.1.0/sorts/tablesort.number.min.js" integrity="sha256-ra1pWQ7MfuVIolZ/phcEXegs9m1ehXaCNI8cmc3gJEs=" crossorigin="anonymous"></script>
		<script>
		new Tablesort(document.getElementById('probe'));
		</script>
	</body>
</html>`

type Probe struct {
	library string
	media   []*Media
	server  *Server
}

func (p *Plex) Probe(libraryKey string) (*Probe, error) {
	lc, err := p.client.GetLibraryContent(libraryKey, "")

	if err != nil {
		return nil, err
	}

	media, err := p.probe(lc.MediaContainer.Metadata)

	if err != nil {
		return nil, err
	}

	return &Probe{
		library: lc.MediaContainer.LibrarySectionTitle,
		media:   media,
		server:  p.server,
	}, nil
}

func (p *Plex) probe(metadata []plex.Metadata) ([]*Media, error) {
	media := make([]*Media, 0)

	for _, m := range metadata {
		switch m.Type {
		case "episode", "movie":
			media = append(media, NewMediaSlice(m)...)
		case "season":
			sub, err := p.client.GetEpisodes(m.RatingKey)

			if err != nil {
				return media, err
			}

			res, err := p.probe(sub.MediaContainer.Metadata)

			if err != nil {
				return media, err
			}

			media = append(media, res...)
		case "show":
			sub, err := p.client.GetMetadataChildren(m.RatingKey)

			if err != nil {
				return media, err
			}

			res, err := p.probe(sub.MediaContainer.Metadata)

			if err != nil {
				return media, err
			}

			media = append(media, res...)
		default:
			return media, fmt.Errorf("unsupported type \"%s\"", m.Type)
		}
	}

	return media, nil
}

func (p *Probe) Ascii(w io.Writer) {
	t := tablewriter.NewWriter(w)

	t.SetHeader([]string{"Title", "Year", "Size", "Quality", "Bit Rate", "Video", "Frame Rate", "Audio", "Channels"})

	for _, m := range p.media {
		t.Append([]string{
			m.Title,
			strconv.Itoa(m.Year),
			m.HumanizeSize(),
			m.Quality,
			m.HumanizeBitRate(),
			m.VideoCodec,
			m.FrameRate,
			m.AudioCodec,
			strconv.Itoa(m.AudioChannels),
		})
	}

	t.Render()
}

func (p *Probe) Html(w io.Writer) error {
	tmpl, err := template.New("probe").Parse(probeTemplate)

	if err != nil {
		return err
	}

	return tmpl.Execute(w, struct {
		Probe *Probe
	}{
		Probe: p,
	})
}

func (p *Probe) Library() string {
	return p.library
}

func (p *Probe) Media() []*Media {
	return p.media
}

func (p *Probe) Server() *Server {
	return p.server
}
