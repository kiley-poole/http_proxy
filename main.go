package main

import (
	"syscall"
	"time"
)

func main() {
	serverAddr := syscall.SockaddrInet4{Port: 8085, Addr: [4]byte{127, 0, 0, 1}}
	forwardAddr := syscall.SockaddrInet4{Port: 9000}

	recvSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	check(err)

	sendSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	check(err)

	defer syscall.Close(sendSocket)
	defer syscall.Close(recvSocket)

	err = syscall.Bind(recvSocket, &serverAddr)
	check(err)

	err = syscall.Listen(recvSocket, syscall.SOMAXCONN)
	check(err)

	acceptSocket, _, err := syscall.Accept(recvSocket)
	check(err)

	err = syscall.Connect(sendSocket, &forwardAddr)
	check(err)

	res := make([]byte, 2048)

	_, _, err = syscall.Recvfrom(acceptSocket, res, 0)
	check(err)

	err = syscall.Sendto(sendSocket, res, 0, &forwardAddr)
	check(err)

	time.Sleep(1 * time.Second)

	_, _, err = syscall.Recvfrom(sendSocket, res, 0)
	check(err)

	err = syscall.Sendto(acceptSocket, res, 0, &forwardAddr)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
