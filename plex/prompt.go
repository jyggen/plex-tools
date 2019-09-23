package plex

import (
	"github.com/jyggen/promptui"
	"os"
)

func (p *Plex) PromptForLibraryKey() (string, error) {
	libraries, err := p.client.GetLibraries()

	if err != nil {
		return "", err
	}

	var options []string

	for _, l := range libraries.MediaContainer.Directory {
		options = append(options, l.Title)
	}

	prompt := promptui.Select{
		Label:  "Choose a library",
		Items:  options,
		Stdin:  os.Stdin,
		Stdout: os.Stderr,
	}

	_, title, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return p.GetLibraryKeyByTitle(title)
}

func (p *Plex) PromptForServer() (*Server, error) {
	servers, err := p.client.GetServers()

	if err != nil {
		return nil, err
	}

	var options []string

	for _, s := range servers {
		options = append(options, s.Name)
	}

	prompt := promptui.Select{
		Label:  "Choose a server",
		Items:  options,
		Stdin:  os.Stdin,
		Stdout: os.Stderr,
	}

	_, name, err := prompt.Run()

	if err != nil {
		return nil, err
	}

	return p.GetServerByName(name)
}
