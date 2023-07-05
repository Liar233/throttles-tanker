package command

import "github.com/Liar233/throttles-tank/internal/server/service"

type FlushDto struct {
	Login   string
	Address string
}

type FlushCommand struct {
	tanker service.TankerInterface
}

func (fc *FlushCommand) Exec(dto FlushDto) error {

	fc.tanker.Reset(dto.Address, dto.Login)

	return nil
}

func NewFlushCommand(tanker service.TankerInterface) *FlushCommand {

	return &FlushCommand{
		tanker: tanker,
	}
}
