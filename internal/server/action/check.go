package action

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/Liar233/throttles-tank/internal/server/command"
	"github.com/Liar233/throttles-tank/pkg/protocol"
	log "github.com/sirupsen/logrus"
)

type CheckAction struct {
	cmd *command.CheckCommand
}

func (ca *CheckAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	request := &protocol.CheckRequest{}

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if net.ParseIP(request.Ip) == nil {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if request.Login == "" || request.Password == "" {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	dto := command.CheckDto{
		Address:  request.Ip,
		Login:    request.Login,
		Password: request.Password,
	}

	err := ca.cmd.Exec(dto)

	if err == command.EntryForbiddenError {

		http.Error(w, "", http.StatusTooManyRequests)
		return
	}

	if err != nil {

		log.Errorf("error checking \"%s %s\" with: %s", request.Login, request.Ip, err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewCheckAction(cmd *command.CheckCommand) *CheckAction {

	return &CheckAction{cmd: cmd}
}
