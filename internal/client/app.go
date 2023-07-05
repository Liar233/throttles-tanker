package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Liar233/throttles-tank/internal/client/command"
	"github.com/urfave/cli/v2"
)

var server, servUrl string

func NewClient() *cli.App {

	httpClient := &http.Client{
		Timeout: time.Second * time.Duration(5),
	}

	app := &cli.App{
		Name:        Name,
		HelpName:    HelpName,
		Usage:       Usage,
		UsageText:   UsageText,
		Version:     Version,
		Description: Description,
		Commands: []*cli.Command{
			command.NewAddCliCommand(httpClient, &servUrl),
			command.NewDeleteCliCommand(httpClient, &servUrl),
			command.NewFlushCliCommand(httpClient, &servUrl),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "server",
				Usage:       "server address",
				Required:    true,
				Value:       "",
				Destination: &server,
				Aliases:     []string{"s"},
				Action:      nil,
			},
		},
		Before: func(context *cli.Context) error {

			var err error

			buf := fmt.Sprintf("http://%s", server)

			servUrl = buf

			if err != nil {
				return fmt.Errorf("server address \"%s\" invalid", server)
			}

			return nil
		},
		Compiled: time.Now(),
	}

	return app
}

const Version = "1.0"
const Name = "tanker-cli"
const HelpName = "tanker-cli"
const Usage = `"throttles tanker client"`

const UsageText = `tanker-client --server [addr...] [command] [command options...] `
const Description = `throttle tanker client provides you add/delete rules for blacklist, whitelist, flush buckets`
