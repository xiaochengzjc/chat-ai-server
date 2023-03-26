package service

import (
	"context"
	"gin-gorm/global"
	"gin-gorm/internal/dao"
	"gin-gorm/internal/model"
	"gin-gorm/pkg/app"
	"gin-gorm/pkg/errcode"
	"gin-gorm/pkg/wechat"
)

type WechatSvc struct {
	userDao *dao.UserDao
	authSvc *AuthSvc
}

type WechatRawData struct {
	NickName  string `json:"nickName"`
	Gender    uint8  `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
	isDemote  bool   `json:"is_demote"`
}

type WechatLoginReq struct {
	AuthReq
	Code string `json:"code" binding:"required,min=1,max=500"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`     // 用户名
	Openid   string `json:"openid"`   // openid
	UnionId  string `json:"unionid"`  // unionid
	Password string `json:"password"` // 密码
	Phone    string `json:"phone"`    // 手机号
	Avatar   string `json:"avatar"`   // 头像
	Gender   uint8  `json:"gender"`   //性别
}

type WechatLoginRsp struct {
	Token    string   `json:"token"`
	UserInfo UserInfo `json:"userInfo"`
}

func NewWechatSvc(authSvc *AuthSvc, userDao *dao.UserDao) *WechatSvc {
	return &WechatSvc{
		userDao: userDao,
		authSvc: authSvc,
	}
}

func (w *WechatSvc) Login(ctx context.Context, r *WechatLoginReq) (WechatLoginRsp, error) {
	rsp := WechatLoginRsp{}
	miniProgram := wechat.NewMiniProgram()
	session, err := miniProgram.Auth.Session(ctx, r.Code)
	if err != nil {
		return rsp, err
	}

	if session.OpenID == "" {
		global.Logger.Errorf("微信登录失败：%v", session)
		return rsp, errcode.ErrorWechatLogin
	}

	user := model.User{
		Openid:  session.OpenID,
		UnionId: session.UnionID,
	}

	err = w.userDao.FirstOrCreateByOpenid(&user)
	if err != nil {
		return rsp, err
	}

	err = w.authSvc.CheckAuth(&r.AuthReq)
	if err != nil {
		return rsp, err
	}

	token, err := app.GenerateToken(r.AppKey, r.AppSecret)
	if err != nil {
		global.Logger.Errorf("生成token失败：%v", err)
		return rsp, err
	}

	rsp.Token = token
	rsp.UserInfo = UserInfo{
		ID:      user.ID,
		Name:    user.Name,
		Openid:  user.Openid,
		UnionId: user.UnionId,
		Avatar:  user.Avatar,
		Gender:  user.Gender,
	}

	return rsp, nil
}
