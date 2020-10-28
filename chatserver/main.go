package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type ClientChan chan string

var AllClients = make(map[string]net.Conn)

type ClientInfo struct {
	ConnChan net.Conn
	Name     string
}

var InfoList []ClientInfo

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("fuck")
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		who := conn.RemoteAddr().String()
		fmt.Println(who, "已经建立连接")
		var tmpinfo ClientInfo
		tmpinfo.ConnChan = conn
		tmpinfo.Name = who
		InfoList = append(InfoList, tmpinfo)
		go handleConn(conn)
		// go Less5SecondEchoEachOther(tmpinfo) //俩启动间隔小于五秒互相通信
		//go BecomeUpper(conn)			//变大写
		// go ReturnTime(conn)          //输出时间
	}
}

//俩启动间隔小于五秒互相通信
func Less5SecondEchoEachOther(tmpinfo ClientInfo) {
	time.Sleep(5 * time.Second)
	for k := range InfoList {
		if InfoList[k].Name != tmpinfo.Name {
			input := bufio.NewScanner(tmpinfo.ConnChan)
			for input.Scan() {
				fmt.Fprintln(InfoList[k].ConnChan, "\t", strings.ToUpper(input.Text()))
			}
		}
	}
}

//输出时间
func ReturnTime(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

//变大写
func BecomeUpper(c net.Conn) {
	defer c.Close()
	input := bufio.NewScanner(c)
	for input.Scan() {
		fmt.Fprintln(c, "\t", strings.ToUpper(input.Text()))
	}
}

var (
	entering = make(chan ClientChan)
	leaving  = make(chan ClientChan)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[ClientChan]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
