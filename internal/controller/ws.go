package controller

import (
	"gin-gorm/global"
	"gin-gorm/internal/service"
	"github.com/gin-gonic/gin"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"strconv"
)

type WSCtr struct {
	userSvc *service.UserSvc
	chatSvc *service.ChatSvc
}

func NewWSCtr(userSvc *service.UserSvc, chatSvc *service.ChatSvc) *WSCtr {
	return &WSCtr{
		userSvc: userSvc,
		chatSvc: chatSvc,
	}
}

func (ws *WSCtr) HandShake(ctx *gin.Context) {
	conn, err := websocket.Accept(ctx.Writer, ctx.Request, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		global.Logger.Errorf("websocket握手失败：%v", err)
		return
	}

	uid := 0
	uid, err = strconv.Atoi(ctx.Request.FormValue("uid"))
	if err != nil || uid <= 0 || !ws.userSvc.CheckUserByUID(uid) {
		wsjson.Write(ctx, conn, "消息结构体")
		conn.Close(websocket.StatusUnsupportedData, "用户信息错误或用户不存在")
		global.Logger.Errorf("用户信息错误或用户不存在")
		return
	}

	user := service.NewWsUser(ws.chatSvc, conn, uid, ctx.Request.RemoteAddr)

	go user.SendMessage(ctx)

	defer close(user.MessageChannel)

	//服务端心跳
	//defer user.Ping()
	//defer close(user.CloseChannel)

	err = user.ReceiveMessage(ctx)

	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		global.Logger.Errorf("读取websocket客户端错误：%v", err)
		conn.Close(websocket.StatusInternalError, "读取websocket客户端错误")
	}
}
