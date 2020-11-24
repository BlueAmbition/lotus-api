package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type UserWallet struct {
	Id       int64     `orm:"column(id)"`
	UserId   int64     `orm:"column(user_id)"`
	CoinId   int64     `orm:"column(coin_id)"`
	Usable   float64   `orm:"column(usable)"`
	Freeze   float64   `orm:"column(freeze)"`
	Status   int       `orm:"column(status)"`
	CreateAt time.Time `orm:"column(created_at);auto_now_add;type(datetime)"`
	UpdateAt time.Time `orm:"column(updated_at);auto_now;type(datetime)"`
}

func (m *UserWallet) TableName() string {
	return "user_wallet"
}

func init() {
	orm.RegisterModel(new(UserWallet))
}
