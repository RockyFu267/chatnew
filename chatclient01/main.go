package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			fmt.Println(err)
			return
		}
	}()
	if _, err := io.Copy(conn, os.Stdin); err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(10 * time.Second)
}
