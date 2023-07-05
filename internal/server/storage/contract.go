package storage

import (
	"errors"
	"net"
)

const WhiteList = "white"
const BlackList = "black"

var InvalidRuleTypeDBError = errors.New("invalid rule type")
var RuleNotFoundDBError = errors.New("rule not found")
var RuleAlreadyExistsDBError = errors.New("rule already exists")

type RuleStorageInterface interface {
	Add(ruleType string, cidr net.IPNet) error
	Delete(cidr net.IPNet) error
	GetList(ruleType string) ([]net.IPNet, error)
	Connect() error
	Close() error
}
