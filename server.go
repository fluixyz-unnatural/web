package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type HandleConnection func(c net.Conn) error

func main() {
	fmt.Println("=== サーバーを起動します ===")

	BASE_DIR, err := filepath.Abs("./")
	if err != nil {
		return
	}
	STATIC_ROOT := BASE_DIR + "/static"

	handleConn := func(c net.Conn) error {
		buf := make([]byte, 1024)
		s := ""
		for {
			n, err := c.Read(buf)
			if err != nil {
				return err
			}
			if n == 0 {
				break
			}
			s += string(buf[:n])
			if string(buf[n-4:n]) == "\r\n\r\n" {
				break
			}
		}
		// リクエストライン
		request_line := strings.Split(strings.Split(s, "\r\n")[0], " ")
		// request_method := request_line[0]
		request_path := request_line[1]
		// http_version := request_line[2]

		// ヘッダー
		// request_header := strings.Split(strings.Split(s, "\r\n\r\n")[0], "\r\n")[1:]

		// ボディ
		// request_body := strings.Split(s, "\r\n\r\n")[1]

		static_file_path := (STATIC_ROOT + request_path)

		response_line := "HTTP/1.1 200 OK \r\n"
		response_header := ""
		response_header += "Date: " + string(time.Now().Format("Mon, 2 Jan 2006 15:04:05 GMT")) + "\r\n"
		response_header += "Host: FlflServer/0.1\r\n"
		response_header += "Connection: Close\r\n"

		response_body := ""

		bytes, err := ioutil.ReadFile(static_file_path)
		if err != nil {
			response_line = "HTTP/1.1 404 Not Found \r\n"
			response_body = "<html><body><h1>404 Not Found</h1></body></html>\r\n"
		} else {
			response_body = string(bytes)
		}

		response_header += "Content-Length: " + strconv.Itoa(len(response_body)) + "\r\n"
		response_header += "Content-Type: " + memetypes(static_file_path) + "\r\n"

		response := response_line + response_header + "\r\n" + response_body

		// fmt.Println("==response==")
		// fmt.Println(response)
		// fmt.Println("===")

		fmt.Fprint(c, response)
		return nil
	}
	for {
		if err := start(handleConn); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func start(f HandleConnection) error {
	fmt.Println("listening... ")
	ln, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		return err
	}
	defer ln.Close()
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

func memetypes(filename string) string {
	switch parts := strings.Split(filename, "."); parts[len(parts)-1] {
	case "html":
		return "text/html"
	case "css":
		return "text/css"
	case "webp":
		return "image/webp"
	case "png":
		return "image/png"
	case "jpg":
		return "image/jpg"
	}
	return "text/html"
}
