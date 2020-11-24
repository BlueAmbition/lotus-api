package routers

import (
	"lotus-api/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &controllers.MainController{})
	v1 := beego.NewNamespace("/erc20_usdt",
		//全部余额
		beego.NSRouter("/total-balance/", &controllers.LotusController{}, "post:TotalBalance"),
		//账户余额
		beego.NSRouter("/account-balance", &controllers.LotusController{}, "post:AccountBalance"),
		//获取或创建account的地址
		beego.NSRouter("/receive-address", &controllers.LotusController{}, "post:GetReceiveAddress"),
		//提币
		beego.NSRouter("/send-to-address", &controllers.LotusController{}, "post:SendToAddress"),
		//人工入账
		beego.NSRouter("/recharge", &controllers.LotusController{}, "post:Recharge"),
		//覆盖交易
		beego.NSRouter("/recover-tx", &controllers.LotusController{}, "post:RecoverTx"),
	)
	//注册命名空间
	beego.AddNamespace(v1)
}
