package usercmd

import (
	"bufio"
	Pt "chatserver/publictype"
	"fmt"
	"strconv"
)

//CyclesInputScan 进入石头剪刀布的流程
func CyclesInputScan(infoChTmp *Pt.ClientChInfo, gamename string, input *bufio.Scanner, ack string) {
	//发送心跳给最外层 以及房间的管道
	//判断对方是否已经断开连接
	if len(Pt.GameCyclesRoom[gamename].ChList) == 1 && Pt.GameCyclesRoom[gamename].ConnectBroken {
		infoChTmp.Ch <- infoChTmp.Name + ":对手断开连接，你赢了"
		//正常赋值
		infoChTmp.RoomLeader = true
		infoChTmp.ReadyStatus = true
		infoChTmp.WinCount = 0
		infoChTmp.LoseCount = 0
		infoChTmp.Draw = 0
		infoChTmp.ActionsHistory = false
		infoChTmp.ActionsStatus = true
		var tmpDataTMP Pt.InfoChListStruct
		tmpDataTMP.ChList = append(tmpDataTMP.ChList, infoChTmp)
		tmpDataTMP.Ack = ack
		tmpDataTMP.JoinStatus = true
		tmpDataTMP.GameStatus = false
		Pt.GameCyclesRoom[gamename] = tmpDataTMP

	}
	//如果是比赛状态
	if Pt.GameCyclesRoom[gamename].GameStatus {
		//判断我本轮是否已经操作过了
		if !infoChTmp.ActionsHistory && infoChTmp.ActionsStatus {
			switch input.Text() {
			case "1":
				GameCycles(infoChTmp, "1", gamename)
			case "2":
				GameCycles(infoChTmp, "2", gamename)
			case "3":
				GameCycles(infoChTmp, "3", gamename)
			default:
				infoChTmp.Ch <- infoChTmp.Name + "无效指令,1 对应石头;2 对应剪刀;3 对应布; 请重新输入"
			}
		} else {
			//这里之后应该加入判断 如果没找到那么说明对方退出了聊天室，游戏结束直接胜利
			infoChTmp.Ch <- infoChTmp.Name + ":等待对手做出决定"
		}

	} else {
		switch input.Text() {
		//后期会加转让房主的功能所以这里创建者和加入者的条件判断一致
		case "ready":
			ReadyCycles(infoChTmp, gamename)
		case "start":
			StartCycles(infoChTmp, gamename)
		case "exit":
			ExitCycles(infoChTmp, gamename)
			return
		default:
			for k := range Pt.GameCyclesRoom[gamename].ChList {
				Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "在游戏房" + gamename + "说:" + input.Text()
			}
		}
	}
}

