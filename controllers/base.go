package controllers

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"lotus-api/functions/str"
)

type BaseController struct {
	beego.Controller
}

//常量
const (
	//成功code
	SuccessCode = 200
	//失败code
	FailureCode = 500
	//数据校验失败
	UnprocessableEntityCode = 422
)

//通用返回值
func (c *BaseController) GeneralReturn(code int, msg string, data interface{}) map[string]interface{} {
	mapData := map[string]interface{}{
		"code": code, "status": "failure", "msg": msg}
	if code == SuccessCode {
		mapData["status"] = "success"
	}
	//数据节点处理
	if data != nil {
		mapData["data"] = data
	}

	return mapData
}

//连接tet ERC20服务器
//func (c *BaseController) TetherClient() *ethclient.Client {
//	host := beego.AppConfig.String("usdt_server::host")
//	port, _ := beego.AppConfig.Int("usdt_server::port")
//	client, err := ethclient.Dial("http://" + host + ":" + strconv.Itoa(port))
//	if err != nil {
//		return nil
//	}
//	return client
//}

//解密数据
func (c *BaseController) DecryptData(encryptData string) ([]byte, error) {
	privateKey, err := ioutil.ReadFile("private.pem")
	if err != nil {
		return nil, err
	}
	decryptStr, err := str.RsaDecrypt(privateKey, []byte(encryptData))
	return decryptStr, err
}
