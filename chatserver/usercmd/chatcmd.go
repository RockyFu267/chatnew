package usercmd

import (
	"bufio"
	Pf "chatserver/publicfunc"
	Pt "chatserver/publictype"
	"encoding/json"
)

//MyName 昵称命令 有回写需要指针
func MyName(infoChTmp *Pt.ClientChInfo, tmpinfo *Pt.ClientInfo, input *bufio.Scanner) {
	//判断是否已经输入过昵称
	if infoChTmp.Name != "" {
		infoChTmp.Ch <- "你已经输入过昵称：" + infoChTmp.Name
		//重新循环用户输入
		return
	}
	infoChTmp.Ch <- infoChTmp.Address + ":输入昵称"
	var myname string
	if input.Scan() {
		myname = input.Text()
	}
	//检查合法性
	judge := Pf.JudgeStringSpecialSymbol(myname)
	if judge == false {
		infoChTmp.Ch <- infoChTmp.Address + ":昵称只支持大小写A-z以及0-9,长度不超过20,不小于2"
		//重新循环用户输入
		return
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
		return
	}
	//一切正常 赋值
	infoChTmp.Name = myname
	tmpinfo.Name = myname
	for k := range Pt.InfoChList {
		if Pt.InfoChList[k].Address == tmpinfo.Address {
			Pt.InfoChList[k].Name = myname
		}
	}
	for k := range Pt.InfoPubChList {
		if Pt.InfoPubChList[k].Address == tmpinfo.Address {
			Pt.InfoPubChList[k].Name = myname
		}
	}
	for k := range Pt.InfoList {
		if Pt.InfoList[k].Address == tmpinfo.Address {
			Pt.InfoList[k].Name = myname
		}
	}
}

//Listuser 列出所有用户昵称命令 没有昵称也能查询
func Listuser(infoChTmp Pt.ClientChInfo) {
	var strlist []string
	for k := range Pt.InfoChList {
		strlist = append(strlist, Pt.InfoChList[k].Name)
	}
	res2B, _ := json.Marshal(strlist)
	infoChTmp.Ch <- "目前在线的用户: " + string(res2B)
}

//Listroom 列出所有组room命令 没有昵称也能查询
func Listroom(infoChTmp Pt.ClientChInfo) {
	var roomlist []string
	for k := range Pt.RoomList {
		roomlist = append(roomlist, Pt.RoomList[k].Name)
	}
	res2B, _ := json.Marshal(roomlist)
	infoChTmp.Ch <- "房间列表: " + string(res2B)
}

//Createroom 创建用户命令
func Createroom(infoChTmp Pt.ClientChInfo, address string, input *bufio.Scanner) {
	//判断是否有昵称 没有昵称不能操作
	if infoChTmp.Name == "" {
		infoChTmp.Ch <- address + ": " + "请先输入昵称"
		infoChTmp.Ch <- address + ": " + Pf.Helpstring()
		return
	}
	infoChTmp.Ch <- infoChTmp.Address + ":输入要创建的房间号"
	var roomname string
	if input.Scan() {
		roomname = input.Text()
	}
	//合法性检查
	judge := Pf.JudgeStringSpecialSymbol(roomname)
	if judge == false {
		infoChTmp.Ch <- infoChTmp.Name + ":房间号只支持大小写A-z以及0-9,长度不超过20,不小于2"
		return
	}
	//检查重名
	for k := range Pt.RoomList {
		if Pt.RoomList[k].Name == roomname {
			infoChTmp.Ch <- infoChTmp.Name + ":已被使用,请重试"
			return
		}
	}
	infoChTmp.Ch <- infoChTmp.Address + ":输入要创建的房间的AccessKey"
	var ack string
	if input.Scan() {
		ack = input.Text()
	}
	judgeack := Pf.JudgeStringSpecialSymbol(ack)
	if judgeack == false {
		infoChTmp.Ch <- infoChTmp.Name + ":AccessKey只支持大小写A-z以及0-9,长度不超过20,不小于2"
		return
	}
	//正常赋值
	infoChTmp.RoomLeader = true
	var tmpData Pt.ChatGroup
	tmpData.Name = roomname
	tmpData.ChList = append(tmpData.ChList, infoChTmp)
	tmpData.AccessKey = ack
	Pt.RoomList = append(Pt.RoomList, tmpData)
	infoChTmp.Ch <- infoChTmp.Name + ":房间创建成功，可通过命令listroom查看"
}

