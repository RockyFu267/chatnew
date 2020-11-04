package main

import (
	"bufio"
	Cf "chatserver/chfunc"
	Pf "chatserver/publicfunc"
	Pt "chatserver/publictype"
	Uc "chatserver/usercmd"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:18000")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("fuck")
	go broadcaster()
	// go Cf.Printint()
	// go Cf.PrintListAddress()
	// go Cf.PrintListName()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		who := conn.RemoteAddr().String()
		fmt.Println(who, "已经建立连接")
		var tmpinfo Pt.ClientInfo
		tmpinfo.ConnChan = conn
		tmpinfo.Address = who
		Pt.InfoList = append(Pt.InfoList, tmpinfo)
		go handleConn(&tmpinfo)
		// go Less5SecondEchoEachOther(tmpinfo) //俩启动间隔小于五秒互相通信
		//go BecomeUpper(conn)			//变大写
		// go ReturnTime(conn)          //输出时间
	}
}

//broadcaster 事件定义
func broadcaster() {
	clients := make(map[Pt.ClientChan]bool) // all connected clients
	for {
		select {
		case msg := <-Pt.Messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-Pt.Entering:
			clients[cli] = true

		case cli := <-Pt.Leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(tmpinfo *Pt.ClientInfo) {
	ch := make(chan string) // outgoing client messages
	go Cf.ClientWriter(tmpinfo.ConnChan, ch)
	var infoChTmp Pt.ClientChInfo
	infoChTmp.Ch = ch
	infoChTmp.Address = tmpinfo.Address
	Pt.InfoChList = append(Pt.InfoChList, infoChTmp)
	Pt.InfoPubChList = append(Pt.InfoPubChList, infoChTmp)
	ch <- "You are " + tmpinfo.Address
	Pt.Messages <- tmpinfo.Address + " has arrived"
	Pt.Entering <- ch

	input := bufio.NewScanner(tmpinfo.ConnChan)
	//循环用户输入
	for input.Scan() {
		switch input.Text() {
		//昵称命令
		case "myname":
			Uc.MyName(&infoChTmp, tmpinfo, input)
		//列出所有用户昵称命令 没有昵称也能查询
		case "listuser":
			Uc.Listuser(infoChTmp)
		//列出所有组room命令 没有昵称也能查询
		case "listroom":
			Uc.Listroom(infoChTmp)
		//创建用户命令
		case "createroom":
			Uc.Createroom(infoChTmp, tmpinfo.Address, input)
		//加入房间命令
		case "joinroom":
			Uc.Joinroom(infoChTmp, tmpinfo.Address, input)
		case "help":
			infoChTmp.Ch <- "命令提示: " + Pf.Helpstring()
		default:
			Uc.DefaultCmd(infoChTmp, tmpinfo.Address, input)
		}
	}
	// NOTE: ignoring potential errors from input.Err()
	//总管道数组中删除断开的连接
	for k := range Pt.InfoChList {
		if Pt.InfoChList[k].Address == infoChTmp.Address {
			Pt.InfoChList = append(Pt.InfoChList[:k], Pt.InfoChList[(k+1):]...)
			break
		}
	}
	//公共管道数组中删除断开的连接
	for k := range Pt.InfoPubChList {
		if Pt.InfoPubChList[k].Address == infoChTmp.Address {
			Pt.InfoPubChList = append(Pt.InfoPubChList[:k], Pt.InfoPubChList[(k+1):]...)
			break
		}
	}
	//TCP连接数组删除断开的连接
	for k := range Pt.InfoList {
		if Pt.InfoList[k].Address == infoChTmp.Address {
			Pt.InfoList = append(Pt.InfoList[:k], Pt.InfoList[(k+1):]...)
			break
		}
	}
	//房间room数组中删除断开的连接
	for k := range Pt.RoomList {
		for i := range Pt.RoomList[k].ChList {
			if Pt.RoomList[k].ChList[i].Address == infoChTmp.Address {
				Pt.RoomList[k].ChList = append(Pt.RoomList[k].ChList[:i], Pt.RoomList[k].ChList[(i+1):]...)
				break
			}
		}
	}
	Pt.Leaving <- ch
	Pt.Messages <- infoChTmp.Name + " has left"
	tmpinfo.ConnChan.Close()
}
