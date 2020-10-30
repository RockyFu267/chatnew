package publictype

import (
	"net"
)

//ClientChan 定义双向管道
type ClientChan chan string

//AllClients 初始化所有的连接map
var AllClients = make(map[string]net.Conn)

//ClientInfo 定义tcp连接的结构体
type ClientInfo struct {
	ConnChan net.Conn
	Address  string
	Name     string
}

//ClientChInfo 定义管道的结构体
type ClientChInfo struct {
	Ch      ClientChan
	Address string
	Name    string
}

//ChatGroup 定义组room的结构体
type ChatGroup struct {
	Name   string
	ChList []ClientChInfo
}

//InfoList 初始化tcp连接的数组 后期可以优化改map
var InfoList []ClientInfo

//InfoChList 初始化管道的数组 后期可以优化改map
var InfoChList []ClientChInfo

//RoomList 初始化组room的数组
var RoomList []ChatGroup

//后期优化可以加其他的事件管道
var (
	Entering = make(chan ClientChan)
	Leaving  = make(chan ClientChan)
	Messages = make(chan string) // all incoming client messages
)
