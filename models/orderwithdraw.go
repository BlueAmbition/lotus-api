package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type OrderWithdraw struct {
	Id           int64     `orm:"column(id);"`
	OrderNo      string    `orm:"column(order_no);"`
	UserId       int64     `orm:"column(user_id);"`
	UserName     string    `orm:"column(user_name);"`
	CoinId       int64     `orm:"column(coin_id);"`
	CoinEn       string    `orm:"column(coin_en);"`
	Address      string    `orm:"column(address);"`
	Memo         string    `orm:"column(memo);"`
	Amount       float64   `orm:"column(amount);"`
	Fee          float64   `orm:"column(fee);"`
	RealAmount   float64   `orm:"column(real_amount);"`
	Description  string    `orm:"column(description);"`
	AutoWithdraw int       `orm:"column(auto_withdraw);"`
	Review       int       `orm:"column(review);"`
	ReviewMark   string    `orm:"column(review_mark);"`
	Status       int       `orm:"column(status);"`
	CreatedAt    time.Time `orm:"column(created_at);auto_now_add;type(datetime)"`
	UpdatedAt    time.Time `orm:"column(updated_at);auto_now;type(datetime)"`
}

func (m *OrderWithdraw) TableName() string {
	return "order_withdraw"
}

func init() {
	orm.RegisterModel(new(OrderWithdraw))
}
