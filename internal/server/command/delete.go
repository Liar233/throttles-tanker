package command

import (
	"net"

	"github.com/Liar233/throttles-tank/internal/server/service"
	"github.com/Liar233/throttles-tank/internal/server/storage"
)

type DeleteRuleDto struct {
	Subnet   net.IPNet
	RuleType string
}

type DeleteRuleCommand struct {
	ruleStorage storage.RuleStorageInterface
	whiteList   service.ListManagerInterface
	blackList   service.ListManagerInterface
}

func (drc *DeleteRuleCommand) Exec(dto DeleteRuleDto) error {

	if err := drc.ruleStorage.Delete(dto.Subnet); err != nil {

		return RuleNotFoundError
	}

	updatedNetList, err := drc.ruleStorage.GetList(dto.RuleType)

	if err != nil {

		return err
	}

	switch dto.RuleType {
	case "white":

		drc.whiteList.FillCheckList(updatedNetList)

	case "black":

		drc.blackList.FillCheckList(updatedNetList)

	default:

		return InvalidRuleTypeError
	}

	return nil
}

func NewDeleteRuleCommand(
	ruleStorage storage.RuleStorageInterface,
	whiteList service.ListManagerInterface,
	blackList service.ListManagerInterface,
) *DeleteRuleCommand {

	return &DeleteRuleCommand{
		ruleStorage: ruleStorage,
		whiteList:   whiteList,
		blackList:   blackList,
	}
}
