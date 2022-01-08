package main

import (
	"syscall"
)

func main() {
	netAddr := syscall.SockaddrInet4{Port: 8085, Addr: [4]byte{127, 0, 0, 1}}

	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	check(err, socket)

	err = syscall.Bind(socket, &netAddr)
	check(err, socket)

	err = syscall.Listen(socket, syscall.SOMAXCONN)
	check(err, socket)

	newSock, incAddr, err := syscall.Accept(socket)
	check(err, socket)

	res := make([]byte, 512)
	_, _, err = syscall.Recvfrom(newSock, res, 0)
	check(err, socket)

	err = syscall.Sendto(newSock, res, 0, incAddr)
	check(err, socket)

	err = syscall.Close(socket)
	check(err, socket)
}

func check(err error, socket int) {
	if err != nil {
		syscall.Shutdown(socket, 2)
		panic(err)
	}
}
