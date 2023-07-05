package service

import (
	"net"
	"testing"
)

func TestCheckEmpty(t *testing.T) {

	listManager := NewListManager()

	ip := net.ParseIP("127.0.0.1")

	if listManager.Check(ip) {

		t.Error("failed check empty list manager")
	}
}

func TestFillCheckList(t *testing.T) {

	chkList1Str := []string{
		"10.10.10.100/32",
		"169.172.100.0/24",
		"192.168.0.0/16",
		"172.0.0.0/8",
	}

	chkList2Str := []string{
		"10.10.10.0/24",
		"169.172.100.0/24",
		"192.168.0.0/16",
		"172.0.0.0/8",
	}

	var chkList1, chkList2 []*net.IPNet
	var err error

	if chkList1, err = makeNetList(chkList1Str); err != nil {

		t.Fatal("failed test data")
	}

	if chkList2, err = makeNetList(chkList2Str); err != nil {

		t.Fatal("failed test data")
	}

	ip := net.ParseIP("10.10.10.101")

	listManager := NewListManager()

	listManager.FillCheckList(chkList1)

	ptr1 := listManager.checkTree.Load()

	if listManager.Check(ip) {

		t.Error("failed check ip")
	}

	listManager.FillCheckList(chkList2)

	ptr2 := listManager.checkTree.Load()

	if !listManager.Check(ip) {

		t.Error("failed check ip after filling list manager")
	}

	if ptr1 == ptr2 {

		t.Error("failed same pointer after filling")
	}
}

func makeNetList(CIDRs []string) ([]*net.IPNet, error) {

	subnets := make([]*net.IPNet, 0)

	for _, CIDR := range CIDRs {

		_, ipnet, err := net.ParseCIDR(CIDR)

		if err != nil {

			return nil, err
		}

		subnets = append(subnets, ipnet)
	}

	return subnets, nil
}
