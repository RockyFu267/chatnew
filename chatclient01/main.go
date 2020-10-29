package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	done := make(chan struct{})
	defer conn.Close()
	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			fmt.Println(err)
			return
		}
		done <- struct{}{}
	}()
	ch := make(chan string)
	go clientWriter(conn, ch)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		switch input.Text() {
		default:
			ch <- input.Text()
		}

	}
	conn.Close()
	<-done

}

//clientWriter 接受所有管道
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
