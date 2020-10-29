package main

import (
	"bufio"
	"encoding/json"
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
	Address  string
	Name     string
}

type ClientChInfo struct {
	Ch      ClientChan
	Address string
	Name    string
}

type ChatGroup struct {
	Name   string
	ChList []ClientChInfo
}

var InfoList []ClientInfo
var InfoChList []ClientChInfo
var RoomList []ChatGroup

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("fuck")
	go broadcaster()
	// go Printint()
	go PrintListAddress()
	go PrintListName()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		who := conn.RemoteAddr().String()
		fmt.Println(who, "已经建立连接")
		var tmpinfo ClientInfo
		tmpinfo.ConnChan = conn
		tmpinfo.Address = who
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
		if InfoList[k].Address != tmpinfo.Address {
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
	infoChTmp.Address = tmpinfo.Address
	InfoChList = append(InfoChList, infoChTmp)
	ch <- "You are " + tmpinfo.Address
	messages <- tmpinfo.Address + " has arrived"
	entering <- ch

	input := bufio.NewScanner(tmpinfo.ConnChan)
	for input.Scan() {
		switch input.Text() {
		case "myname":
			if infoChTmp.Name != "" {
				infoChTmp.Ch <- "你已经输入过昵称：" + tmpinfo.Name
				continue
			}
			infoChTmp.Ch <- infoChTmp.Address + ":输入昵称"
			var myname string
			if input.Scan() {
				myname = input.Text()
			}
			judge := JudgeStringSpecialSymbol(myname)
			if judge == false {
				infoChTmp.Ch <- infoChTmp.Address + ":昵称只支持大小写A-z以及0-9,长度不超过20"
				continue
			}
			var sign bool = false
			for k := range InfoChList {
				if InfoChList[k].Name == myname {
					infoChTmp.Ch <- infoChTmp.Address + ":已被使用,请重新输入昵称"
					sign = true
					break
				}
			}
			if sign == true {
				continue
			}
			infoChTmp.Name = myname
			tmpinfo.Name = myname
			for k := range InfoChList {
				if InfoChList[k].Address == tmpinfo.Address {
					InfoChList[k].Name = myname
				}
			}
			for k := range InfoList {
				if InfoList[k].Address == tmpinfo.Address {
					InfoList[k].Name = myname
				}
			}
		case "listuser":
			var strlist []string
			for k := range InfoChList {
				strlist = append(strlist, InfoChList[k].Name)
			}
			res2B, _ := json.Marshal(strlist)
			infoChTmp.Ch <- tmpinfo.Name + ": " + string(res2B)
		case "listroom":
			var roomlist []string
			for k := range RoomList {
				roomlist = append(roomlist, RoomList[k].Name)
			}
			res2B, _ := json.Marshal(roomlist)
			infoChTmp.Ch <- tmpinfo.Name + ": " + string(res2B)
		case "createroom":
			if infoChTmp.Name == "" {
				infoChTmp.Ch <- tmpinfo.Address + ": " + "请先输入昵称"
				infoChTmp.Ch <- tmpinfo.Address + ": " + help()
				continue
			}
			infoChTmp.Ch <- infoChTmp.Address + ":输入要创建的房间号"
			var roomname string
			if input.Scan() {
				roomname = input.Text()
			}
			judge := JudgeStringSpecialSymbol(roomname)
			if judge == false {
				infoChTmp.Ch <- infoChTmp.Name + ":房间号只支持大小写A-z以及0-9,长度不超过20"
				continue
			}
			var sign bool = false
			for k := range RoomList {
				if RoomList[k].Name == roomname {
					infoChTmp.Ch <- infoChTmp.Name + ":已被使用,请重试"
					sign = true
					break
				}
			}
			if sign == true {
				continue
			}
			var tmpData ChatGroup
			tmpData.Name = roomname
			tmpData.ChList = append(tmpData.ChList, infoChTmp)
			RoomList = append(RoomList, tmpData)
			infoChTmp.Ch <- infoChTmp.Name + ":房间创建成功，可通过命令listroom查看"
		case "joinroom":
			if infoChTmp.Name == "" {
				infoChTmp.Ch <- tmpinfo.Address + ": " + "请先输入昵称"
				infoChTmp.Ch <- tmpinfo.Address + ": " + help()
				continue
			}
			infoChTmp.Ch <- infoChTmp.Name + ":输入要加入的房间号"
			var roomname string
			if input.Scan() {
				roomname = input.Text()
			}
			var sign bool = false
			for k := range RoomList {
				if RoomList[k].Name == roomname {
					for i := range RoomList[k].ChList {
						if RoomList[k].ChList[i].Name == infoChTmp.Name {
							infoChTmp.Ch <- infoChTmp.Name + ":你已经加入过该房间"
							sign = true
							break
						}
					}
					if sign == true {
						break
					}
					RoomList[k].ChList = append(RoomList[k].ChList, infoChTmp)
					infoChTmp.Ch <- infoChTmp.Name + ":房间加入成功"
				}
			}
			if sign == false {
				infoChTmp.Ch <- "房间不存在，可通过命令listroom查看"
			}
		case "help":
			infoChTmp.Ch <- tmpinfo.Name + ": " + help()
		default:
			if infoChTmp.Name == "" {
				infoChTmp.Ch <- tmpinfo.Address + ": " + "请先输入昵称"
				infoChTmp.Ch <- tmpinfo.Address + ": " + help()
				continue
			}
			if len(input.Text()) == 0 {
				messages <- tmpinfo.Name + ": " + input.Text()
				continue
			}
			if string(input.Text())[0] == '@' {
				strtmp := stringToDestinationAddr(input.Text())
				contenttmp := stringToDestinationContent(input.Text())
				var sign bool = false
				for k := range InfoChList {
					if strtmp == InfoChList[k].Name {
						InfoChList[k].Ch <- tmpinfo.Name + "悄悄对你说: " + contenttmp
						sign = true
						break
					}
				}
				if sign == false {
					infoChTmp.Ch <- "not found"
				}
			} else {
				messages <- tmpinfo.Name + ": " + input.Text()
			}
		}
	}
	// NOTE: ignoring potential errors from input.Err()
	for k := range InfoChList {
		if InfoChList[k].Address == infoChTmp.Address {
			InfoChList = append(InfoChList[:k], InfoChList[(k+1):]...)
			break
		}
	}
	for k := range InfoList {
		if InfoList[k].Address == infoChTmp.Address {
			InfoList = append(InfoList[:k], InfoList[(k+1):]...)
			break
		}
	}
	leaving <- ch
	messages <- tmpinfo.Address + " has left"
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

func help() string {
	return (`
	please choose options:
		- joinroom : 加入房间			格式:"joinroom" 
		- listroom : 获取所有房间号			格式:"listroom" 
		- listuser : 获取所有在线用户Name			格式:"listuser"
		- myname : 注册自己的聊天昵称			格式:"myname" 
		
        `)
}

//判断用户名合法性
func JudgeStringSpecialSymbol(input string) bool {
	f := func(r rune) bool {
		return (r < 'A' && r > '9') || r > 'z' || (r > 'Z' && r < 'a') || r < '0'
	}
	if strings.IndexFunc(input, f) != -1 {
		return false
	}
	if len(input) >= 20 {
		return false
	} else {
		return true
	}
}

func PrintListName() {
	var strlist []string
	for k := range InfoChList {
		strlist = append(strlist, InfoChList[k].Name)
		res2B, _ := json.Marshal(strlist)
		fmt.Println(string(res2B))
		time.Sleep(1 * time.Second)
	}
}
func PrintListAddress() {
	var strlist []string
	for k := range InfoChList {
		strlist = append(strlist, InfoChList[k].Address)
		res2B, _ := json.Marshal(strlist)
		fmt.Println(string(res2B))
		time.Sleep(1 * time.Second)
	}
}

func Printint() {
	for k := 0; k < 1000000; k++ {
		fmt.Println(k)
		time.Sleep(1 * time.Second)
	}
}
