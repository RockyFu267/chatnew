package main

import (
	"bufio"
	Cf "chatserver/chfunc"
	Pf "chatserver/publicfunc"
	Pt "chatserver/publictype"
	Uc "chatserver/usercmd"
	"fmt"
	"net"
	"time"
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
	infoChTmp.Friends = make(map[string]bool)
	//限制聊天室人数上限
	if len(Pt.InfoChList) >= 20 {
		fmt.Println(infoChTmp.Address + ":" + infoChTmp.Name + "连接达到上限")
		infoChTmp.Ch <- "连接达到上限"
		time.Sleep(1 * time.Second)
		tmpinfo.ConnChan.Close()
		return
	}

	Pt.InfoChList = append(Pt.InfoChList, infoChTmp)
	Pt.InfoPubChList = append(Pt.InfoPubChList, infoChTmp)
	ch <- "You are " + tmpinfo.Address
	Pt.Messages <- tmpinfo.Address + " has arrived"
	Pt.Entering <- ch

	input := bufio.NewScanner(tmpinfo.ConnChan)
	chScanTurnBool := make(chan bool)
	chScanCloseBool := make(chan bool)
	go func() {
		//循环用户输入
		for input.Scan() {
			chScanTurnBool <- true
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
			case "addfriends":
				Uc.AddFriends(&infoChTmp, tmpinfo.Address, input)
			case "delfriends":
				Uc.DeleteFriends(&infoChTmp, tmpinfo.Address, input)
			case "help":
				infoChTmp.Ch <- "命令提示: " + Pf.Helpstring()
			//游戏
			case "createcycles":
				Uc.CreateCycles(infoChTmp, tmpinfo.Address, input)
			case "joincycles":
				Uc.JoinCycles(infoChTmp, tmpinfo.Address, input)
			default:
				Uc.DefaultCmd(infoChTmp, tmpinfo.Address, input)
			}
		}
		chScanCloseBool <- true
	}()
	for {
		select {
		case <-chScanTurnBool:
		case <-chScanCloseBool:
			Pt.Messages <- infoChTmp.Address + ":" + infoChTmp.Name + "has left"
			fmt.Println(infoChTmp.Address + ":" + infoChTmp.Name + "主动断开连接")
			Cf.DeleteConn(infoChTmp, ch, tmpinfo)
			return
		case <-time.After(time.Duration(180 * time.Second)):
			Pt.Messages <- infoChTmp.Address + ":" + infoChTmp.Name + "timeout has left"
			fmt.Println(infoChTmp.Address + ":" + infoChTmp.Name + "超时断开连接")
			Cf.DeleteConn(infoChTmp, ch, tmpinfo)
			return
		}
	}

}
