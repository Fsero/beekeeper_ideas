package main
import "net"
import "fmt"
import "os"
import "time"

func main () {
	client_port := 10000
	pot_port := 20000
	ln, err := net.Listen("tcp",fmt.Sprintf(":%d",client_port))
	if err != nil {
		fmt.Errorf("something went wrong: %s",err)
		os.Exit(1)
	}
	defer ln.Close()
	duration, err := time.ParseDuration("5s")
	backendConnection, err := net.DialTimeout("tcp",fmt.Sprintf(":%d",pot_port), duration)
	handleError(err,100)

	for {
		conn, err := ln.Accept()
		handleError(err,101)
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
	go connect_two_sockets(conn,backend)
}

func connect_two_sockets(client net.Conn, pot net.Conn) {
	var client_data = make([]byte,256)
	var pot_data = make([]byte,256)


	for {

		_, err := client.Write(client_data)
		handleError(err, 104)

		_, err = pot.Read(client_data)
		handleError(err, 105)
		now := time.Now()
		now = now.Add(time.Second * 5)
		pot.SetReadDeadline(now)
		fmt.Printf("Reading data from pot\n")
		_ , err = pot.Read(pot_data)
		handleError(err, 106)
		fmt.Printf("Sending client data from pot\n")
		_ , err = client.Write(pot_data)

		handleError(err, 107)


	}
}
