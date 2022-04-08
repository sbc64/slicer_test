package main

import (
	//"github.com/gobwas/ws"
	"log"
	"net"
	"os"
	"reflect"
	"syscall"
)

// GetFdFromConn get net.Conn's file descriptor.
/*
func GetFdFromConn(l net.Conn) int {
	log.Println("Here")
	v := reflect.ValueOf(l)
	log.Printf("2 %s\n", v)
	netFD := reflect.Indirect(reflect.Indirect(v).FieldByName("fd"))
	log.Printf("3 %+v\n", netFD)
	fd := int(netFD.FieldByName("sysfd").Int())
	log.Println("4")
	return fd
}
*/
func GetFDFromConn(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}

// GetFdFromListener get net.Listener's file descriptor.
func GetFdFromListener(l net.Listener) int {
	v := reflect.ValueOf(l)
	netFD := reflect.Indirect(reflect.Indirect(v).FieldByName("fd"))
	fd := int(netFD.FieldByName("sysfd").Int())
	return fd
}

/*
first convert the fd to *os.File via os.NewFile,
and then convert the *os.File to a net.Conn via net.FileConn.
*/

func getCopyConn(originalConn *net.Conn) (net.Conn, error) {
	oldFd := GetFDFromConn(*originalConn)
	newFd, err := syscall.Dup(oldFd)
	if err != nil {
		return nil, err
	}
	connFile := os.NewFile(uintptr(newFd), "websocketconn")
	return net.FileConn(connFile)
}

func handleConn(conn net.Conn) {
	relayAddress := "127.0.0.1:5000"
	relayConn, err := net.ResolveTCPAddr("tcp", relayAddress)
	if err != nil {
		panic(err)
	}
	relayTcpConn, err := net.DialTCP("tcp", nil, relayConn)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to upstream")
	if err != nil {
		panic(err)
	}
	// Consider replacing this with net.Pipe
	go func() {
		readAmount, err := relayTcpConn.ReadFrom(conn)
		log.Printf("Read amount: %d\n", readAmount)
		if err != nil {
			panic(err)
		}
	}()
	defer relayTcpConn.Close()

	go func() {
		readAmount, err := conn.(*net.TCPConn).ReadFrom(relayTcpConn)
		log.Printf("Read amount 2: %d\n", readAmount)
		if err != nil {
			panic(err)
		}
	}()
	defer conn.Close()
}

func main() {
	listenAddr := ":8081"
	log.Printf("Launching listener on: %s\n", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	for {
		log.Printf("Waiting...")
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		log.Printf("Accepted New Connection")
		//newFd, err := getCopyConn(&conn)
		if err != nil {
			panic("2" + err.Error())
		}
		go handleConn(conn)

		//log.Println(newFd)

	}
}
