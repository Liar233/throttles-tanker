package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/Liar233/throttles-tank/pkg/protocol"
	"github.com/urfave/cli/v2"
)

func NewDeleteCliCommand(client *http.Client, servUrl *string) *cli.Command {

	var white, black bool

	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"a"},
		Usage:   "delete subnet from white/black list",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "white",
				Usage:       "",
				Value:       false,
				Destination: &white,
				Aliases:     []string{"w"},
			},
			&cli.BoolFlag{
				Name:        "black",
				Usage:       "",
				Value:       false,
				Destination: &black,
				Aliases:     []string{"b"},
			},
		},
		Before: func(cCtx *cli.Context) error {

			if white == true && white == black {

				return errors.New("subnet can be deleted from white or black lists only")
			}

			if white == false && white == black {

				return errors.New("have to set white or black list")
			}

			return nil
		},
		Action: func(cCtx *cli.Context) error {

			subnet := cCtx.Args().First()

			if subnet == "" {

				return errors.New("subnet not set")
			}

			if _, _, err := net.ParseCIDR(subnet); err != nil {

				return errors.New("subnet CIDR not valid")
			}

			listType := "white"

			if black {
				listType = "black"
			}

			uri := fmt.Sprintf("%s/delete/%s", *servUrl, listType)

			reqDto := &protocol.DeleteRuleRequest{Subnet: subnet}

			data, err := json.Marshal(reqDto)

			if err != nil {

				return err
			}

			reader := bytes.NewReader(data)

			req, err := http.NewRequest(http.MethodDelete, uri, reader)

			if err != nil {

				return err
			}

			resp, err := client.Do(req)

			if err != nil {

				return err
			}

			if resp.StatusCode == http.StatusNotFound {

				return fmt.Errorf("subnet %s not found", subnet)
			}

			if resp.StatusCode != http.StatusOK {

				return fmt.Errorf("server respose with code %d", resp.StatusCode)
			}

			_, _ = fmt.Fprintf(cCtx.App.Writer, "subnet %s deleted successful\n", subnet)

			return nil
		},
	}
}
