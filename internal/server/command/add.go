package command

import (
	"errors"
	"net"

	"github.com/Liar233/throttles-tank/internal/server/service"
	"github.com/Liar233/throttles-tank/internal/server/storage"
)

var InvalidRuleTypeError = errors.New("invalid rule type")

type AddRuleDto struct {
	Subnet   net.IPNet
	RuleType string
}

type AddRuleCommand struct {
	ruleStorage storage.RuleStorageInterface
	whiteList   service.ListManagerInterface
	blackList   service.ListManagerInterface
}

func (arc *AddRuleCommand) Exec(dto AddRuleDto) error {

	if err := arc.ruleStorage.Add(dto.RuleType, dto.Subnet); err != nil {

		if err == storage.RuleAlreadyExistsDBError {

			return RuleAlreadyExistsError
		}

		return err
	}

	updatedNetList, err := arc.ruleStorage.GetList(dto.RuleType)

	if err != nil {

		return err
	}

	switch dto.RuleType {
	case "white":

		arc.whiteList.FillCheckList(updatedNetList)

	case "black":

		arc.blackList.FillCheckList(updatedNetList)

	default:

		return InvalidRuleTypeError
	}

	return nil
}

func NewAddRuleCommand(
	ruleStorage storage.RuleStorageInterface,
	whiteList service.ListManagerInterface,
	blackList service.ListManagerInterface,
) *AddRuleCommand {

	return &AddRuleCommand{
		ruleStorage: ruleStorage,
		whiteList:   whiteList,
		blackList:   blackList,
	}
}
