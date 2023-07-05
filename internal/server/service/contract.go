package service

import "net"

type TankerInterface interface {
	Check(address, login, password string) bool
	Reset(address, login string)
}

type ListManagerInterface interface {
	Check(ip net.IP) bool
	FillCheckList(items []net.IPNet)
}
