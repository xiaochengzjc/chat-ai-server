package service

import (
	"gin-gorm/internal/dao"
	model "gin-gorm/internal/model"
)

type UserSvc struct {
	userDao *dao.UserDao
}

type UserRequest struct {
	NickName  string `json:"nickName"`
	Gender    uint8  `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
	isDemote  bool   `json:"is_demote"`
	Openid    string `json:"openid"`
	Unionid   string `json:"unionid"`
}

func NewUserSvc(userDao *dao.UserDao) *UserSvc {
	return &UserSvc{
		userDao: userDao,
	}
}

func (s *UserSvc) CheckUserByUID(uid int) bool {
	return s.userDao.FindUserByUid(uid)
}

func (s *UserSvc) AddOrUpdateUser(userRequest *UserRequest) error {
	user := model.User{
		Name:    userRequest.NickName,
		Gender:  userRequest.Gender,
		Avatar:  userRequest.AvatarUrl,
		UnionId: userRequest.Unionid,
		Openid:  userRequest.Openid,
	}
	return s.userDao.FirstOrCreateByOpenid(&user)
}
