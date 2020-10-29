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

//ClientChan 定义双向管道
type ClientChan chan string

//AllClients 初始化所有的连接map
var AllClients = make(map[string]net.Conn)

//ClientInfo 定义tcp连接的结构体
type ClientInfo struct {
	ConnChan net.Conn
	Address  string
	Name     string
}

//ClientChInfo 定义管道的结构体
type ClientChInfo struct {
	Ch      ClientChan
	Address string
	Name    string
}

//ChatGroup 定义组room的结构体
type ChatGroup struct {
	Name   string
	ChList []ClientChInfo
}

//InfoList 初始化tcp连接的数组 后期可以优化改map
var InfoList []ClientInfo

//InfoChList 初始化管道的数组 后期可以优化改map
var InfoChList []ClientChInfo

//RoomList 初始化组room的数组
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

//后期优化可以加其他的事件管道
var (
	entering = make(chan ClientChan)
	leaving  = make(chan ClientChan)
	messages = make(chan string) // all incoming client messages
)

//broadcaster 事件定义
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
	//循环用户输入
	for input.Scan() {
		switch input.Text() {
		//昵称命令
		case "myname":
			//判断是否已经输入过昵称
			if infoChTmp.Name != "" {
				infoChTmp.Ch <- "你已经输入过昵称：" + tmpinfo.Name
				//重新循环用户输入
				continue
			}
			infoChTmp.Ch <- infoChTmp.Address + ":输入昵称"
			var myname string
			if input.Scan() {
				myname = input.Text()
			}
			//检查合法性
			judge := JudgeStringSpecialSymbol(myname)
			if judge == false {
				infoChTmp.Ch <- infoChTmp.Address + ":昵称只支持大小写A-z以及0-9,长度不超过20"
				//重新循环用户输入
				continue
			}
			var sign bool = false
			for k := range InfoChList {
				//判断是否有重名
				if InfoChList[k].Name == myname {
					infoChTmp.Ch <- infoChTmp.Address + ":已被使用,请重新输入昵称"
					//标记有重名
					sign = true
					//跳出此次循环
					break
				}
			}
			//出现重名重新循环用户输入
			if sign == true {
				continue
			}
			//一切正常 赋值
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
		//列出所有用户昵称命令 没有昵称也能查询
		case "listuser":
			var strlist []string
			for k := range InfoChList {
				strlist = append(strlist, InfoChList[k].Name)
			}
			res2B, _ := json.Marshal(strlist)
			infoChTmp.Ch <- tmpinfo.Name + ": " + string(res2B)
		//列出所有组room命令 没有昵称也能查询
		case "listroom":
			var roomlist []string
			for k := range RoomList {
				roomlist = append(roomlist, RoomList[k].Name)
			}
			res2B, _ := json.Marshal(roomlist)
			infoChTmp.Ch <- tmpinfo.Name + ": " + string(res2B)
		//创建用户命令
		case "createroom":
			//判断是否有昵称 没有昵称不能操作
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
			//合法性检查
			judge := JudgeStringSpecialSymbol(roomname)
			if judge == false {
				infoChTmp.Ch <- infoChTmp.Name + ":房间号只支持大小写A-z以及0-9,长度不超过20"
				continue
			}
			var sign bool = false
			//检查重名
			for k := range RoomList {
				if RoomList[k].Name == roomname {
					infoChTmp.Ch <- infoChTmp.Name + ":已被使用,请重试"
					sign = true
					//跳出此次循环
					break
				}
			}
			//出现重名重新循环用户输入
			if sign == true {
				continue
			}
			//正常赋值
			var tmpData ChatGroup
			tmpData.Name = roomname
			tmpData.ChList = append(tmpData.ChList, infoChTmp)
			RoomList = append(RoomList, tmpData)
			infoChTmp.Ch <- infoChTmp.Name + ":房间创建成功，可通过命令listroom查看"
		//加入房间命令
		case "joinroom":
			//判断是否有昵称 没有昵称不能操作
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
			//检查是否存在
			for k := range RoomList {
				if RoomList[k].Name == roomname {
					for i := range RoomList[k].ChList {
						if RoomList[k].ChList[i].Name == infoChTmp.Name {
							infoChTmp.Ch <- infoChTmp.Name + ":你已经加入过该房间"
							sign = true
							//跳出房间成员的数组循环
							break
						}
					}
					if sign == true {
						//已经在该房间 跳出检查的循环
						break
					}
					//正常赋值
					RoomList[k].ChList = append(RoomList[k].ChList, infoChTmp)
					infoChTmp.Ch <- infoChTmp.Name + ":房间加入成功"
					sign = true
					break
				}
			}
			//标记状态未变 不存在该room
			if sign == false {
				infoChTmp.Ch <- "房间不存在，可通过命令listroom查看"
			}
		case "help":
			infoChTmp.Ch <- tmpinfo.Name + ": " + help()
		default:
			//先检查有没有昵称
			if infoChTmp.Name == "" {
				infoChTmp.Ch <- tmpinfo.Address + ": " + "请先输入昵称"
				infoChTmp.Ch <- tmpinfo.Address + ": " + help()
				//重来 判断
				continue
			}
			//如果输入为空
			if len(input.Text()) == 0 {
				messages <- tmpinfo.Name + ": " + input.Text()
				//重来 判断
				continue
			}
			//私聊1v1
			if string(input.Text())[0] == '@' {
				//截取输入
				strtmp := stringToDestinationAddr(input.Text())
				contenttmp := stringToDestinationContent(input.Text())
				var sign bool = false
				//在公共管道数组里找目标管道
				for k := range InfoChList {
					if strtmp == InfoChList[k].Name {
						InfoChList[k].Ch <- tmpinfo.Name + "悄悄对你说: " + contenttmp
						sign = true
						break
					}
				}
				//状态未变 找不到目标管道
				if sign == false {
					infoChTmp.Ch <- "user not found"
				}
				//重来 判断
				continue
			}
			//小房间私聊
			if string(input.Text())[0] == '#' {
				//截取
				strtmp := stringToDestinationAddr(input.Text())
				contenttmp := stringToDestinationContent(input.Text())
				var sign bool = false
				//在room数组中找目标管道
				for k := range RoomList {
					if strtmp == RoomList[k].Name {
						for i := range RoomList[k].ChList {
							RoomList[k].ChList[i].Ch <- tmpinfo.Name + "在房间" + strtmp + "小声说: " + contenttmp
						}
						sign = true
						//跳出查找循环
						break
					}
				}
				//状态未变 没找到
				if sign == false {
					infoChTmp.Ch <- "room rnot found"
				}
				//最后一个不需要跳出重来判断
			} else {
				messages <- tmpinfo.Name + ": " + input.Text()
			}
		}
	}
	// NOTE: ignoring potential errors from input.Err()
	//公共管道数组中删除断开的连接
	for k := range InfoChList {
		if InfoChList[k].Address == infoChTmp.Address {
			InfoChList = append(InfoChList[:k], InfoChList[(k+1):]...)
			break
		}
	}
	//TCP连接数组删除断开的连接
	for k := range InfoList {
		if InfoList[k].Address == infoChTmp.Address {
			InfoList = append(InfoList[:k], InfoList[(k+1):]...)
			break
		}
	}
	//房间room数组中删除断开的连接
	for k := range RoomList {
		for i := range RoomList[k].ChList {
			if RoomList[k].ChList[i].Address == infoChTmp.Address {
				RoomList[k].ChList = append(RoomList[k].ChList[:i], RoomList[k].ChList[(i+1):]...)
				break
			}
		}
	}
	leaving <- ch
	messages <- infoChTmp.Name + " has left"
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

//JudgeStringSpecialSymbol 判断用户名合法性
func JudgeStringSpecialSymbol(input string) bool {
	f := func(r rune) bool {
		return (r < 'A' && r > '9') || r > 'z' || (r > 'Z' && r < 'a') || r < '0'
	}
	if strings.IndexFunc(input, f) != -1 {
		return false
	}
	if len(input) >= 20 {
		return false
	}
	return true

}

//Less5SecondEchoEachOther 俩启动间隔小于五秒互相通信
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
	for k := range InfoChList {
		strlist = append(strlist, InfoChList[k].Name)
		res2B, _ := json.Marshal(strlist)
		fmt.Println(string(res2B))
		time.Sleep(1 * time.Second)
	}
}

//PrintListAddress debug时候用的
func PrintListAddress() {
	var strlist []string
	for k := range InfoChList {
		strlist = append(strlist, InfoChList[k].Address)
		res2B, _ := json.Marshal(strlist)
		fmt.Println(string(res2B))
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