//Joinroom 加入房间命令
func Joinroom(infoChTmp Pt.ClientChInfo, address string, input *bufio.Scanner) {
	//判断是否有昵称 没有昵称不能操作
	if infoChTmp.Name == "" {
		infoChTmp.Ch <- address + ": " + "请先输入昵称"
		infoChTmp.Ch <- address + ": " + Pf.Helpstring()
		return
	}
	infoChTmp.Ch <- infoChTmp.Name + ":输入要加入的房间号"
	var roomname string
	if input.Scan() {
		roomname = input.Text()
	}
	if roomname == "public" {
		for k := range Pt.InfoPubChList {
			if Pt.InfoPubChList[k].Address == infoChTmp.Address {
				infoChTmp.Ch <- infoChTmp.Name + ":你已经加入过该房间"
				return
			}
		}
		Pt.InfoPubChList = append(Pt.InfoPubChList, infoChTmp)
		infoChTmp.Ch <- infoChTmp.Name + ":已加入公共聊天室"
		return
	}
	//检查是否存在
	for k := range Pt.RoomList {
		if Pt.RoomList[k].Name == roomname {
			for i := range Pt.RoomList[k].ChList {
				if Pt.RoomList[k].ChList[i].Name == infoChTmp.Name {
					infoChTmp.Ch <- infoChTmp.Name + ":你已经加入过该房间"
					return
				}
			}
			//房间人数上限
			if len(Pt.RoomList[k].ChList) >= 10 {
				infoChTmp.Ch <- infoChTmp.Name + ":房间人数已达上限"
				return
			}
			//检查accesskey的输入
			infoChTmp.Ch <- infoChTmp.Name + ":输入要加入的房间的AccessKey"
			var ack string
			if input.Scan() {
				ack = input.Text()
			}
			if ack != Pt.RoomList[k].AccessKey {
				infoChTmp.Ch <- infoChTmp.Name + ":AccessKey错误，加入失败"
				return
			}
			//正常赋值
			Pt.RoomList[k].ChList = append(Pt.RoomList[k].ChList, infoChTmp)
			infoChTmp.Ch <- infoChTmp.Name + ":房间加入成功"
			return
		}
	}
	//循环中没有房间name
	infoChTmp.Ch <- "房间不存在，可通过命令listroom查看"

}

//DefaultCmd 加入房间命令
func DefaultCmd(infoChTmp Pt.ClientChInfo, address string, input *bufio.Scanner) {
	//先检查有没有昵称
	if infoChTmp.Name == "" {
		infoChTmp.Ch <- address + ": " + "请先输入昵称"
		infoChTmp.Ch <- address + ": " + Pf.Helpstring()
		//重来 判断
		return
	}
	//如果输入为空
	if len(input.Text()) == 0 {
		//Pt.Messages <- infoChTmp.Name + ": " + input.Text()
		//重来 判断
		return
	}
	//退出组room
	if string(input.Text())[0] == '!' {
		//截取输入
		strtmp := Pf.StringToDestinationName(input.Text())
		if strtmp == "public" {
			//在公共管道数组里找目标管道
			for k := range Pt.InfoPubChList {
				if infoChTmp.Name == Pt.InfoPubChList[k].Name {
					for k := range Pt.InfoPubChList {
						if Pt.InfoPubChList[k].Address == infoChTmp.Address {
							Pt.InfoPubChList = append(Pt.InfoPubChList[:k], Pt.InfoPubChList[(k+1):]...)
							infoChTmp.Ch <- infoChTmp.Name + ":成功退出公共聊天"
							return
						}
					}
					infoChTmp.Ch <- infoChTmp.Name + ":本来就不在该房间"
					return
				}
			}
			infoChTmp.Ch <- "user not found"
			//结束public的命令判断
			return
		}
		//在room数组中找目标管道
		for k := range Pt.RoomList {
			//判断房间是否存在
			if strtmp == Pt.RoomList[k].Name {
				//在房间管道中寻找目标名
				for i := range Pt.RoomList[k].ChList {
					//判断管道是否在房间中
					if Pt.RoomList[k].ChList[i].Address == infoChTmp.Address {
						Pt.RoomList[k].ChList = append(Pt.RoomList[k].ChList[:i], Pt.RoomList[k].ChList[(i+1):]...)
						infoChTmp.Ch <- infoChTmp.Name + ":成功退出房间" + strtmp
						return
					}
				}
				infoChTmp.Ch <- infoChTmp.Name + ":本来就不在该房间"
				return
			}
		}
		infoChTmp.Ch <- "房间不存在"
		//结束public的命令判断
		return
	}
	//私聊1v1
	if string(input.Text())[0] == '@' {
		//截取输入
		strtmp := Pf.StringToDestinationAddr(input.Text())
		contenttmp := Pf.StringToDestinationContent(input.Text())
		//检查是否在自己的好友列表里
		if _, ok := infoChTmp.Friends[strtmp]; ok {
			//存在
			//在公共管道数组里找目标管道
			for k := range Pt.InfoChList {
				if strtmp == Pt.InfoChList[k].Name {
					//检查是否在对方的好友列表里
					if _, ok := Pt.InfoChList[k].Friends[infoChTmp.Name]; ok {
						//直接发送消息
						Pt.InfoChList[k].Ch <- infoChTmp.Name + "悄悄对你说: " + contenttmp
						return
					}
					//不在对方的列表里
					infoChTmp.Ch <- "你与" + strtmp + "不是好友关系，请先成为好友"
					return
				}
			}
			infoChTmp.Ch <- "对方已下线"
			return
		}
		//不在自己的列表里
		infoChTmp.Ch <- "你与" + strtmp + "不是好友关系，请先成为好友"
		return
	}
	//小房间私聊
	if string(input.Text())[0] == '#' {
		//截取
		strtmp := Pf.StringToDestinationAddr(input.Text())
		contenttmp := Pf.StringToDestinationContent(input.Text())
		//检查是否已加入目标管道
		for k := range Pt.RoomList {
			if strtmp == Pt.RoomList[k].Name {
				for i := range Pt.RoomList[k].ChList {
					if Pt.RoomList[k].ChList[i].Name == infoChTmp.Name {
						for j := range Pt.RoomList[k].ChList {
							Pt.RoomList[k].ChList[j].Ch <- infoChTmp.Name + "在房间" + strtmp + "小声说: " + contenttmp
						}
						return
					}
				}
				infoChTmp.Ch <- "请先加入房间" + strtmp
				return
			}
		}
		infoChTmp.Ch <- "room not found"
		//最后一个不需要跳出重来判断
	} else {
		//先判断是不是正在公共组里
		var sign bool
		for k := range Pt.InfoPubChList {
			if Pt.InfoPubChList[k].Name == infoChTmp.Name {
				sign = true
				break
			}
		}
		if sign == false {
			infoChTmp.Ch <- "请先加入public房间"
			return
		}
		//遍历公共管道数组 公共广播
		for k := range Pt.InfoPubChList {
			if Pt.InfoPubChList[k].Name != "" {
				Pt.InfoPubChList[k].Ch <- infoChTmp.Name + "在公共房间广播说: " + input.Text()
			}
		}
	}
}

