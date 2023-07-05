package storage

import (
	"net"
	"sync"
)

type MemoryRuleStorage struct {
	white sync.Map
	black sync.Map
}

func (m *MemoryRuleStorage) Connect() error {

	return nil
}

func (m *MemoryRuleStorage) Close() error {

	return nil
}

func (m *MemoryRuleStorage) Add(ruleType string, cidr net.IPNet) error {

	key := cidr.String()

	switch ruleType {
	case WhiteList:

		m.white.Store(key, cidr)

	case BlackList:

		m.black.Store(key, cidr)

	default:
		return InvalidRuleTypeDBError
	}

	return nil
}

func (m *MemoryRuleStorage) Delete(cidr net.IPNet) error {

	key := cidr.String()

	m.black.Delete(key)
	m.white.Delete(key)

	return RuleNotFoundDBError
}

func (m *MemoryRuleStorage) GetList(ruleType string) ([]net.IPNet, error) {

	result := make([]net.IPNet, 0)

	var source *sync.Map

	switch ruleType {
	case WhiteList:

		source = &m.white

	case BlackList:

		source = &m.black

	default:

		return nil, InvalidRuleTypeDBError
	}

	source.Range(func(key, value any) bool {

		result = append(result, value.(net.IPNet))
		return true
	})

	return result, nil
}

func NewMemoryRuleStorage() *MemoryRuleStorage {
	return &MemoryRuleStorage{}
}
