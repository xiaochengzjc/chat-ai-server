package dao

import (
	"gin-gorm/global"
	"gin-gorm/internal/model"
)

type AuthDao struct {
	*Dao
}

func NewAuthDao(dao *Dao) *AuthDao {
	return &AuthDao{
		Dao: dao,
	}
}

func (d *AuthDao) GetAuth(appKey string, appSecret string) (model.Auth, error) {
	var auth model.Auth
	db := d.db.Where("app_key = ? AND app_secret=?", appKey, appSecret)
	err := db.First(&auth).Error
	if err != nil {
		global.Logger.Errorf("AuthDao.GetAuth.查询数据库错误：%v", err)
		return auth, err
	}
	return auth, nil
}
