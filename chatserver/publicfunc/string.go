package publicfunc

import (
	"strings"
)

//JudgeStringSpecialSymbol 判断用户名合法性
func JudgeStringSpecialSymbol(input string) bool {
	f := func(r rune) bool {
		return (r < 'A' && r > '9') || r > 'z' || (r > 'Z' && r < 'a') || r < '0'
	}
	if strings.IndexFunc(input, f) != -1 {
		return false
	}
	if len(input) >= 20 {
		return false
	}
	return true

}

//StringToDestinationAddr 截取@的地址
func StringToDestinationAddr(input string) (output string) {
	for k := range string(input) {
		if string(input[k]) == " " {
			output = string(input)[1:k]
			break
		}

	}
	return output
}

//StringToDestinationContent 截取@的内容
func StringToDestinationContent(input string) (output string) {
	for k := range string(input) {
		if string(input[k]) == " " {
			output = string(input)[k+1:]
			break
		}

	}
	return output
}

//返回help
func Helpstring() string {
	return (`
	please choose options:
		- createroom : 创建房间
		- joinroom   : 加入房间
		- listroom   : 获取所有房间号
		- listuser   : 获取所有在线用户Name
		- myname     : 注册自己的聊天昵称
        `)
}
