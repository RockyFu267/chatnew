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

type ClientChInfo struct {
	Ch   ClientChan
	Name string
}

var InfoList []ClientInfo
var InfoChList []ClientChInfo

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
		go handleConn(tmpinfo)
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
func handleConn(tmpinfo ClientInfo) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(tmpinfo.ConnChan, ch)
	var infoChTmp ClientChInfo
	infoChTmp.Ch = ch
	infoChTmp.Name = tmpinfo.Name
	InfoChList = append(InfoChList, infoChTmp)
	ch <- "You are " + tmpinfo.Name
	messages <- tmpinfo.Name + " has arrived"
	entering <- ch

	input := bufio.NewScanner(tmpinfo.ConnChan)
	for input.Scan() {
		if len(input.Text()) == 0 {
			messages <- tmpinfo.Name + ": " + input.Text()
			continue
		}
		if string(input.Text())[0] == '@' {
			strtmp := stringToDestinationAddr(input.Text())
			contenttmp := stringToDestinationContent(input.Text())
			fmt.Println(strtmp)
			for k := range InfoChList {
				if strtmp == InfoChList[k].Name {
					InfoChList[k].Ch <- tmpinfo.Name + "悄悄对你说: " + contenttmp
				}
			}
		} else {
			messages <- tmpinfo.Name + ": " + input.Text()
		}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <- tmpinfo.Name + " has left"
	tmpinfo.ConnChan.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//截取@的地址
func stringToDestinationAddr(input string) (output string) {
	for k := range string(input) {
		if string(input[k]) == " " {
			output = string(input)[1:k]
			break
		}

	}
	return output
}

//截取@的内容
func stringToDestinationContent(input string) (output string) {
	for k := range string(input) {
		if string(input[k]) == " " {
			output = string(input)[k+1:]
			break
		}

	}
	return output
}
