package dao

import (
	"gin-gorm/internal/model"
)

type UserDao struct {
	*Dao
}

func NewUserDao(dao *Dao) *UserDao {
	return &UserDao{
		Dao: dao,
	}
}

func (d *UserDao) FirstOrCreateByOpenid(user *model.User) error {
	tx := d.db.Where(model.User{Openid: user.Openid}).Assign(*user).FirstOrCreate(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *UserDao) AddUser(user *model.User) error {
	tx := d.db.Create(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *UserDao) UpdateUser(user *model.User) error {
	tx := d.db.Save(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *UserDao) FindUserByOpenid(openid string) (model.User, error) {
	user := model.User{}
	tx := d.db.Where("openid = ?", openid).First(&user)
	if tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (d *UserDao) FindUserByUid(uid int) bool {
	user := model.User{}
	tx := d.db.Where("id", uid).First(&user)
	if tx.Error != nil {
		return false
	}
	return true

}
