package utils

import "net"

func GetAvailablePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func(l *net.TCPListener) {
		err = l.Close()
		if err != nil {

		}
	}(l)
	return l.Addr().(*net.TCPAddr).Port, nil
}
