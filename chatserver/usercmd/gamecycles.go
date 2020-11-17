package usercmd

import (
	"bufio"
	Pt "chatserver/publictype"
)

//GameCycles 游戏-石头剪刀布
func GameCycles(player Pt.InfoChListStruct) {

}

//CyclesJudge 游戏-石头剪刀布
func CyclesJudge(str1 string, str2 string) {

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

func JudgeRes(play1 *Pt.ClientChInfo, play2 *Pt.ClientChInfo) (winner []*Pt.ClientChInfo) {
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
