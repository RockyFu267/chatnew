package usercmd

import (
	"bufio"
	Pf "chatserver/publicfunc"
	Pt "chatserver/publictype"
	"fmt"
)

//CreateCycles 创建石头剪刀布的房间 1v1
func CreateCycles(infoChTmpData Pt.ClientChInfo, address string, input *bufio.Scanner) {
	var infoChTmp = infoChTmpData
	//判断是否有昵称 没有昵称不能操作
	if infoChTmp.Name == "" {
		infoChTmp.Ch <- address + ": " + "请先输入昵称"
		infoChTmp.Ch <- address + ": " + Pf.Helpstring()
		return
	}
	infoChTmp.Ch <- infoChTmp.Address + ":输入要创建的游戏房间号"
	var gamename string
	if input.Scan() {
		gamename = input.Text()
	}
	//合法性检查
	judge := Pf.JudgeStringSpecialSymbol(gamename)
	if !judge {
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
	if !judgeack {
		infoChTmp.Ch <- infoChTmp.Name + ":AccessKey只支持大小写A-z以及0-9,长度不超过20,不小于2"
		return
	}
	//wait ----初始化其他信息 因为是公用的对象
	//正常赋值
	infoChTmp.RoomLeader = true
	infoChTmp.ReadyStatus = true
	var tmpData Pt.InfoChListStruct
	tmpData.ChList = append(tmpData.ChList, &infoChTmp)
	tmpData.Ack = ack
	tmpData.JoinStatus = true
	//创建房间的chan
	var tmpCH = make(chan string, 1)
	Pt.CyclesRoomChMap["cycles"+gamename] = tmpCH
	//数值初始化--join时候需要，目前是1v1 后者加入游戏即开始,多人应加上房主标签，其他人ready状态全true，房主可以start
	Pt.GameCyclesRoom[gamename] = tmpData
	// infoChTmp.Ch <- infoChTmp.Name + ":房间创建成功，可通过命令listroom查看"
	infoChTmp.Ch <- infoChTmp.Name + ":房间创建成功，等待对手"
	//进入游戏房间随时开始
	for input.Scan() {
		CyclesInputScan(&infoChTmp, gamename, input, ack)
	}
	//主动或被动断开连接退出房间或者直接判负
	//wait
	//自己断线 重新赋值
	var tmpDataTMP Pt.InfoChListStruct
	tmpDataTMP.ChList = Pt.GameCyclesRoom[gamename].ChList
	tmpDataTMP.Ack = ack
	tmpDataTMP.JoinStatus = true
	tmpDataTMP.GameStatus = false
	tmpDataTMP.ConnectBroken = true
	for k := range tmpDataTMP.ChList {
		if tmpDataTMP.ChList[k].Name == infoChTmp.Name {
			tmpDataTMP.ChList = append(tmpDataTMP.ChList[:k], tmpDataTMP.ChList[(k+1):]...)
			break
		}
	}
	Pt.GameCyclesRoom[gamename] = tmpDataTMP
	//如果房间只有自己 那就退出并删除房间
	if len(Pt.GameCyclesRoom[gamename].ChList) == 0 {
		delete(Pt.GameCyclesRoom, gamename)
		return
	}
}

//JoinCycles 加入房间命令
func JoinCycles(infoChTmpData Pt.ClientChInfo, address string, input *bufio.Scanner) {
	var infoChTmp = infoChTmpData
	//判断是否有昵称 没有昵称不能操作
	if infoChTmp.Name == "" {
		infoChTmp.Ch <- address + ": " + "请先输入昵称"
		infoChTmp.Ch <- address + ": " + Pf.Helpstring()
		return
	}
	infoChTmp.Ch <- infoChTmp.Name + ":输入要加入的游戏房间号"
	var gamename string
	if input.Scan() {
		gamename = input.Text()
	}
	//检查是否存在
	if _, ok := Pt.GameCyclesRoom[gamename]; ok {
		for k := range Pt.GameCyclesRoom[gamename].ChList {
			if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
				//2020-11-10 这个逻辑好像不会触发  -待修改
				infoChTmp.Ch <- infoChTmp.Name + ":你已经加入过该房间"
				return
			}
		}
		//房间人数上限或游戏中有人退出(德扑的房间存在被淘汰的玩家离场的情况)
		if len(Pt.GameCyclesRoom[gamename].ChList) >= 2 || !Pt.GameCyclesRoom[gamename].JoinStatus {
			//debug
			fmt.Println(Pt.GameCyclesRoom[gamename].JoinStatus, "----------------")
			//------
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
		//wait ----初始化其他信息 因为是公用的对象
		//正常赋值
		//这段应该可以直接修改值吧  -----在公司的IDE报错------原因待查
		var tmpData Pt.InfoChListStruct
		tmpData.ChList = Pt.GameCyclesRoom[gamename].ChList
		tmpData.Ack = Pt.GameCyclesRoom[gamename].Ack
		tmpData.JoinStatus = Pt.GameCyclesRoom[gamename].JoinStatus
		tmpData.ChList = append(tmpData.ChList, &infoChTmp)
		Pt.GameCyclesRoom[gamename] = tmpData

		infoChTmp.Ch <- infoChTmp.Name + ":房间加入成功"
		//进入游戏房间随时开始
		for input.Scan() {
			CyclesInputScan(&infoChTmp, gamename, input, ack)
		}
		//主动或被动断开连接退出房间或者直接判负
		//wait
		//自己断线 重新赋值
		var tmpDataTMP Pt.InfoChListStruct
		tmpDataTMP.ChList = Pt.GameCyclesRoom[gamename].ChList
		tmpDataTMP.Ack = ack
		tmpDataTMP.JoinStatus = true
		tmpDataTMP.GameStatus = false
		tmpDataTMP.ConnectBroken = true
		for k := range tmpDataTMP.ChList {
			if tmpDataTMP.ChList[k].Name == infoChTmp.Name {
				tmpDataTMP.ChList = append(tmpDataTMP.ChList[:k], tmpDataTMP.ChList[(k+1):]...)
				break
			}
		}
		Pt.GameCyclesRoom[gamename] = tmpDataTMP
		//如果房间只有自己 那就退出并删除房间
		if len(Pt.GameCyclesRoom[gamename].ChList) == 0 {
			delete(Pt.GameCyclesRoom, gamename)
			return
		}
		return
	}

	//循环中没有房间name
	infoChTmp.Ch <- "房间不存在，可通过命令listcycles查看"

}
