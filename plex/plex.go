package plex

import (
	"errors"
	"fmt"
	"github.com/jrudio/go-plex-client"
)

type Plex struct {
	client *plex.Plex
	server *Server
	token  string
}

type Server = plex.PMSDevices

func New(token string) (*Plex, error) {
	if token == "" {
		return nil, errors.New("no Plex access token configured")
	}

	c, err := plex.New("", token)

	if err != nil {
		return nil, err
	}

	return &Plex{
		client: c,
		token:  token,
	}, nil
}

func (p *Plex) GetLibraryKeyByTitle(title string) (string, error) {
	libraries, err := p.client.GetLibraries()

	if err != nil {
		return "", err
	}

	for _, l := range libraries.MediaContainer.Directory {
		if l.Title == title {
			return l.Key, nil
		}
	}

	return "", fmt.Errorf("no library titled \"%s\" found", title)
}

func (p *Plex) GetServerByName(name string) (*Server, error) {
	servers, err := p.client.GetServers()

	if err != nil {
		return nil, err
	}

	for _, s := range servers {
		if s.Name == name {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("no server named \"%s\" found", name)
}

func (p *Plex) UseServer(server *Server) {
	p.client.Token = server.AccessToken

	for _, c := range server.Connection {
		if c.Local == 0 {
			p.client.URL = c.URI
			break
		}
	}

	p.server = server
}
