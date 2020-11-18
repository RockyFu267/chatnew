package usercmd

import (
	"bufio"
	Pt "chatserver/publictype"
	"fmt"
	"strconv"
)

//GameCycles 游戏-石头剪刀布
func GameCycles(infoChTmp *Pt.ClientChInfo, input string, gamename string) {
	infoChTmp.Value = input
	strPlay1 := TypeNameRes(infoChTmp.Value)
	//输入1 等待比较 返回结果
	//判断对方是否已经输入
	if len(Pt.CyclesRoomChMap["cycles"+gamename]) == 1 {
		for k := range Pt.GameCyclesRoom[gamename].ChList {
			//找到对手的指针
			if Pt.GameCyclesRoom[gamename].ChList[k].Name != infoChTmp.Name {
				//先赋值对手的值
				Pt.GameCyclesRoom[gamename].ChList[k].Value = <-Pt.CyclesRoomChMap["cycles"+gamename]
				strPlay2 := TypeNameRes(Pt.GameCyclesRoom[gamename].ChList[k].Value)
				//先比大小然后输出结果
				res := JudgeCyclesRes(infoChTmp, Pt.GameCyclesRoom[gamename].ChList[k])
				if len(res) > 1 {
					fmt.Println("双方平局")
					//初始化
					infoChTmp.ActionsHistory = false
					infoChTmp.ActionsStatus = true
					infoChTmp.Draw = infoChTmp.Draw + 1
					Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
					Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
					Pt.GameCyclesRoom[gamename].ChList[k].Draw = Pt.GameCyclesRoom[gamename].ChList[k].Draw + 1
					fmt.Println("双方平局")
					for k := range Pt.GameCyclesRoom[gamename].ChList {
						if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
							continue
						}
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "双方平局"
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:\n" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:\n" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
					}
					infoChTmp.Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
					infoChTmp.Ch <- "双方平局"
					infoChTmp.Ch <- "当前战绩:\n" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
					infoChTmp.Ch <- "当前战绩:\n" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
					//跳出循环
					return
				}
				fmt.Println("winner is " + res[0].Name)
				fmt.Println("winner issssss " + infoChTmp.Name)
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
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:\n" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:\n" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
					}
					infoChTmp.Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
					infoChTmp.Ch <- infoChTmp.Name + "胜利"
					infoChTmp.Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
					infoChTmp.Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
					//跳出循环
					return
				}
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
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:\n" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:\n" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)

				}
				infoChTmp.Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
				infoChTmp.Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "胜利"
				infoChTmp.Ch <- "当前战绩:\n" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
				infoChTmp.Ch <- "当前战绩:\n" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)

				//跳出循环
				return

			}

		}
		//这里之后应该加入sign 如果没找到那么说明对方退出了聊天室，游戏结束直接胜利
		return
		//更新所有玩家状态以及初始化房间
	}
	infoChTmp.Value = input
	strPlay1TMP := TypeNameRes(infoChTmp.Value)
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

	return
}

//ReadyCycles 准备开始
func ReadyCycles(infoChTmp Pt.ClientChInfo, address string, input *bufio.Scanner) {

}

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
