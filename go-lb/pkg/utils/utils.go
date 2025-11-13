package utils

import (
	"net"
)

func IsValidIPPort(addr string) bool {
	_, err := net.ResolveTCPAddr("tcp", addr)
	return err == nil
}
