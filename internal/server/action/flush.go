package action

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/Liar233/throttles-tank/internal/server/command"
	"github.com/Liar233/throttles-tank/pkg/protocol"
	log "github.com/sirupsen/logrus"
)

type FlushAction struct {
	cmd *command.FlushCommand
}

func (fa *FlushAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	request := &protocol.FlushRequest{}

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if net.ParseIP(request.Ip) == nil {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if request.Login == "" {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	dto := command.FlushDto{
		Login:   request.Login,
		Address: request.Ip,
	}

	if err := fa.cmd.Exec(dto); err != nil {

		log.Errorf("error flushing \"%s %s\" with: %s", request.Login, request.Ip, err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewFlushAction(cmd *command.FlushCommand) *FlushAction {

	return &FlushAction{cmd: cmd}
}
