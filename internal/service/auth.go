package service

import (
	"gin-gorm/internal/dao"
)

type AuthSvc struct {
	authDao *dao.AuthDao
}

type AuthReq struct {
	AppKey    string `json:"app_key" binding:"required"`
	AppSecret string `json:"app_secret" binding:"required"`
}

func NewAuthSvc(authDao *dao.AuthDao) *AuthSvc {
	return &AuthSvc{
		authDao: authDao,
	}
}

func (a *AuthSvc) CheckAuth(param *AuthReq) error {
	_, err := a.authDao.GetAuth(param.AppKey, param.AppSecret)
	if err != nil {
		return err
	}
	return nil
}
