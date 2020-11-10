package usercmd

import (
	"bufio"
	Pf "chatserver/publicfunc"
	Pt "chatserver/publictype"
	"fmt"
)

//CreateCycles 创建石头剪刀布的房间 1v1
func CreateCycles(infoChTmp Pt.ClientChInfo, address string, input *bufio.Scanner) {
	//判断是否有昵称 没有昵称不能操作
	if infoChTmp.Name == "" {
		infoChTmp.Ch <- address + ": " + "请先输入昵称"
		infoChTmp.Ch <- address + ": " + Pf.Helpstring()
		return
	}
	infoChTmp.Ch <- infoChTmp.Address + ":输入要创建的房间号"
	var gamename string
	if input.Scan() {
		gamename = input.Text()
	}
	//合法性检查
	judge := Pf.JudgeStringSpecialSymbol(gamename)
	if judge == false {
		infoChTmp.Ch <- infoChTmp.Name + ":房间号只支持大小写A-z以及0-9,长度不超过20,不小于2"
		return
	}
	//检查重名
	if _, ok := Pt.GameCyclesRoom[gamename]; ok {
		infoChTmp.Ch <- infoChTmp.Name + ":已存在,请重试"
		return
	}
	infoChTmp.Ch <- infoChTmp.Name + ":输入要创建的游戏房间的AccessKey"
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
	var tmpData Pt.InfoChListStruct
	tmpData.ChList = append(tmpData.ChList, infoChTmp)
	tmpData.Ack = ack
	tmpData.JoinStatus = true
	//数值初始化--join时候需要，目前是1v1 后者加入游戏即开始,多人应加上房主标签，其他人ready状态全true，房主可以start
	Pt.GameCyclesRoom[gamename] = tmpData
	// infoChTmp.Ch <- infoChTmp.Name + ":房间创建成功，可通过命令listroom查看"
	infoChTmp.Ch <- infoChTmp.Name + ":房间创建成功，等待对手"
	for input.Scan() {
		switch input.Text() {
		// //准备
		// case "ready":
		// 	//debug
		// 	fmt.Println("ready")
		case "exit":
			//debug
			fmt.Println("exit")
			//-----------------
			//如果房间只有自己 那就退出并删除房间
			if len(Pt.GameCyclesRoom[gamename].ChList) == 1 {
				delete(Pt.GameCyclesRoom, gamename)
				infoChTmp.Ch <- infoChTmp.Name + ":退出房间，房间被注销"
				return
			}
			//否则就移交给下一个(2号)元素位的用户
			var tmpData Pt.InfoChListStruct
			tmpData = Pt.GameCyclesRoom[gamename]
			tmpData.ChList = tmpData.ChList[1:]
			tmpData.ChList[0].RoomLeader = true
			Pt.GameCyclesRoom[gamename] = tmpData
			infoChTmp.Ch <- infoChTmp.Name + ":退出房间，房间被移交给其他用户"
			//debug
			fmt.Println(tmpData.ChList[0].RoomLeader)
			//-----------------
			return
		default:
			//debug
			fmt.Println("default")
			for k := range Pt.GameCyclesRoom[gamename].ChList {
				Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + ":" + input.Text()
			}
		}
	}
	return
}

//JoinCycles 加入房间命令
func JoinCycles(infoChTmp Pt.ClientChInfo, address string, input *bufio.Scanner) {
	//判断是否有昵称 没有昵称不能操作
	if infoChTmp.Name == "" {
		infoChTmp.Ch <- address + ": " + "请先输入昵称"
		infoChTmp.Ch <- address + ": " + Pf.Helpstring()
		return
	}
	infoChTmp.Ch <- infoChTmp.Name + ":输入要加入的房间号"
	var gamename string
	if input.Scan() {
		gamename = input.Text()
	}
	//检查是否存在
	if _, ok := Pt.GameCyclesRoom[gamename]; ok {
		for k := range Pt.GameCyclesRoom[gamename].ChList {
			if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
				infoChTmp.Ch <- infoChTmp.Name + ":你已经加入过该房间"
				return
			}
		}
		//房间人数上限或游戏中有人退出(德扑的房间存在被淘汰的玩家离场的情况)
		if len(Pt.GameCyclesRoom[gamename].ChList) >= 2 || Pt.GameCyclesRoom[gamename].JoinStatus == false {
			infoChTmp.Ch <- infoChTmp.Name + ":房间人数已达上限"
			return
		}
		//检查accesskey的输入
		infoChTmp.Ch <- infoChTmp.Name + ":输入要加入的房间的AccessKey"
		var ack string
		if input.Scan() {
			ack = input.Text()
		}
		if ack != Pt.GameCyclesRoom[gamename].Ack {
			infoChTmp.Ch <- infoChTmp.Name + ":AccessKey错误，加入失败"
			return
		}
		//正常赋值
		var tmpData Pt.InfoChListStruct
		tmpData.ChList = Pt.GameCyclesRoom[gamename].ChList
		tmpData.Ack = Pt.GameCyclesRoom[gamename].Ack
		tmpData.ChList = append(tmpData.ChList, infoChTmp)
		Pt.GameCyclesRoom[gamename] = tmpData

		infoChTmp.Ch <- infoChTmp.Name + ":房间加入成功"
		for input.Scan() {
			switch input.Text() {
			// //准备
			// case "ready":
			// 	//debug
			// 	fmt.Println("ready")
			case "exit":
				//debug
				fmt.Println("exit")
				//-----------------
				//如果房间只有自己 那就退出并删除房间
				if len(Pt.GameCyclesRoom[gamename].ChList) == 1 {
					delete(Pt.GameCyclesRoom, gamename)
					infoChTmp.Ch <- infoChTmp.Name + ":退出房间，房间被注销"
					return
				}
				//先判断自己是不是0号元素并且roomleader为true
				if Pt.GameCyclesRoom[gamename].ChList[0].Name == infoChTmp.Name {
					if Pt.GameCyclesRoom[gamename].ChList[0].RoomLeader == true {
						//移交给下一个(2号)元素位的用户
						var tmpData Pt.InfoChListStruct
						tmpData = Pt.GameCyclesRoom[gamename]
						tmpData.ChList = tmpData.ChList[1:]
						tmpData.ChList[0].RoomLeader = true
						Pt.GameCyclesRoom[gamename] = tmpData
						infoChTmp.Ch <- infoChTmp.Name + ":退出房间，房间被移交给其他用户"
						//debug
						fmt.Println(tmpData.ChList[0].RoomLeader)
						//-----------------
						return
					}
				}
				//否则就移交给下一个(2号)元素位的用户
				var tmpData Pt.InfoChListStruct
				tmpData = Pt.GameCyclesRoom[gamename]
				tmpData.ChList = tmpData.ChList[1:]
				tmpData.ChList[0].RoomLeader = true
				Pt.GameCyclesRoom[gamename] = tmpData
				infoChTmp.Ch <- infoChTmp.Name + ":退出房间，房间被移交给其他用户"
				//debug
				fmt.Println(tmpData.ChList[0].RoomLeader)
				//-----------------
				return
			default:
				//debug
				fmt.Println("default")
				for k := range Pt.GameCyclesRoom[gamename].ChList {
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + ":" + input.Text()
				}
			}
		}
		return
	}

	//循环中没有房间name
	infoChTmp.Ch <- "房间不存在，可通过命令listcycles查看"
	return

}

// //RaadyCycles 准备开始
// func ReadyCycles(infoChTmp Pt.ClientChInfo, address string, input *bufio.Scanner){

// }
