package main

import (
	"bufio"
	Cf "chatserver/chfunc"
	Pf "chatserver/publicfunc"
	Pt "chatserver/publictype"
	"encoding/json"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("fuck")
	go broadcaster()
	// go Cf.Printint()
	go Cf.PrintListAddress()
	go Cf.PrintListName()
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
		go handleConn(tmpinfo)
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
func handleConn(tmpinfo Pt.ClientInfo) {
	ch := make(chan string) // outgoing client messages
	go Cf.ClientWriter(tmpinfo.ConnChan, ch)
	var infoChTmp Pt.ClientChInfo
	infoChTmp.Ch = ch
	infoChTmp.Address = tmpinfo.Address
	Pt.InfoChList = append(Pt.InfoChList, infoChTmp)
	ch <- "You are " + tmpinfo.Address
	Pt.Messages <- tmpinfo.Address + " has arrived"
	Pt.Entering <- ch

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
			judge := Pf.JudgeStringSpecialSymbol(myname)
			if judge == false {
				infoChTmp.Ch <- infoChTmp.Address + ":昵称只支持大小写A-z以及0-9,长度不超过20"
				//重新循环用户输入
				continue
			}
			var sign bool = false
			for k := range Pt.InfoChList {
				//判断是否有重名
				if Pt.InfoChList[k].Name == myname {
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
			for k := range Pt.InfoChList {
				if Pt.InfoChList[k].Address == tmpinfo.Address {
					Pt.InfoChList[k].Name = myname
				}
			}
			for k := range Pt.InfoList {
				if Pt.InfoList[k].Address == tmpinfo.Address {
					Pt.InfoList[k].Name = myname
				}
			}
		//列出所有用户昵称命令 没有昵称也能查询
		case "listuser":
			var strlist []string
			for k := range Pt.InfoChList {
				strlist = append(strlist, Pt.InfoChList[k].Name)
			}
			res2B, _ := json.Marshal(strlist)
			infoChTmp.Ch <- tmpinfo.Name + ": " + string(res2B)
		//列出所有组room命令 没有昵称也能查询
		case "listroom":
			var roomlist []string
			for k := range Pt.RoomList {
				roomlist = append(roomlist, Pt.RoomList[k].Name)
			}
			res2B, _ := json.Marshal(roomlist)
			infoChTmp.Ch <- tmpinfo.Name + ": " + string(res2B)
		//创建用户命令
		case "createroom":
			//判断是否有昵称 没有昵称不能操作
			if infoChTmp.Name == "" {
				infoChTmp.Ch <- tmpinfo.Address + ": " + "请先输入昵称"
				infoChTmp.Ch <- tmpinfo.Address + ": " + Pf.Helpstring()
				continue
			}
			infoChTmp.Ch <- infoChTmp.Address + ":输入要创建的房间号"
			var roomname string
			if input.Scan() {
				roomname = input.Text()
			}
			//合法性检查
			judge := Pf.JudgeStringSpecialSymbol(roomname)
			if judge == false {
				infoChTmp.Ch <- infoChTmp.Name + ":房间号只支持大小写A-z以及0-9,长度不超过20"
				continue
			}
			var sign bool = false
			//检查重名
			for k := range Pt.RoomList {
				if Pt.RoomList[k].Name == roomname {
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
			var tmpData Pt.ChatGroup
			tmpData.Name = roomname
			tmpData.ChList = append(tmpData.ChList, infoChTmp)
			Pt.RoomList = append(Pt.RoomList, tmpData)
			infoChTmp.Ch <- infoChTmp.Name + ":房间创建成功，可通过命令listroom查看"
		//加入房间命令
		case "joinroom":
			//判断是否有昵称 没有昵称不能操作
			if infoChTmp.Name == "" {
				infoChTmp.Ch <- tmpinfo.Address + ": " + "请先输入昵称"
				infoChTmp.Ch <- tmpinfo.Address + ": " + Pf.Helpstring()
				continue
			}
			infoChTmp.Ch <- infoChTmp.Name + ":输入要加入的房间号"
			var roomname string
			if input.Scan() {
				roomname = input.Text()
			}
			var sign bool = false
			//检查是否存在
			for k := range Pt.RoomList {
				if Pt.RoomList[k].Name == roomname {
					for i := range Pt.RoomList[k].ChList {
						if Pt.RoomList[k].ChList[i].Name == infoChTmp.Name {
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
					Pt.RoomList[k].ChList = append(Pt.RoomList[k].ChList, infoChTmp)
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
			infoChTmp.Ch <- tmpinfo.Name + ": " + Pf.Helpstring()
		default:
			//先检查有没有昵称
			if infoChTmp.Name == "" {
				infoChTmp.Ch <- tmpinfo.Address + ": " + "请先输入昵称"
				infoChTmp.Ch <- tmpinfo.Address + ": " + Pf.Helpstring()
				//重来 判断
				continue
			}
			//如果输入为空
			if len(input.Text()) == 0 {
				Pt.Messages <- tmpinfo.Name + ": " + input.Text()
				//重来 判断
				continue
			}
			//私聊1v1
			if string(input.Text())[0] == '@' {
				//截取输入
				strtmp := Pf.StringToDestinationAddr(input.Text())
				contenttmp := Pf.StringToDestinationContent(input.Text())
				var sign bool = false
				//在公共管道数组里找目标管道
				for k := range Pt.InfoChList {
					if strtmp == Pt.InfoChList[k].Name {
						Pt.InfoChList[k].Ch <- tmpinfo.Name + "悄悄对你说: " + contenttmp
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
				strtmp := Pf.StringToDestinationAddr(input.Text())
				contenttmp := Pf.StringToDestinationContent(input.Text())
				var sign bool = false
				//在room数组中找目标管道
				for k := range Pt.RoomList {
					if strtmp == Pt.RoomList[k].Name {
						for i := range Pt.RoomList[k].ChList {
							Pt.RoomList[k].ChList[i].Ch <- tmpinfo.Name + "在房间" + strtmp + "小声说: " + contenttmp
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
				Pt.Messages <- tmpinfo.Name + ": " + input.Text()
			}
		}
	}
	// NOTE: ignoring potential errors from input.Err()
	//公共管道数组中删除断开的连接
	for k := range Pt.InfoChList {
		if Pt.InfoChList[k].Address == infoChTmp.Address {
			Pt.InfoChList = append(Pt.InfoChList[:k], Pt.InfoChList[(k+1):]...)
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
