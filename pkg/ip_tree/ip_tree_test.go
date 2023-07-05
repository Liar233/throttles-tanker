package ip_tree

import (
	"net"
	"testing"
)

var CIDRs = []string{
	"10.10.10.255/8",
	"10.10.10.255/16",
	"10.10.10.255/24",
	"10.10.10.255/32",
	"10.0.0.255/8",
	"10.0.0.255/16",
	"10.0.0.255/24",
	"10.0.0.255/32",
	"169.172.111.255/8",
	"169.172.111.255/16",
	"169.172.111.255/24",
	"169.172.111.255/32",
	"172.169.222.255/8",
	"172.169.222.255/16",
	"172.169.222.255/24",
	"172.169.222.255/32",
	"192.168.100.255/8",
	"192.168.100.255/16",
	"192.168.100.255/24",
	"192.168.100.255/32",
}

func TestIpTree(t *testing.T) {

	tree := NewIpLeafNode(nil, 0)

	for _, cidr := range CIDRs {

		_, ipNet, _ := net.ParseCIDR(cidr)

		tree.Insert(ipNet, nil)
	}

	ip := net.ParseIP("10.10.10.255")

	if !tree.Check(ip) {

		t.Error("failed check ip")
	}

	ip = net.ParseIP("168.177.100.128")

	if tree.Check(ip) {

		t.Error("failed check ip")
	}
}
