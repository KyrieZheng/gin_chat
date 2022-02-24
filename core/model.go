package core

import "time"

type uid = string

const (
	EventTypeMsg = "event-msg" // 用户发言
	EventTypeSystem = "event-system" // 系统信息推送 房间人数
	EventTypeJoin = "event-join" // 用户加入
	EventTypeTyping = "event-typing" // 用户正在输入中
	EventTypeLeave = "event-leave" // 用户离开
	EventTypeImage = "event-image" // 消息图片
)

// 聊天室事件定义
type Event struct {
	Type string `json:"type"`
	User string `json:"user"`
	Timestamp int64 `json:"timestamp"`
	Text string `json:"text"`
	UserCount int `json:"userCount"`
}

func NewEvent(typ string, user, msg string) Event {
	return Event{
		Type: typ,
		User: user,
		Timestamp: time.Now().UnixNano() / 1e6,
		Text: msg,
	}
}

type Subscription struct {
	Id string // 用户在聊天室中的ID
	UserName string // 用户名
	Pipe <-chan Event // 事件接收通道 用户从这个通道接受消息
	EmitCHn chan Event // 用户消息推送通道
	LeaveCHn chan uid // 用户离开事件推送
}

func (s *Subscription) Leave() {
	s.LeaveCHn <- s.Id
}

func (s *Subscription) Say(message string) {
	s.EmitCHn <- NewEvent(EventTypeMsg, s.UserName, message)
}