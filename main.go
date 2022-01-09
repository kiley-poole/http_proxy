package main

import (
	"syscall"
)

func main() {
	serverAddr := syscall.SockaddrInet4{Port: 8085, Addr: [4]byte{127, 0, 0, 1}}
	forwardAddr := syscall.SockaddrInet4{Port: 9000}

	inboundSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	check(err, inboundSocket)

	forwardSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	check(err, forwardSocket)

	defer syscall.Close(forwardSocket)
	defer syscall.Close(inboundSocket)

	err = syscall.Bind(inboundSocket, &serverAddr)
	check(err, inboundSocket)

	err = syscall.Listen(inboundSocket, syscall.SOMAXCONN)
	check(err, inboundSocket)

	newSock, _, err := syscall.Accept(inboundSocket)
	check(err, inboundSocket)

	res := make([]byte, 2048)

	_, _, err = syscall.Recvfrom(newSock, res, 0)
	check(err, newSock)

	err = syscall.Connect(forwardSocket, &forwardAddr)
	check(err, forwardSocket)

	err = syscall.Sendto(forwardSocket, res, 0, &forwardAddr)
	check(err, forwardSocket)

	_, _, err = syscall.Recvfrom(forwardSocket, res, 0)
	check(err, newSock)

	err = syscall.Sendto(newSock, res, 0, &forwardAddr)
	check(err, forwardSocket)
}

func check(err error, socket int) {
	if err != nil {
		syscall.Close(socket)
		panic(err)
	}
}
