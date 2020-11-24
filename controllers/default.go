package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"lotus-api/functions/req"
	"lotus-api/functions/str"
	"os"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	var
	(
		err       error
		publicKey []byte
		//privateKey []byte
		encryptStr []byte
		//decryptStr []byte
		jsonData []byte
		postData map[string][]string
		data     map[string]interface{}
		res      []byte
	)
	publicKey, err = ioutil.ReadFile("public.pem")
	if err != nil {
		os.Exit(-1)
	}
	//privateKey, err = ioutil.ReadFile("private.pem")
	//ERC20 TET总额
	data = map[string]interface{}{"is_valid": true}
	jsonData, err = json.Marshal(data)
	encryptStr, err = str.RsaEncrypt(publicKey, jsonData)
	postData = make(map[string][]string)
	postData["data"] = []string{string(encryptStr)}
	res, err = req.PostForm("http://127.0.0.1:8083/erc20_usdt/total-balance", postData)
	fmt.Println("返回数据:" + string(res))
	if err != nil {
		fmt.Println("错误信息:" + err.Error())
	}

	//账户余额
	data = map[string]interface{}{"account": "0x6FF8fCEC1FaF3623f277a7bc06f840bb602059CB"}
	jsonData, err = json.Marshal(data)
	encryptStr, err = str.RsaEncrypt(publicKey, jsonData)
	postData = make(map[string][]string)
	postData["data"] = []string{string(encryptStr)}
	res, err = req.PostForm("http://127.0.0.1:8083/erc20_usdt/account-balance", postData)
	fmt.Println("返回数据:" + string(res))
	if err != nil {
		fmt.Println("错误信息:" + err.Error())
	}

	//获取ERC20 TET地址
	//data = map[string]interface{}{"account": "12345678"}
	//jsonData, err = json.Marshal(data)
	//encryptStr, err = str.RsaEncrypt(publicKey, jsonData)
	//postData = make(map[string][]string)
	//postData["data"] = []string{string(encryptStr)}
	//res, err = req.PostForm("http://127.0.0.1:8083/erc20_usdt/receive-address", postData)
	//fmt.Println("返回数据:" + string(res))
	//if err != nil {
	//	fmt.Println("错误信息:" + err.Error())
	//}

	//ERC20 TET转账测试
	//data = map[string]interface{}{"address": "0x6FF8fCEC1FaF3623f277a7bc06f840bb602059CB", "origin_amount": 1, "amount": 1, "memo": "", "withdraw_id": 1010}
	//jsonData, err = json.Marshal(data)
	//encryptStr, err = str.RsaEncrypt(publicKey, jsonData)
	//postData = make(map[string][]string)
	//postData["data"] = []string{string(encryptStr)}
	//res, err = req.PostForm("http://127.0.0.1:8083/erc20_usdt/send-to-address", postData)
	//fmt.Println("返回数据:" + string(res))
	//if err != nil {
	//	fmt.Println("错误信息:" + err.Error())
	//}

	//人工入账
	data = map[string]interface{}{"tx_id": "0x7eb31e1910ecf8879f181f5e664159c17941c85e741950112325d785f9d19986"}
	//data = map[string]interface{}{"tx_id": "0xff64f5017c5062c7e49ba4618793ae9b27c278509095d54310c0556ab307d347"}
	jsonData, err = json.Marshal(data)
	encryptStr, err = str.RsaEncrypt(publicKey, jsonData)
	postData = make(map[string][]string)
	postData["data"] = []string{string(encryptStr)}
	res, err = req.PostForm("http://127.0.0.1:8083/erc20_usdt/recharge", postData)
	fmt.Println("返回数据:" + string(res))
	if err != nil {
		fmt.Println("错误信息:" + err.Error())
	}
}
