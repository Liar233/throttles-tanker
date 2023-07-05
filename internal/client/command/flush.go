package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Liar233/throttles-tank/pkg/protocol"
	"github.com/urfave/cli/v2"
)

func NewFlushCliCommand(client *http.Client, servUrl *string) *cli.Command {

	var login, ipStr string

	return &cli.Command{
		Name:    "flush",
		Aliases: []string{"f"},
		Usage:   "flush bucket",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "login",
				Usage:       "",
				Destination: &login,
				Aliases:     []string{"l"},
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "ip",
				Usage:       "",
				Destination: &ipStr,
				Aliases:     []string{"i"},
				Required:    true,
			},
		},
		Action: func(cCtx *cli.Context) error {

			subnet := cCtx.Args().First()

			uri := fmt.Sprintf("%s/flush", *servUrl)

			reqDto := &protocol.AddRuleRequest{Subnet: subnet}

			data, err := json.Marshal(reqDto)

			if err != nil {

				return err
			}

			reader := bytes.NewReader(data)

			req, err := http.NewRequest(http.MethodPost, uri, reader)

			if err != nil {

				return err
			}

			resp, err := client.Do(req)

			if err != nil {

				return err
			}

			if resp.StatusCode != http.StatusOK {

				return fmt.Errorf("server respose with code %d", resp.StatusCode)
			}

			return nil
		},
	}
}