//AddFriends 添加好友 需要指针会写
func AddFriends(infoChTmp *Pt.ClientChInfo, address string, input *bufio.Scanner) {
	//判断是否有昵称 没有昵称不能操作
	if infoChTmp.Name == "" {
		infoChTmp.Ch <- address + ": " + "请先输入昵称"
		infoChTmp.Ch <- address + ": " + Pf.Helpstring()
		return
	}
	infoChTmp.Ch <- "请输入要添加为好友的昵称"
	var firendName string
	if input.Scan() {
		firendName = input.Text()
	}
	if firendName == infoChTmp.Name {
		infoChTmp.Ch <- "不能添加自己"
		return
	}
	//先检查自身
	if _, ok := infoChTmp.Friends[firendName]; ok {
		//存在
		infoChTmp.Ch <- infoChTmp.Name + ":你已经加过该好友"
		return
	}
	for k := range Pt.InfoChList {
		if firendName == Pt.InfoChList[k].Name {
			//对方是否已经添加过
			if _, ok := Pt.InfoChList[k].Friends[infoChTmp.Name]; ok {
				//对方有，直接添加好友
				infoChTmp.Friends[firendName] = true
				infoChTmp.Ch <- infoChTmp.Name + ":" + firendName + "与你成为好友"
				Pt.InfoChList[k].Ch <- infoChTmp.Name + "已同意添加你为好友"
				return
			}
			//对方没有
			infoChTmp.Ch <- "请输入验证消息:"
			var addContent string
			if input.Scan() {
				addContent = input.Text()
			}
			infoChTmp.Friends[firendName] = true
			Pt.InfoChList[k].Ch <- infoChTmp.Name + "向你发送添加好友请求,并附言:" + addContent
			infoChTmp.Ch <- infoChTmp.Name + ":已向" + firendName + "发送添加好友请求,等待对方确认"
			return
		}
	}
	infoChTmp.Ch <- "user not found"
}

//DeleteFriends 添加好友 需要指针会写
func DeleteFriends(infoChTmp *Pt.ClientChInfo, address string, input *bufio.Scanner) {
	//判断是否有昵称 没有昵称不能操作
	if infoChTmp.Name == "" {
		infoChTmp.Ch <- address + ": " + "请先输入昵称"
		infoChTmp.Ch <- address + ": " + Pf.Helpstring()
		return
	}
	infoChTmp.Ch <- "请输入要删除为好友的昵称"
	var firendName string
	if input.Scan() {
		firendName = input.Text()
	}
	if firendName == infoChTmp.Name {
		infoChTmp.Ch <- "不能删除自己"
		return
	}
	//检查自身
	if _, ok := infoChTmp.Friends[firendName]; ok {
		for k := range Pt.InfoChList {
			if firendName == Pt.InfoChList[k].Name {
				//对方是否已经添加过
				if _, ok := Pt.InfoChList[k].Friends[infoChTmp.Name]; ok {
					delete(infoChTmp.Friends, firendName)
					delete(Pt.InfoChList[k].Friends, infoChTmp.Name)
					infoChTmp.Ch <- infoChTmp.Name + ":你已经删除该好友"
					Pt.InfoChList[k].Ch <- infoChTmp.Name + "已将你从好友列表中移除"
					return
				}
				infoChTmp.Ch <- infoChTmp.Name + ":对方还未回应你之前的添加好友请求，你的好友请求已收回"
				delete(infoChTmp.Friends, firendName)
				return
			}
		}
		//对方已离线
		delete(infoChTmp.Friends, firendName)
		infoChTmp.Ch <- infoChTmp.Name + ":你已经删除该好友"
		return
	}
	infoChTmp.Ch <- "user not found in your friednslist"
	return
}
