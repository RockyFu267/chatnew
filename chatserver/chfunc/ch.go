package chfunc

import (
	"bufio"
	Pt "chatserver/publictype"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

//BecomeUpper 变大写
func BecomeUpper(c net.Conn) {
	defer c.Close()
	input := bufio.NewScanner(c)
	for input.Scan() {
		fmt.Fprintln(c, "\t", strings.ToUpper(input.Text()))
	}
}

//PrintListName debug时候用的
func PrintListName() {
	var strlist []string
	for k := range Pt.InfoChList {
		strlist = append(strlist, Pt.InfoChList[k].Name)
		res2B, _ := json.Marshal(strlist)
		fmt.Println(string(res2B))
		time.Sleep(1 * time.Second)
	}
}

//PrintListAddress debug时候用的
func PrintListAddress() {
	var strlist []string
	for k := range Pt.InfoChList {
		strlist = append(strlist, Pt.InfoChList[k].Address)
		res2B, _ := json.Marshal(strlist)
		fmt.Println(string(res2B))
		time.Sleep(1 * time.Second)
	}
}

//ReturnTime 输出时间
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

//Printint debug时候用的
func Printint() {
	for k := 0; k < 1000000; k++ {
		fmt.Println(k)
		time.Sleep(1 * time.Second)
	}
}

//Less5SecondEchoEachOther 俩启动间隔小于五秒互相通信
func Less5SecondEchoEachOther(tmpinfo Pt.ClientInfo) {
	time.Sleep(5 * time.Second)
	for k := range Pt.InfoList {
		if Pt.InfoList[k].Address != tmpinfo.Address {
			input := bufio.NewScanner(tmpinfo.ConnChan)
			for input.Scan() {
				fmt.Fprintln(Pt.InfoList[k].ConnChan, "\t", strings.ToUpper(input.Text()))
			}
		}
	}
}

//ClientWriter 把管道数据写入tcp连接
func ClientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