//GameCycles 游戏-石头剪刀布
func GameCycles(infoChTmp *Pt.ClientChInfo, input string, gamename string) {
	infoChTmp.Value = input
	strPlay1 := TypeNameRes(infoChTmp.Value)
	//输入1 等待比较 返回结果
	//判断对方是否已经输入 取房间管道的值长度
	if len(Pt.CyclesRoomChMap["cycles"+gamename]) == 1 {
		for k := range Pt.GameCyclesRoom[gamename].ChList {
			//找到对手的指针
			if Pt.GameCyclesRoom[gamename].ChList[k].Name != infoChTmp.Name {
				//先赋值对手的值
				Pt.GameCyclesRoom[gamename].ChList[k].Value = <-Pt.CyclesRoomChMap["cycles"+gamename]
				strPlay2 := TypeNameRes(Pt.GameCyclesRoom[gamename].ChList[k].Value)
				//先比大小然后输出结果
				res := JudgeCyclesRes(infoChTmp, Pt.GameCyclesRoom[gamename].ChList[k])
				//如果平局
				if len(res) > 1 {
					//初始化
					infoChTmp.ActionsHistory = false
					infoChTmp.ActionsStatus = true
					infoChTmp.Draw = infoChTmp.Draw + 1
					Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
					Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
					Pt.GameCyclesRoom[gamename].ChList[k].Draw = Pt.GameCyclesRoom[gamename].ChList[k].Draw + 1
					for k := range Pt.GameCyclesRoom[gamename].ChList {
						if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
							continue
						}
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "双方平局"
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "当前战绩:" + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "当前战绩:" + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
					}
					infoChTmp.Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
					infoChTmp.Ch <- "双方平局"
					infoChTmp.Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "当前战绩:" + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
					infoChTmp.Ch <- infoChTmp.Name + "当前战绩:" + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
					//跳出循环
					break
				}
				//自己胜利
				if res[0].Name == infoChTmp.Name {
					//初始化
					infoChTmp.ActionsHistory = false
					infoChTmp.ActionsStatus = true
					infoChTmp.WinCount = infoChTmp.WinCount + 1
					Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
					Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
					Pt.GameCyclesRoom[gamename].ChList[k].LoseCount = Pt.GameCyclesRoom[gamename].ChList[k].LoseCount + 1
					for k := range Pt.GameCyclesRoom[gamename].ChList {
						if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
							continue
						}
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "胜利"
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "当前战绩:" + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "当前战绩:" + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
					}
					infoChTmp.Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
					infoChTmp.Ch <- infoChTmp.Name + "胜利"
					infoChTmp.Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "当前战绩:" + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
					infoChTmp.Ch <- infoChTmp.Name + "当前战绩:" + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
					//跳出循环
					break
				}
				//对方胜利
				//初始化
				infoChTmp.ActionsHistory = false
				infoChTmp.ActionsStatus = true
				infoChTmp.LoseCount = infoChTmp.LoseCount + 1
				Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
				Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
				Pt.GameCyclesRoom[gamename].ChList[k].WinCount = Pt.GameCyclesRoom[gamename].ChList[k].WinCount + 1
				for k := range Pt.GameCyclesRoom[gamename].ChList {
					if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
						continue
					}
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "胜利"
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "当前战绩:" + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "当前战绩:" + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)

				}
				infoChTmp.Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
				infoChTmp.Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "胜利"
				infoChTmp.Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "当前战绩:" + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
				infoChTmp.Ch <- infoChTmp.Name + "当前战绩:" + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)

				//跳出循环
				break

			}

		}
		//这里之后应该加入sign 如果没找到那么说明对方退出了聊天室，游戏结束直接胜利
		//如果比赛局数达到N局 对局结束房间信息重置
		if infoChTmp.WinCount+infoChTmp.LoseCount+infoChTmp.Draw >= 10 {
			for k := range Pt.GameCyclesRoom[gamename].ChList {
				if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
					continue
				}
				Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "本局最终战绩:" + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
				Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "本局最终战绩:" + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
				infoChTmp.Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "本局最终战绩:" + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
				infoChTmp.Ch <- infoChTmp.Name + "本局最终战绩:" + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
				if Pt.GameCyclesRoom[gamename].ChList[k].WinCount == infoChTmp.WinCount {
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "最终双方平手"
					infoChTmp.Ch <- "最终双方平手"
				}
				if Pt.GameCyclesRoom[gamename].ChList[k].WinCount > infoChTmp.WinCount {
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "最终的BattleKing是:" + Pt.GameCyclesRoom[gamename].ChList[k].Name
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "恭喜" + Pt.GameCyclesRoom[gamename].ChList[k].Name
					infoChTmp.Ch <- "最终的BattleKing是:" + Pt.GameCyclesRoom[gamename].ChList[k].Name
					infoChTmp.Ch <- "恭喜" + Pt.GameCyclesRoom[gamename].ChList[k].Name
				}
				if Pt.GameCyclesRoom[gamename].ChList[k].WinCount < infoChTmp.WinCount {
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "最终的BattleKing是:" + infoChTmp.Name
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "恭喜" + infoChTmp.Name
					infoChTmp.Ch <- "最终的BattleKing是:" + infoChTmp.Name
					infoChTmp.Ch <- "恭喜" + infoChTmp.Name
				}
			}
			for k := range Pt.GameCyclesRoom[gamename].ChList {
				Pt.GameCyclesRoom[gamename].ChList[k].WinCount = 0
				Pt.GameCyclesRoom[gamename].ChList[k].LoseCount = 0
				Pt.GameCyclesRoom[gamename].ChList[k].Draw = 0
				Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
				Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
				if !Pt.GameCyclesRoom[gamename].ChList[k].RoomLeader {
					Pt.GameCyclesRoom[gamename].ChList[k].ReadyStatus = false
				}
			}
			var tmpDataTMP Pt.InfoChListStruct
			tmpDataTMP.ChList = Pt.GameCyclesRoom[gamename].ChList
			tmpDataTMP.Ack = Pt.GameCyclesRoom[gamename].Ack
			tmpDataTMP.JoinStatus = false
			tmpDataTMP.GameStatus = false
			Pt.GameCyclesRoom[gamename] = tmpDataTMP
		}
		return
		//更新所有玩家状态以及初始化房间
	}
	infoChTmp.Value = input
	strPlay1TMP := TypeNameRes(infoChTmp.Value)
	//快乐输出
	fmt.Println(infoChTmp.Name + ":" + strPlay1TMP)
	Pt.CyclesRoomChMap["cycles"+gamename] <- infoChTmp.Value
	//对方没输入
	infoChTmp.Ch <- infoChTmp.Name + ":你选择出" + strPlay1TMP + " 等待对手做出决定"
	for k := range Pt.GameCyclesRoom[gamename].ChList {
		if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
			continue
		}
		Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "已做出决定,请你选择"
	}
	infoChTmp.ActionsHistory = true
	infoChTmp.ActionsStatus = false

}

