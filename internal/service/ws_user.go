package service

import (
	"context"
	"gin-gorm/global"
	"gin-gorm/pkg/util"
	"gin-gorm/pkg/ws"
	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

type WsUser struct {
	UID            int             `json:"uid"`
	EnterAt        util.JsonTime   `json:"enter_at"`
	Addr           string          `json:"addr"`
	MessageChannel chan *Message   `json:"-"`
	conn           *websocket.Conn `json:"-"`
	chatSvc        *ChatSvc        `json:"-"`
	CloseChannel   chan struct{}   `json:"-"`
}

type Message struct {
	MsgId   string        `json:"msg_id"`
	User    *WsUser       `json:"user"`
	Type    string        `json:"type"`
	Content string        `json:"content"`
	MsgTime util.JsonTime `json:"msg_time"`
}

func NewWsUser(chatSvc *ChatSvc, conn *websocket.Conn, uid int, addr string) *WsUser {
	return &WsUser{
		UID:            uid,
		EnterAt:        util.JsonTime(time.Now()),
		Addr:           addr,
		conn:           conn,
		MessageChannel: make(chan *Message),
		chatSvc:        chatSvc,
		CloseChannel:   make(chan struct{}),
	}
}

func (u *WsUser) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		err := wsjson.Write(ctx, u.conn, msg)
		if err != nil {
			u.CloseChannel <- struct{}{}
		}
	}
}
func (u *WsUser) ReceiveMessage(ctx context.Context) error {
	var (
		err error
	)
	msg := Message{}
	weighted := semaphore.NewWeighted(1)
	for {
		err = wsjson.Read(ctx, u.conn, &msg)
		if err != nil {
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) { //正常关闭，非错误
				return nil
			}
			return err
		}
		go u.HandleMsg(ctx, &Message{MsgId: msg.MsgId, Content: msg.Content, Type: msg.Type}, weighted)
	}
}

func (u *WsUser) HandleMsg(ctx context.Context, msg *Message, weighted *semaphore.Weighted) {
	//客户端请求消息
	if msg.Type == ws.MsgTypeReq {
		if ok := weighted.TryAcquire(1); ok { //限制单个用户同时只能请求一次
			u.chatSvc.Chat(ctx, msg, u.MessageChannel)
			defer weighted.Release(1)
		} else {
			u.MessageChannel <- &Message{
				Type:    ws.MsgTypeError,
				Content: "请求频繁，请稍后再试",
				MsgId:   msg.MsgId,
			}
		}
	}

	//心跳包
	if msg.Type == ws.MsgTypeHeart {
		u.chatSvc.Heart(ctx, msg, u.MessageChannel)
	}
}

func (u *WsUser) Ping() {
	ticker := time.NewTicker(ws.PingPeriod)
	defer func() {
		ticker.Stop()
		u.conn.Close(websocket.StatusNormalClosure, "服务器断开")
	}()

	for {
		select {
		case <-ticker.C:
			u.MessageChannel <- &Message{
				Type:    ws.MsgTypeHeart,
				Content: "ping",
			}
		case <-u.CloseChannel:
			global.Logger.Info("websocket服务端关闭")
			return
		}
	}
}
