package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) >= 3 {
		var conn_sever = fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])
		var ln, err = net.Listen("tcp", conn_sever)
		if err != nil {
			log.Fatal(err)
		}
		defer ln.Close()

		fmt.Printf("Listening on %s\n", conn_sever)

		for {
			var conn, err = ln.Accept()
			if err != nil {
				log.Fatal(err)
			}

			go handleRequest(conn)
		}
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		var reply = make([]byte, 1024)
		var r = bufio.NewReader(conn)
		var n, err = r.Read(reply)
		if err != nil {
			var _, ok = err.(*net.OpError)
			if ok || err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		var text = reply[0 : n-1]
		var elements = strings.Split(string(text), " ")

		if len(elements) >= 3 {
			x, err := strconv.Atoi(elements[1])
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.Atoi(elements[2])
			if err != nil {
				log.Fatal(err)
			}

			var mess string
			switch elements[0] {
			case "add":
				var res = add(x, y)
				mess = fmt.Sprintf("Result to add: %d\nReply by %s\n", res, conn.LocalAddr())
				conn.Write([]byte(mess))
			case "subtract":
				var res = subtract(x, y)
				mess = fmt.Sprintf("Result to subtract: %d\nReply by %s\n", res, conn.LocalAddr())
				conn.Write([]byte(mess))
			default:
				mess = fmt.Sprintf("Func no found\nReply by %s\n", conn.LocalAddr())
				conn.Write([]byte(mess))
			}

			fmt.Printf("%s petition successful\n", conn.RemoteAddr())
		}
	}
}

func add(x, y int) int {
	return x + y
}

func subtract(x, y int) int {
	return x - y
}
