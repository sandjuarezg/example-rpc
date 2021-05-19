package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

var rConn *bufio.Reader
var rStdin = bufio.NewReader(os.Stdin)

func main() {
	if len(os.Args) >= 4 {
		var wg = sync.WaitGroup{}
		var connMyServer = fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])
		var connOtherServer = fmt.Sprintf("%s:%s", os.Args[1], os.Args[3])
		var reply = make([]byte, 1024)

		conn, err := net.Dial("tcp", connOtherServer)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		rConn = bufio.NewReader(conn)

		ln, err := net.Listen("tcp", connMyServer)
		if err != nil {
			log.Fatal(err)
		}
		defer ln.Close()

		fmt.Printf("Listening on %s\n", connMyServer)

		wg.Add(1)
		go func(ln net.Listener, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				var _, err = ln.Accept()
				if err != nil {
					log.Fatal(err)
				}

				//Make request
			}
		}(ln, &wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				//Read from standard input
				var n, err = rStdin.Read(reply)
				if err != nil {
					var _, ok = err.(*net.OpError)
					if ok || err == io.EOF {
						break
					} else {
						log.Fatal(err)
					}
				}

				//Write on connection
				_, err = conn.Write([]byte(reply[0:n]))
				if err != nil {
					log.Fatal(err)
				}
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				//Read message host
				_, err = rConn.Read(reply)
				if err != nil {
					var _, ok = err.(*net.OpError)
					if ok || err == io.EOF {
						fmt.Println("Server closed")
						os.Exit(1)
					} else {
						log.Fatal(err)
					}
				}
				fmt.Printf("%s", reply)
			}
		}(&wg)

		wg.Wait()

	}
}
