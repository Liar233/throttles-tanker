package service

import (
	"net"
	"sync/atomic"

	"github.com/Liar233/throttles-tank/pkg/ip_tree"
)

type ListManager struct {
	checkTree atomic.Pointer[ip_tree.IpLeafNode]
}

func (lm *ListManager) FillCheckList(items []net.IPNet) {

	newTree := &ip_tree.IpLeafNode{}

	for _, item := range items {

		newTree.Insert(&item, nil)
	}

	lm.checkTree.Store(newTree)
}

func (lm *ListManager) Check(ip net.IP) bool {

	chkTree := lm.checkTree.Load()

	if chkTree == nil {

		return false
	}

	return chkTree.Check(ip)
}

func NewListManager() *ListManager {

	return &ListManager{
		checkTree: atomic.Pointer[ip_tree.IpLeafNode]{},
	}
}
