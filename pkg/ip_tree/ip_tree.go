package ip_tree

import "net"

type IpLeafNode struct {
	parent *IpLeafNode
	leaves []*IpLeafNode

	oct uint8

	subnet *net.IPNet
}

func (node *IpLeafNode) Insert(subnet *net.IPNet, unmasked []byte) {

	if node.parent == nil {

		unmasked = GetUnmaskedIp(subnet)
	}

	if len(unmasked) == 0 {

		node.subnet = subnet

		return
	}

	if node.leaves == nil {

		node.leaves = make([]*IpLeafNode, 0)
	}

	for _, leaf := range node.leaves {

		if leaf.oct == unmasked[0] {

			leaf.Insert(subnet, unmasked[1:])

			return
		}
	}

	newLeaf := NewIpLeafNode(node, unmasked[0])

	node.leaves = append(node.leaves, newLeaf)

	newLeaf.Insert(subnet, unmasked[1:])
}

func (node *IpLeafNode) Check(ip net.IP) bool {

	if node.subnet != nil {

		return node.subnet.Contains(ip)
	}

	for _, leaf := range node.leaves {

		if leaf.Check(ip) {

			return true
		}
	}

	return false
}

func NewIpLeafNode(parent *IpLeafNode, oct uint8) *IpLeafNode {

	return &IpLeafNode{
		parent: parent,
		oct:    oct,
	}
}

func GetUnmaskedIp(subnet *net.IPNet) []byte {

	addr := subnet.IP.To4()

	for i := 0; i < len(addr); i++ {

		if subnet.Mask[i] != 255 {

			return addr[:i]
		}
	}

	return addr
}
