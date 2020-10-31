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
		infoChTmp.Ch <- infoChTmp.Address + ":昵称只支持大小写A-z以及0-9,长度不超过20,小于2"
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
		infoChTmp.Ch <- infoChTmp.Name + ":房间号只支持大小写A-z以及0-9,长度不超过20，小于2"
		return
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
		return
	}
	//正常赋值
	var tmpData Pt.ChatGroup
	tmpData.Name = roomname
	tmpData.ChList = append(tmpData.ChList, infoChTmp)
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
		Pt.Messages <- infoChTmp.Name + ": " + input.Text()
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
		var sign bool = false
		//在公共管道数组里找目标管道
		for k := range Pt.InfoChList {
			if strtmp == Pt.InfoChList[k].Name {
				Pt.InfoChList[k].Ch <- infoChTmp.Name + "悄悄对你说: " + contenttmp
				sign = true
				break
			}
		}
		//状态未变 找不到目标管道
		if sign == false {
			infoChTmp.Ch <- "user not found"
		}
		//重来 判断
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
			Pt.InfoPubChList[k].Ch <- infoChTmp.Name + ": " + input.Text()
		}
	}
}
