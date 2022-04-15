package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type HandleConnection func(c net.Conn) error

func main() {
	fmt.Println("=== サーバーを起動します ===")

	handleConn := func(c net.Conn) error {
		buf := make([]byte,1024)
		for {
			n,err:=c.Read(buf)
			if err != nil {
				return err
			}
			if n == 0 {
				break
			}
			s := string(buf[:n])
			fmt.Println(s)
			fmt.Fprintf(c,"accept:%s\n", s)
		}
		return nil
	}

	if err := start(handleConn); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func start(f HandleConnection) error {
	ln, err := net.Listen("tcp","0.0.0.0:8080")
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
			fmt.Println("==データここまで==")
			return err
		}
		if err == io.EOF {
			return nil
		}
	}
}