//ReadyCycles 准备开始
func ReadyCycles(infoChTmp *Pt.ClientChInfo, gamename string) {
	//判断自己是不是房主 是房主不能准备 只能开始
	for k := range Pt.GameCyclesRoom[gamename].ChList {
		if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
			if Pt.GameCyclesRoom[gamename].ChList[k].RoomLeader {
				infoChTmp.Ch <- infoChTmp.Name + ":你是房主，不能准备(可以在全员准备情况下开始)"
				break
			} else {
				infoChTmp.Ch <- infoChTmp.Name + ":你已经准备(在全员准备情况下房主才能开始)"
				for k := range Pt.GameCyclesRoom[gamename].ChList {
					if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
						continue
					}
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "已准备"
				}
				Pt.GameCyclesRoom[gamename].ChList[k].ReadyStatus = true
				//更改选手状态
				infoChTmp.ActionsHistory = false
				infoChTmp.ActionsStatus = true
				break
			}
		}
	}
}

//StartCycles 开始游戏
func StartCycles(infoChTmp *Pt.ClientChInfo, gamename string) {
	if len(Pt.GameCyclesRoom[gamename].ChList) == 1 {
		infoChTmp.Ch <- infoChTmp.Name + ":等待对手加入房间"
		return
	}
	var sign bool
	//判断自己是不是房主 是房主才能开始
	for k := range Pt.GameCyclesRoom[gamename].ChList {
		if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
			if Pt.GameCyclesRoom[gamename].ChList[k].RoomLeader {
				//先判断是不是所有人都准备了
				for j := range Pt.GameCyclesRoom[gamename].ChList {
					if Pt.GameCyclesRoom[gamename].ChList[j].ReadyStatus {
						continue
					} else {
						//标记存在没准备的
						sign = true
						//提示还没有准备的玩家
						for k := range Pt.GameCyclesRoom[gamename].ChList {
							Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[j].Name + "还未准备"
						}
					}
				}
				if sign {
					break
				}
				//更改房间比赛状态
				tmpstruct := Pt.GameCyclesRoom[gamename]
				tmpstruct.GameStatus = true
				Pt.GameCyclesRoom[gamename] = tmpstruct
				//更改选手状态
				infoChTmp.ActionsHistory = false
				infoChTmp.ActionsStatus = true
				//全员就绪
				for k := range Pt.GameCyclesRoom[gamename].ChList {
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "比赛开始"
				}
			} else {
				infoChTmp.Ch <- infoChTmp.Name + ":你不是房主，没有开始权限"
				break
			}
		}
	}
}

//ExitCycles 退出房间
func ExitCycles(infoChTmp *Pt.ClientChInfo, gamename string) {
	//如果房间只有自己 那就退出并删除房间
	if len(Pt.GameCyclesRoom[gamename].ChList) == 1 {
		delete(Pt.GameCyclesRoom, gamename)
		infoChTmp.Ch <- infoChTmp.Name + ":退出房间，房间被注销"
		return
	}

	var tmpData = Pt.GameCyclesRoom[gamename]
	//记录修改前的状态
	var sign bool
	for k := range tmpData.ChList {
		if tmpData.ChList[k].Name == infoChTmp.Name {
			//判断自己还是不是房主
			if tmpData.ChList[k].RoomLeader {
				sign = true
			}
			tmpData.ChList = append(tmpData.ChList[:k], tmpData.ChList[k+1:]...)
			break
		}
	}
	//还是房主,就移交给目前0元素位的用户
	if sign {
		tmpData.ChList[0].RoomLeader = true
		tmpData.ChList[0].ReadyStatus = true
	}
	Pt.GameCyclesRoom[gamename] = tmpData
	infoChTmp.Ch <- infoChTmp.Name + ":退出房间，房间被移交给其他用户"
}

//TypeNameRes 翻译结果
func TypeNameRes(Input string) (OutPut string) {
	switch Input {
	case "1":
		return "石头/rock"
	case "2":
		return "剪刀/scissors"
	case "3":
		return "布/paper"
	default:
		return ""
	}
}

//JudgeCyclesRes 判断胜负
func JudgeCyclesRes(play1 *Pt.ClientChInfo, play2 *Pt.ClientChInfo) (winner []*Pt.ClientChInfo) {
	switch play1.Value {
	case "1":
		switch play2.Value {
		case "1":
			winner = append(winner, play1, play2)
			return winner
		case "2":
			winner = append(winner, play1)
			return winner
		case "3":
			winner = append(winner, play2)
			return winner
		}
	case "2":
		switch play2.Value {
		case "1":
			winner = append(winner, play2)
			return winner
		case "2":
			winner = append(winner, play1, play2)
			return winner
		case "3":
			winner = append(winner, play1)
			return winner
		}
	case "3":
		switch play2.Value {
		case "1":
			winner = append(winner, play1)
			return winner
		case "2":
			winner = append(winner, play2)
			return winner
		case "3":
			winner = append(winner, play1, play2)
			return winner
		}
	default:
		return winner
	}

	return winner
}
