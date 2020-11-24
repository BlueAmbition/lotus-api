package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type UserWalletRecord struct {
	Id             int64     `orm:"column(id)"`
	UserId         int64     `orm:"column(user_id)"`
	UserName       string    `orm:"column(user_name)"`
	WalletId       int64     `orm:"column(wallet_id)"`
	CoinId         int64     `orm:"column(coin_id)"`
	CoinEn         string    `orm:"column(coin_en)"`
	UsableChange   float64   `orm:"column(usable_change)"`
	FreezeChange   float64   `orm:"column(freeze_change)"`
	CurrentUsable  float64   `orm:"column(current_usable)"`
	CurrentFreeze  float64   `orm:"column(current_freeze)"`
	Description    string    `orm:"column(description)"`
	Behavior       int64       `orm:"column(behavior)"`
	AdminId        int64     `orm:"column(admin_id)"`
	OptDescription string    `orm:"column(opt_description)"`
	CreatedAt      time.Time `orm:"column(created_at);auto_now_add;type(datetime)"`
	UpdatedAt      time.Time `orm:"column(updated_at);auto_now;type(datetime)"`
}

func (m *UserWalletRecord) TableName() string {
	return "user_wallet_record"
}

func init() {
	orm.RegisterModel(new(UserWalletRecord))
}
