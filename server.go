package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
)

type HandleConnection func(c net.Conn) error

func main() {
	fmt.Println("=== サーバーを起動します ===")

	handleConn := func(c net.Conn) error {
		buf := make([]byte,1024)
		fmt.Println("handle Conn")
		defer fmt.Println("handle End")
		for {
			n,err:=c.Read(buf)
			fmt.Println(n,err)
			if err != nil {
				fmt.Println("c read err",err)
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
		f, err := ioutil.ReadFile("./response.txt")
		fmt.Println("open response.txt")
		if err != nil {
			fmt.Println("err", err)
		}
		fmt.Fprintln(c,string(f))
		return io.EOF
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
			return err
		}
		if err == io.EOF {
			return nil
		}
	}
}