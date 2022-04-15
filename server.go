package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
	"unicode/utf8"
)

type HandleConnection func(c net.Conn) error

func main() {
	fmt.Println("=== サーバーを起動します ===")

	handleConn := func(c net.Conn) error {
		buf := make([]byte, 1024)
		fmt.Println("handle Conn")
		defer fmt.Println("handle End")
		for {
			n, err := c.Read(buf)
			fmt.Println(n, err)
			if err != nil {
				fmt.Println("c read err", err)
				return err
			}
			if n == 0 {
				fmt.Println("0")
				break
			}
			s := string(buf[:n])
			fmt.Println(s)
			if string(buf[n-4:n]) == "\r\n\r\n" {
				break
			}
		}

		response_body := "<html><body><h1>It works!</h1></body></html>\r\n"
		response_line := "HTTP/1.1 200 OK \r\n"
		response_header := ""
		response_header += "Date: " + string(time.Now().Format("Mon, 2 Jan 2006 15:04:05 GMT")) + "\r\n"
		response_header += "Host: FlflServer/0.1\r\n"
		response_header += "Content-Length: " + strconv.Itoa(utf8.RuneCountInString(response_body)) + "\r\n"
		response_header += "Connection: Close\r\n"
		response_header += "Content-Type: text/html\r\n"

		response := response_line + response_header + "\r\n" + response_body
		fmt.Println("==response==")
		fmt.Println(response)
		fmt.Println("===")

		fmt.Fprint(c, response)
		return nil
	}

	if err := start(handleConn); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func start(f HandleConnection) error {
	ln, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		return err
	}
	defer ln.Close()
	fmt.Println("accept!")
	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()
	for {
		err := f(conn)
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			return nil
		}
	}
}
