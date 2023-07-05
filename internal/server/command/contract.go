package command

import "errors"

var EntryForbiddenError = errors.New("entry forbidden")
var RuleNotFoundError = errors.New("rule not found")
var RuleAlreadyExistsError = errors.New("rule already exists")

type ServerCommandInterface[T CheckDto | AddRuleDto | DeleteRuleDto | FlushDto] interface {
	Exec(T) error
}
