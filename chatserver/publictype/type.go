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
	Ch         ClientChan
	Address    string
	Name       string
	Friends    map[string]bool
	RoomLeader bool
	//游戏专用
	Value string
	//是否可以行动
	ActionsStatus bool `json:"actionstatus,omitempty"`
	//准备状态
	ReadyStatus bool `json:"readystatus,omitempty"`
	//胜利 失败 平局
	WinCount  int `json:"wincount,omitempty"`
	LoseCount int `json:"losecount,omitempty"`
	Draw      int `json:"draw,omitempty"`
	//本局胜利
	Winner int `json:"winner,omitempty"`
	//本轮行动状态 是否已行动过
	ActionsHistory bool `json:"actionshistory,omitempty"`
	//德扑专用
	//剩余筹码
	//座次
	//枪口加N位
	//allin状态
	//行为状态 --是否已经弃牌
}

//ChatGroup 定义组room的结构体
type ChatGroup struct {
	Name      string
	AccessKey string
	ChList    []ClientChInfo
}

//InfoChListStruct 管道list的结构体 游戏属性
type InfoChListStruct struct {
	Ack    string
	ChList []*ClientChInfo
	//能否继续加入
	JoinStatus bool
	//GameStatus
	GameStatus bool
	//先出手的值
	//ActionFirst chan string
}

//TMPCyclesCh
var TMPCyclesCh = make(chan string, 1)

//InfoList 初始化tcp连接的数组 后期可以优化改map  list不用考虑并发锁的问题
var InfoList []ClientInfo

//InfoChList 初始化管道的数组 后期可以优化改map	list不用考虑并发锁的问题
var InfoChList []ClientChInfo

//InfoPubChList 初始化公共管道的数组 后期可以优化改map	list不用考虑并发锁的问题
var InfoPubChList []ClientChInfo

//RoomList 初始化组room的数组
var RoomList []ChatGroup

//后期优化可以加其他的事件管道
var (
	Entering = make(chan ClientChan)
	Leaving  = make(chan ClientChan)
	Messages = make(chan string) // all incoming client messages
)

//InfoMap 优化InfoList改map	要考虑并发锁的问题
var InfoMap = make(map[string]ClientInfo)

//InfoChMap 优化InfoChList改map	要考虑并发锁的问题
var InfoChMap = make(map[string]ClientInfo)

//RoomMap 初始化组room的数组	要考虑并发锁的问题
var RoomMap = make(map[string]ChatGroup)

//GameCyclesRoom 石头剪刀布的房间 暂时只支持1v1
var GameCyclesRoom = make(map[string]InfoChListStruct)
