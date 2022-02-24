package core

import (
	"container/list"
	"fmt"

	"github.com/google/uuid"
)

const (
	archiveSize = 20
	chanSize = 10
	msgJoin = "[加入房间]"
	msgLeave = "[离开房间]"
	msgTyping = "[正在输入]"
)

type Room struct {
	users  map[uid]chan Event
	userCount int
	publishChn  chan Event             // 聊天室的消息推送入口
	archive     *list.List             // 历史记录 todo 未持久化 重启失效
	archiveChan chan chan []Event      // 通过接受chan来同步聊天内容
	joinChn     chan chan Subscription // 接收订阅事件的通道 用户加入聊天室后要把历史事件推送给用户
	leaveChn    chan uid               // 用户取消订阅通道 把通道中的历史事件释放并把用户从聊天室用户列表中删除
}

func NewRoom() *Room {
	r := &Room{
		users: map[uid]chan Event{},
		userCount: 0,
		publishChn: make(chan Event, chanSize),
		archiveChan: make(chan chan []Event, chanSize),
		archive: list.New(),
		joinChn: make(chan chan Subscription, chanSize),
		leaveChn: make(chan string, chanSize),
	}
	go r.Serve()

	return r
}

func (r *Room) MsgJoin(user string) {
	fmt.Println("user:", user, "join")
	r.publishChn <- NewEvent(EventTypeJoin, user, msgJoin)
}

func (r *Room) MsgSay(user, message string) {
	r.publishChn <- NewEvent(EventTypeMsg, user, message)
}

func (r *Room) MsgLeave(user string) {
	r.publishChn <- NewEvent(EventTypeLeave, user, msgLeave)
}

func (r *Room) Remove(id uid) {
	r.leaveChn <- id
}

func (r *Room) Join(username string) Subscription {
	resp := make(chan Subscription)
	r.joinChn <- resp
	s := <-resp
	s.UserName = username
	return s
}

func (r *Room) Serve() {
	for {
		select {
		case ch := <- r.joinChn:
			chn := make(chan Event, chanSize)
			r.userCount++
			uid := uuid.New().String()
			r.users[uid] = chn
			ch <- Subscription{
				Id: uid,
				Pipe: chn,
				EmitCHn: r.publishChn,
				LeaveCHn: r.leaveChn,
			}

			ev := NewEvent(EventTypeSystem, "", "")
			ev.UserCount = r.userCount
			for _, v := range r.users {
				v <- ev
			}
			fmt.Println(r.users)
		case event := <-r.publishChn:
			event.UserCount = r.userCount
			for _,v := range r.users {
				v <- event
			}
		case k := <-r.leaveChn:
			if _, ok := r.users[k]; ok {
				delete(r.users, k)
				r.userCount--
			}
			ev := NewEvent(EventTypeSystem, "", "")
			ev.UserCount = r.userCount
			for _, v := range r.users {
				v <- ev
			}
		}
	}
}

