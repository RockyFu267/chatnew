package usercmd

import (
	"bufio"
	Pf "chatserver/publicfunc"
	Pt "chatserver/publictype"
	"fmt"
	"strconv"
)

//CreateCycles 创建石头剪刀布的房间 1v1
func CreateCycles(infoChTmpData Pt.ClientChInfo, address string, input *bufio.Scanner) {
	var infoChTmp Pt.ClientChInfo
	infoChTmp = infoChTmpData
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
	//wait ----初始化其他信息 因为是公用的对象
	//正常赋值
	infoChTmp.RoomLeader = true
	infoChTmp.ReadyStatus = true
	var tmpData Pt.InfoChListStruct
	tmpData.ChList = append(tmpData.ChList, &infoChTmp)
	tmpData.Ack = ack
	tmpData.JoinStatus = true
	//数值初始化--join时候需要，目前是1v1 后者加入游戏即开始,多人应加上房主标签，其他人ready状态全true，房主可以start
	Pt.GameCyclesRoom[gamename] = tmpData
	// infoChTmp.Ch <- infoChTmp.Name + ":房间创建成功，可通过命令listroom查看"
	infoChTmp.Ch <- infoChTmp.Name + ":房间创建成功，等待对手"
	//进入游戏房间随时开始
	for input.Scan() {
		if Pt.GameCyclesRoom[gamename].GameStatus == true {
			if infoChTmp.ActionsHistory == false && infoChTmp.ActionsStatus == true {
				switch input.Text() {
				case "1":
					infoChTmp.Value = "1"
					strPlay1 := TypeNameRes(infoChTmp.Value)
					//输入1 等待比较 返回结果
					//判断对方是否已经输入
					if len(Pt.TMPCyclesCh) == 1 {
						for k := range Pt.GameCyclesRoom[gamename].ChList {
							//找到对手的指针
							if Pt.GameCyclesRoom[gamename].ChList[k].Name != infoChTmp.Name {
								//先赋值对手的值
								Pt.GameCyclesRoom[gamename].ChList[k].Value = <-Pt.TMPCyclesCh
								strPlay2 := TypeNameRes(Pt.GameCyclesRoom[gamename].ChList[k].Value)
								//先比大小然后输出结果
								res := JudgeRes(&infoChTmp, Pt.GameCyclesRoom[gamename].ChList[k])
								if len(res) > 1 {
									fmt.Println("双方平局")
									for k := range Pt.GameCyclesRoom[gamename].ChList {
										Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
										Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "双方平局"
									}
									//初始化
									infoChTmp.ActionsHistory = false
									infoChTmp.ActionsStatus = true
									infoChTmp.Draw = infoChTmp.Draw + 1
									Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
									Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
									Pt.GameCyclesRoom[gamename].ChList[k].Draw = Pt.GameCyclesRoom[gamename].ChList[k].Draw + 1
									//跳出循环
									break
								} else {
									fmt.Println("winner is " + res[0].Name)
									fmt.Println("winner issssss " + infoChTmp.Name)
									if res[0].Name == infoChTmp.Name {
										for k := range Pt.GameCyclesRoom[gamename].ChList {
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "胜利"
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
										}
										//初始化
										infoChTmp.ActionsHistory = false
										infoChTmp.ActionsStatus = true
										infoChTmp.WinCount = infoChTmp.WinCount + 1
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
										Pt.GameCyclesRoom[gamename].ChList[k].LoseCount = Pt.GameCyclesRoom[gamename].ChList[k].LoseCount + 1
										//跳出循环
										break
									} else {
										for k := range Pt.GameCyclesRoom[gamename].ChList {
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "胜利"
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)

										}
										//初始化
										infoChTmp.ActionsHistory = false
										infoChTmp.ActionsStatus = true
										infoChTmp.LoseCount = infoChTmp.LoseCount + 1
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
										Pt.GameCyclesRoom[gamename].ChList[k].WinCount = Pt.GameCyclesRoom[gamename].ChList[k].WinCount + 1
										//跳出循环
										break
									}
								}
							}

						}
						//这里之后应该加入sign 如果没找到那么说明对方退出了聊天室，游戏结束直接胜利
						continue
						//更新所有玩家状态以及初始化房间
					} else {
						infoChTmp.Value = "1"
						strPlay1 := TypeNameRes(infoChTmp.Value)
						Pt.TMPCyclesCh <- infoChTmp.Value
						//对方没输入
						infoChTmp.Ch <- infoChTmp.Name + ":你选择出" + strPlay1 + " 等待对手做出决定"
						for k := range Pt.GameCyclesRoom[gamename].ChList {
							if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
								continue
							}
							Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "已做出决定,请你选择"
						}
						infoChTmp.ActionsHistory = true
						infoChTmp.ActionsStatus = false
					}
					fmt.Println("11111111111")
				case "2":
					infoChTmp.Value = "2"
					strPlay1 := TypeNameRes(infoChTmp.Value)
					//输入1 等待比较 返回结果
					//判断对方是否已经输入
					if len(Pt.TMPCyclesCh) == 1 {
						for k := range Pt.GameCyclesRoom[gamename].ChList {
							//找到对手的指针
							if Pt.GameCyclesRoom[gamename].ChList[k].Name != infoChTmp.Name {
								//先赋值对手的值
								Pt.GameCyclesRoom[gamename].ChList[k].Value = <-Pt.TMPCyclesCh
								strPlay2 := TypeNameRes(Pt.GameCyclesRoom[gamename].ChList[k].Value)
								//先比大小然后输出结果
								res := JudgeRes(&infoChTmp, Pt.GameCyclesRoom[gamename].ChList[k])
								if len(res) > 1 {
									fmt.Println("双方平局")
									for k := range Pt.GameCyclesRoom[gamename].ChList {
										Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
										Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "双方平局"
									}
									//初始化
									infoChTmp.ActionsHistory = false
									infoChTmp.ActionsStatus = true
									infoChTmp.Draw = infoChTmp.Draw + 1
									Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
									Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
									Pt.GameCyclesRoom[gamename].ChList[k].Draw = Pt.GameCyclesRoom[gamename].ChList[k].Draw + 1
									//跳出循环
									break
								} else {
									fmt.Println("winner is " + res[0].Name)
									fmt.Println("winner issssss " + infoChTmp.Name)
									if res[0].Name == infoChTmp.Name {
										for k := range Pt.GameCyclesRoom[gamename].ChList {
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "胜利"
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
										}
										//初始化
										infoChTmp.ActionsHistory = false
										infoChTmp.ActionsStatus = true
										infoChTmp.WinCount = infoChTmp.WinCount + 1
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
										Pt.GameCyclesRoom[gamename].ChList[k].LoseCount = Pt.GameCyclesRoom[gamename].ChList[k].LoseCount + 1
										//跳出循环
										break
									} else {
										for k := range Pt.GameCyclesRoom[gamename].ChList {
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "胜利"
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)

										}
										//初始化
										infoChTmp.ActionsHistory = false
										infoChTmp.ActionsStatus = true
										infoChTmp.LoseCount = infoChTmp.LoseCount + 1
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
										Pt.GameCyclesRoom[gamename].ChList[k].WinCount = Pt.GameCyclesRoom[gamename].ChList[k].WinCount + 1
										//跳出循环
										break
									}
								}
							}

						}
						//这里之后应该加入sign 如果没找到那么说明对方退出了聊天室，游戏结束直接胜利
						continue
						//更新所有玩家状态以及初始化房间
					} else {
						infoChTmp.Value = "2"
						strPlay1 := TypeNameRes(infoChTmp.Value)
						Pt.TMPCyclesCh <- infoChTmp.Value
						//对方没输入
						infoChTmp.Ch <- infoChTmp.Name + ":你选择出" + strPlay1 + " 等待对手做出决定"
						for k := range Pt.GameCyclesRoom[gamename].ChList {
							if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
								continue
							}
							Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "已做出决定,请你选择"
						}
						infoChTmp.ActionsHistory = true
						infoChTmp.ActionsStatus = false
					}
					fmt.Println("22222222222")
				case "3":
					infoChTmp.Value = "3"
					strPlay1 := TypeNameRes(infoChTmp.Value)
					//输入1 等待比较 返回结果
					//判断对方是否已经输入
					if len(Pt.TMPCyclesCh) == 1 {
						for k := range Pt.GameCyclesRoom[gamename].ChList {
							//找到对手的指针
							if Pt.GameCyclesRoom[gamename].ChList[k].Name != infoChTmp.Name {
								//先赋值对手的值
								Pt.GameCyclesRoom[gamename].ChList[k].Value = <-Pt.TMPCyclesCh
								strPlay2 := TypeNameRes(Pt.GameCyclesRoom[gamename].ChList[k].Value)
								//先比大小然后输出结果
								res := JudgeRes(&infoChTmp, Pt.GameCyclesRoom[gamename].ChList[k])
								if len(res) > 1 {
									fmt.Println("双方平局")
									for k := range Pt.GameCyclesRoom[gamename].ChList {
										Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
										Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "双方平局"
									}
									//初始化
									infoChTmp.ActionsHistory = false
									infoChTmp.ActionsStatus = true
									infoChTmp.Draw = infoChTmp.Draw + 1
									Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
									Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
									Pt.GameCyclesRoom[gamename].ChList[k].Draw = Pt.GameCyclesRoom[gamename].ChList[k].Draw + 1
									//跳出循环
									break
								} else {
									fmt.Println("winner is " + res[0].Name)
									fmt.Println("winner issssss " + infoChTmp.Name)
									if res[0].Name == infoChTmp.Name {
										for k := range Pt.GameCyclesRoom[gamename].ChList {
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "胜利"
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
										}
										//初始化
										infoChTmp.ActionsHistory = false
										infoChTmp.ActionsStatus = true
										infoChTmp.WinCount = infoChTmp.WinCount + 1
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
										Pt.GameCyclesRoom[gamename].ChList[k].LoseCount = Pt.GameCyclesRoom[gamename].ChList[k].LoseCount + 1
										//跳出循环
										break
									} else {
										for k := range Pt.GameCyclesRoom[gamename].ChList {
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "胜利"
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)

										}
										//初始化
										infoChTmp.ActionsHistory = false
										infoChTmp.ActionsStatus = true
										infoChTmp.LoseCount = infoChTmp.LoseCount + 1
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
										Pt.GameCyclesRoom[gamename].ChList[k].WinCount = Pt.GameCyclesRoom[gamename].ChList[k].WinCount + 1
										//跳出循环
										break
									}
								}
							}

						}
						//这里之后应该加入sign 如果没找到那么说明对方退出了聊天室，游戏结束直接胜利
						continue
						//更新所有玩家状态以及初始化房间
					} else {
						infoChTmp.Value = "3"
						strPlay1 := TypeNameRes(infoChTmp.Value)
						Pt.TMPCyclesCh <- infoChTmp.Value
						//对方没输入
						infoChTmp.Ch <- infoChTmp.Name + ":你选择出" + strPlay1 + " 等待对手做出决定"
						for k := range Pt.GameCyclesRoom[gamename].ChList {
							if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
								continue
							}
							Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "已做出决定,请你选择"
						}
						infoChTmp.ActionsHistory = true
						infoChTmp.ActionsStatus = false
					}
					fmt.Println("11111111111")
				default:
					fmt.Println("无效指令")
				}
			} else {
				//这里之后应该加入判断 如果没找到那么说明对方退出了聊天室，游戏结束直接胜利
				infoChTmp.Ch <- infoChTmp.Name + ":等待对手做出决定"
			}

		} else {
			switch input.Text() {
			//后期会加转让房主的功能所以这里创建者和加入者的条件判断一致
			case "ready":
				//判断自己是不是房主 是房主不能准备 只能开始
				for k := range Pt.GameCyclesRoom[gamename].ChList {
					if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
						if Pt.GameCyclesRoom[gamename].ChList[k].RoomLeader == true {
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
				//debug
				fmt.Println("ready")
			//如果还是房主则可以开始
			case "start":
				var sign bool
				//判断自己是不是房主 是房主才能开始
				for k := range Pt.GameCyclesRoom[gamename].ChList {
					if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
						if Pt.GameCyclesRoom[gamename].ChList[k].RoomLeader == true {
							//先判断是不是所有人都准备了
							for j := range Pt.GameCyclesRoom[gamename].ChList {
								if Pt.GameCyclesRoom[gamename].ChList[j].ReadyStatus == true {
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
							if sign == true {
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
							//debug
							fmt.Println("全员已准备")
							//---------
						} else {
							infoChTmp.Ch <- infoChTmp.Name + ":你不是房主，没有开始权限"
							break
						}
					}
				}
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
				var tmpData Pt.InfoChListStruct
				tmpData = Pt.GameCyclesRoom[gamename]
				//记录修改前的状态
				var sign bool
				for k := range tmpData.ChList {
					if tmpData.ChList[k].Name == infoChTmp.Name {
						//判断自己还是不是房主
						if tmpData.ChList[k].RoomLeader == true {
							sign = true
						}
						tmpData.ChList = append(tmpData.ChList[:k], tmpData.ChList[k+1:]...)
						break
					}
				}
				//还是房主,就移交给目前0元素位的用户
				if sign == true {
					tmpData.ChList[0].RoomLeader = true
					tmpData.ChList[0].ReadyStatus = true
				}
				Pt.GameCyclesRoom[gamename] = tmpData
				infoChTmp.Ch <- infoChTmp.Name + ":退出房间，房间被移交给其他用户"
				//debug
				fmt.Println(Pt.GameCyclesRoom[gamename].ChList[0].RoomLeader)
				fmt.Println(Pt.GameCyclesRoom[gamename].ChList[0].Name)
				fmt.Println(len(Pt.GameCyclesRoom[gamename].ChList))
				//-----------------
				return
			default:
				//debug
				fmt.Println("default")
				for k := range Pt.GameCyclesRoom[gamename].ChList {
					Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "在游戏房" + gamename + "说:" + input.Text()
				}
			}
		}

	}
	//主动或被动断开连接退出房间或者直接判负
	//wait
	return
}

//JoinCycles 加入房间命令
func JoinCycles(infoChTmpData Pt.ClientChInfo, address string, input *bufio.Scanner) {
	var infoChTmp Pt.ClientChInfo
	infoChTmp = infoChTmpData
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
				//2020-11-10 这个逻辑好像不会触发  -待修改
				infoChTmp.Ch <- infoChTmp.Name + ":你已经加入过该房间"
				return
			}
		}
		//房间人数上限或游戏中有人退出(德扑的房间存在被淘汰的玩家离场的情况)
		if len(Pt.GameCyclesRoom[gamename].ChList) >= 2 || Pt.GameCyclesRoom[gamename].JoinStatus == false {
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
			if Pt.GameCyclesRoom[gamename].GameStatus == true {
				if infoChTmp.ActionsHistory == false && infoChTmp.ActionsStatus == true {
					switch input.Text() {
					case "1":
						infoChTmp.Value = "1"
						strPlay1 := TypeNameRes(infoChTmp.Value)
						//输入1 等待比较 返回结果
						//判断对方是否已经输入
						if len(Pt.TMPCyclesCh) == 1 {
							for k := range Pt.GameCyclesRoom[gamename].ChList {
								//找到对手的指针
								if Pt.GameCyclesRoom[gamename].ChList[k].Name != infoChTmp.Name {
									//先赋值对手的值
									Pt.GameCyclesRoom[gamename].ChList[k].Value = <-Pt.TMPCyclesCh
									strPlay2 := TypeNameRes(Pt.GameCyclesRoom[gamename].ChList[k].Value)
									//先比大小然后输出结果
									res := JudgeRes(&infoChTmp, Pt.GameCyclesRoom[gamename].ChList[k])
									if len(res) > 1 {
										fmt.Println("双方平局")
										for k := range Pt.GameCyclesRoom[gamename].ChList {
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "双方平局"
										}
										//初始化
										infoChTmp.ActionsHistory = false
										infoChTmp.ActionsStatus = true
										infoChTmp.Draw = infoChTmp.Draw + 1
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
										Pt.GameCyclesRoom[gamename].ChList[k].Draw = Pt.GameCyclesRoom[gamename].ChList[k].Draw + 1
										//跳出循环
										break
									} else {
										fmt.Println("winner is " + res[0].Name)
										fmt.Println("winner issssss " + infoChTmp.Name)
										if res[0].Name == infoChTmp.Name {
											for k := range Pt.GameCyclesRoom[gamename].ChList {
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "胜利"
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
											}
											//初始化
											infoChTmp.ActionsHistory = false
											infoChTmp.ActionsStatus = true
											infoChTmp.WinCount = infoChTmp.WinCount + 1
											Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
											Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
											Pt.GameCyclesRoom[gamename].ChList[k].LoseCount = Pt.GameCyclesRoom[gamename].ChList[k].LoseCount + 1
											//跳出循环
											break
										} else {
											for k := range Pt.GameCyclesRoom[gamename].ChList {
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "胜利"
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)

											}
											//初始化
											infoChTmp.ActionsHistory = false
											infoChTmp.ActionsStatus = true
											infoChTmp.LoseCount = infoChTmp.LoseCount + 1
											Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
											Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
											Pt.GameCyclesRoom[gamename].ChList[k].WinCount = Pt.GameCyclesRoom[gamename].ChList[k].WinCount + 1
											//跳出循环
											break
										}
									}
								}

							}
							//这里之后应该加入sign 如果没找到那么说明对方退出了聊天室，游戏结束直接胜利
							continue
							//更新所有玩家状态以及初始化房间
						} else {
							infoChTmp.Value = "1"
							strPlay1 := TypeNameRes(infoChTmp.Value)
							Pt.TMPCyclesCh <- infoChTmp.Value
							//对方没输入
							infoChTmp.Ch <- infoChTmp.Name + ":你选择出" + strPlay1 + " 等待对手做出决定"
							for k := range Pt.GameCyclesRoom[gamename].ChList {
								if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
									continue
								}
								Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "已做出决定,请你选择"
							}
							infoChTmp.ActionsHistory = true
							infoChTmp.ActionsStatus = false
						}
						fmt.Println("11111111111")
					case "2":
						infoChTmp.Value = "2"
						strPlay1 := TypeNameRes(infoChTmp.Value)
						//输入1 等待比较 返回结果
						//判断对方是否已经输入
						if len(Pt.TMPCyclesCh) == 1 {
							for k := range Pt.GameCyclesRoom[gamename].ChList {
								//找到对手的指针
								if Pt.GameCyclesRoom[gamename].ChList[k].Name != infoChTmp.Name {
									//先赋值对手的值
									Pt.GameCyclesRoom[gamename].ChList[k].Value = <-Pt.TMPCyclesCh
									strPlay2 := TypeNameRes(Pt.GameCyclesRoom[gamename].ChList[k].Value)
									//先比大小然后输出结果
									res := JudgeRes(&infoChTmp, Pt.GameCyclesRoom[gamename].ChList[k])
									if len(res) > 1 {
										fmt.Println("双方平局")
										for k := range Pt.GameCyclesRoom[gamename].ChList {
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "双方平局"
										}
										//初始化
										infoChTmp.ActionsHistory = false
										infoChTmp.ActionsStatus = true
										infoChTmp.Draw = infoChTmp.Draw + 1
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
										Pt.GameCyclesRoom[gamename].ChList[k].Draw = Pt.GameCyclesRoom[gamename].ChList[k].Draw + 1
										//跳出循环
										break
									} else {
										fmt.Println("winner is " + res[0].Name)
										fmt.Println("winner issssss " + infoChTmp.Name)
										if res[0].Name == infoChTmp.Name {
											for k := range Pt.GameCyclesRoom[gamename].ChList {
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "胜利"
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
											}
											//初始化
											infoChTmp.ActionsHistory = false
											infoChTmp.ActionsStatus = true
											infoChTmp.WinCount = infoChTmp.WinCount + 1
											Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
											Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
											Pt.GameCyclesRoom[gamename].ChList[k].LoseCount = Pt.GameCyclesRoom[gamename].ChList[k].LoseCount + 1
											//跳出循环
											break
										} else {
											for k := range Pt.GameCyclesRoom[gamename].ChList {
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "胜利"
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)

											}
											//初始化
											infoChTmp.ActionsHistory = false
											infoChTmp.ActionsStatus = true
											infoChTmp.LoseCount = infoChTmp.LoseCount + 1
											Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
											Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
											Pt.GameCyclesRoom[gamename].ChList[k].WinCount = Pt.GameCyclesRoom[gamename].ChList[k].WinCount + 1
											//跳出循环
											break
										}
									}
								}

							}
							//这里之后应该加入sign 如果没找到那么说明对方退出了聊天室，游戏结束直接胜利
							continue
							//更新所有玩家状态以及初始化房间
						} else {
							infoChTmp.Value = "2"
							strPlay1 := TypeNameRes(infoChTmp.Value)
							Pt.TMPCyclesCh <- infoChTmp.Value
							//对方没输入
							infoChTmp.Ch <- infoChTmp.Name + ":你选择出" + strPlay1 + " 等待对手做出决定"
							for k := range Pt.GameCyclesRoom[gamename].ChList {
								if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
									continue
								}
								Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "已做出决定,请你选择"
							}
							infoChTmp.ActionsHistory = true
							infoChTmp.ActionsStatus = false
						}
						fmt.Println("22222222222")
					case "3":
						infoChTmp.Value = "3"
						strPlay1 := TypeNameRes(infoChTmp.Value)
						//输入1 等待比较 返回结果
						//判断对方是否已经输入
						if len(Pt.TMPCyclesCh) == 1 {
							for k := range Pt.GameCyclesRoom[gamename].ChList {
								//找到对手的指针
								if Pt.GameCyclesRoom[gamename].ChList[k].Name != infoChTmp.Name {
									//先赋值对手的值
									Pt.GameCyclesRoom[gamename].ChList[k].Value = <-Pt.TMPCyclesCh
									strPlay2 := TypeNameRes(Pt.GameCyclesRoom[gamename].ChList[k].Value)
									//先比大小然后输出结果
									res := JudgeRes(&infoChTmp, Pt.GameCyclesRoom[gamename].ChList[k])
									if len(res) > 1 {
										fmt.Println("双方平局")
										for k := range Pt.GameCyclesRoom[gamename].ChList {
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "双方平局"
										}
										//初始化
										infoChTmp.ActionsHistory = false
										infoChTmp.ActionsStatus = true
										infoChTmp.Draw = infoChTmp.Draw + 1
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
										Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
										Pt.GameCyclesRoom[gamename].ChList[k].Draw = Pt.GameCyclesRoom[gamename].ChList[k].Draw + 1
										//跳出循环
										break
									} else {
										fmt.Println("winner is " + res[0].Name)
										fmt.Println("winner issssss " + infoChTmp.Name)
										if res[0].Name == infoChTmp.Name {
											for k := range Pt.GameCyclesRoom[gamename].ChList {
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "胜利"
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)
											}
											//初始化
											infoChTmp.ActionsHistory = false
											infoChTmp.ActionsStatus = true
											infoChTmp.WinCount = infoChTmp.WinCount + 1
											Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
											Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
											Pt.GameCyclesRoom[gamename].ChList[k].LoseCount = Pt.GameCyclesRoom[gamename].ChList[k].LoseCount + 1
											//跳出循环
											break
										} else {
											for k := range Pt.GameCyclesRoom[gamename].ChList {
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "选择出了:" + strPlay2 + " , " + infoChTmp.Name + "选择出了:" + strPlay1
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[k].Name + "胜利"
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + Pt.GameCyclesRoom[gamename].ChList[k].Name + ":\n胜利-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].WinCount) + "\n失败-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].LoseCount) + "\n平局-" + strconv.Itoa(Pt.GameCyclesRoom[gamename].ChList[k].Draw)
												Pt.GameCyclesRoom[gamename].ChList[k].Ch <- "当前战绩:" + infoChTmp.Name + ":\n胜利-" + strconv.Itoa(infoChTmp.WinCount) + "\n失败-" + strconv.Itoa(infoChTmp.LoseCount) + "\n平局-" + strconv.Itoa(infoChTmp.Draw)

											}
											//初始化
											infoChTmp.ActionsHistory = false
											infoChTmp.ActionsStatus = true
											infoChTmp.LoseCount = infoChTmp.LoseCount + 1
											Pt.GameCyclesRoom[gamename].ChList[k].ActionsHistory = false
											Pt.GameCyclesRoom[gamename].ChList[k].ActionsStatus = true
											Pt.GameCyclesRoom[gamename].ChList[k].WinCount = Pt.GameCyclesRoom[gamename].ChList[k].WinCount + 1
											//跳出循环
											break
										}
									}
								}

							}
							//这里之后应该加入sign 如果没找到那么说明对方退出了聊天室，游戏结束直接胜利
							continue
							//更新所有玩家状态以及初始化房间
						} else {
							infoChTmp.Value = "3"
							strPlay1 := TypeNameRes(infoChTmp.Value)
							Pt.TMPCyclesCh <- infoChTmp.Value
							//对方没输入
							infoChTmp.Ch <- infoChTmp.Name + ":你选择出" + strPlay1 + " 等待对手做出决定"
							for k := range Pt.GameCyclesRoom[gamename].ChList {
								if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
									continue
								}
								Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "已做出决定,请你选择"
							}
							infoChTmp.ActionsHistory = true
							infoChTmp.ActionsStatus = false
						}
						fmt.Println("11111111111")
					default:
						fmt.Println("无效指令")
					}
				} else {
					//这里之后应该加入判断 如果没找到那么说明对方退出了聊天室，游戏结束直接胜利
					infoChTmp.Ch <- infoChTmp.Name + ":等待对手做出决定"
				}

			} else {
				switch input.Text() {
				//准备
				case "ready":
					//判断自己是不是房主 是房主不能准备 只能开始
					for k := range Pt.GameCyclesRoom[gamename].ChList {
						if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
							if Pt.GameCyclesRoom[gamename].ChList[k].RoomLeader == true {
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
					//debug
					fmt.Println("ready")
				//如果获得房主则可以开始
				case "start":
					var sign bool
					//判断自己是不是房主 是房主才能开始
					for k := range Pt.GameCyclesRoom[gamename].ChList {
						if Pt.GameCyclesRoom[gamename].ChList[k].Name == infoChTmp.Name {
							if Pt.GameCyclesRoom[gamename].ChList[k].RoomLeader == true {
								//先判断是不是所有人都准备了
								for j := range Pt.GameCyclesRoom[gamename].ChList {
									if Pt.GameCyclesRoom[gamename].ChList[j].ReadyStatus == true {
										continue
									} else {
										//标记存在未准备
										sign = true
										//提示还没有准备的玩家
										for k := range Pt.GameCyclesRoom[gamename].ChList {
											Pt.GameCyclesRoom[gamename].ChList[k].Ch <- Pt.GameCyclesRoom[gamename].ChList[j].Name + "还未准备"
										}
									}
								}
								if sign == true {
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
								//debug
								fmt.Println("全员已准备")
								//---------
							} else {
								infoChTmp.Ch <- infoChTmp.Name + ":你不是房主，没有开始权限"
								break
							}
						}
					}
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
					var tmpData Pt.InfoChListStruct
					tmpData = Pt.GameCyclesRoom[gamename]
					//记录修改前的状态
					var sign bool
					for k := range tmpData.ChList {
						if tmpData.ChList[k].Name == infoChTmp.Name {
							//判断自己还是不是房主
							if tmpData.ChList[k].RoomLeader == true {
								sign = true
							}
							tmpData.ChList = append(tmpData.ChList[:k], tmpData.ChList[k+1:]...)
							break
						}
					}
					//还是房主,就移交给目前0元素位的用户
					if sign == true {
						tmpData.ChList[0].RoomLeader = true
						tmpData.ChList[0].ReadyStatus = true
					}
					Pt.GameCyclesRoom[gamename] = tmpData
					infoChTmp.Ch <- infoChTmp.Name + ":退出房间"
					//debug
					fmt.Println(Pt.GameCyclesRoom[gamename].ChList[0].RoomLeader)
					fmt.Println(Pt.GameCyclesRoom[gamename].ChList[0].Name)
					fmt.Println(len(Pt.GameCyclesRoom[gamename].ChList))
					//-----------------
					return
				default:
					//debug
					fmt.Println("default")
					for k := range Pt.GameCyclesRoom[gamename].ChList {
						Pt.GameCyclesRoom[gamename].ChList[k].Ch <- infoChTmp.Name + "在游戏房" + gamename + "说:" + input.Text()
					}
				}
			}
		}
		//主动或被动断开连接退出房间或者直接判负
		//wait
		return
	}

	//循环中没有房间name
	infoChTmp.Ch <- "房间不存在，可通过命令listcycles查看"
	return

}
