package action

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/Liar233/throttles-tank/internal/server/command"
	"github.com/Liar233/throttles-tank/pkg/protocol"
	log "github.com/sirupsen/logrus"
)

type DeleteRuleAction struct {
	listType string
	cmd      *command.DeleteRuleCommand
}

func (dra *DeleteRuleAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	request := &protocol.DeleteRuleRequest{}

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	_, subnet, err := net.ParseCIDR(request.Subnet)

	if err != nil {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	dto := command.DeleteRuleDto{
		Subnet:   *subnet,
		RuleType: dra.listType,
	}

	err = dra.cmd.Exec(dto)

	if err == command.InvalidRuleTypeError {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err == command.RuleNotFoundError {

		http.Error(w, "", http.StatusNotFound)
		return
	}

	if err != nil {

		log.Errorf("error deleting rule\"%s\" with: %s", dto.Subnet.String(), err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewDeleteRuleAction(listType string, cmd *command.DeleteRuleCommand) *DeleteRuleAction {

	return &DeleteRuleAction{
		listType: listType,
		cmd:      cmd,
	}
}
