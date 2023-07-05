package action

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/Liar233/throttles-tank/internal/server/command"
	"github.com/Liar233/throttles-tank/pkg/protocol"
	log "github.com/sirupsen/logrus"
)

type AddRuleAction struct {
	listType string
	cmd      *command.AddRuleCommand
}

func (ara *AddRuleAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	request := &protocol.AddRuleRequest{}

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	_, subnet, err := net.ParseCIDR(request.Subnet)

	if err != nil {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	dto := command.AddRuleDto{
		Subnet:   *subnet,
		RuleType: ara.listType,
	}

	err = ara.cmd.Exec(dto)

	if err == command.InvalidRuleTypeError {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err != nil {

		if err == command.RuleAlreadyExistsError {

			http.Error(w, "", http.StatusConflict)
			return
		}

		log.Errorf("error adding rule\"%s\" with: %s", dto.Subnet.String(), err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func NewAddRuleAction(listType string, cmd *command.AddRuleCommand) *AddRuleAction {

	return &AddRuleAction{
		listType: listType,
		cmd:      cmd,
	}
}
