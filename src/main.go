package main
import "net"
import "fmt"
import "os"
import (
	"time"
	"io"
)

func main () {
	client_port := 22
	pot_port := 20000
	ln, err := net.Listen("tcp",fmt.Sprintf(":%d",client_port))
	handleError(err,1)
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		handleError(err,101)
		duration, err := time.ParseDuration("5s")
		//TODO: Replace this call for launching the docker interface.
		backendConnection, err := net.DialTimeout("tcp",fmt.Sprintf(":%d",pot_port), duration)
		handleError(err,100)

		go handleConnection(conn, backendConnection)
	}
}

func handleError(err error, exit_code int) {
	if err != nil {
		fmt.Printf("something went wrong: %s\n",err)
		os.Exit(exit_code)
	}

}

func handleConnection(conn net.Conn, backend net.Conn) {

	copyConn:=func(writer, reader net.Conn) {
		_, err:= io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(backend, conn)
	go copyConn(conn, backend)
}



