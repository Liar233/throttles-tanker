package command

import (
	"net"

	"github.com/Liar233/throttles-tank/internal/server/service"
)

type CheckDto struct {
	Address  string
	Login    string
	Password string
}

type CheckCommand struct {
	whiteList service.ListManagerInterface
	blackList service.ListManagerInterface
	tanker    service.TankerInterface
}

func (cc *CheckCommand) Exec(dto CheckDto) error {

	ip := net.ParseIP(dto.Address)

	if cc.blackList.Check(ip) {

		return EntryForbiddenError
	}

	if cc.whiteList.Check(ip) {

		return nil
	}

	if !cc.tanker.Check(dto.Address, dto.Login, dto.Password) {

		return EntryForbiddenError
	}

	return nil
}

func NewCheckCommand(
	whiteList service.ListManagerInterface,
	blackList service.ListManagerInterface,
	tanker service.TankerInterface,
) *CheckCommand {

	return &CheckCommand{
		whiteList: whiteList,
		blackList: blackList,
		tanker:    tanker,
	}
}
