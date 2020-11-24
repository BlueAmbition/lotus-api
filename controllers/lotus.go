package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"lotus-api/functions/array"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	cache2 "github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
)

type LotusController struct {
	BaseController
}

//安全模式,测试时可关闭
var SafeMode bool
var CoinId int

//预处理
func (c *LotusController) Prepare() {
	ip := c.Ctx.Input.IP()
	//fmt.Println(ip)
	allowIps := coin.AllowIps()
	if len(allowIps) < 1 {
		c.Data["json"] = c.GeneralReturn(FailureCode, "未设置节点白名单，禁止访问", nil)
		c.ServeJSON()
		return
	}
	if !array.InArray(ip, allowIps) {
		c.Data["json"] = c.GeneralReturn(FailureCode, "非白名单内IP，禁止访问", nil)
		c.ServeJSON()
		return
	}
	//安全模式证书加密解密
	safeModel, err := beego.AppConfig.Bool("usdt_server::safe_mode")
	if err != nil {
		safeModel = true
	}
	coinId, err := beego.AppConfig.Int("usdt_server::coin_id")
	SafeMode = safeModel
	CoinId = coinId
}

//总余额
func (c *LotusController) TotalBalance() {
	var (
		byteData  []byte
		err       error
		objStruct struct {
			IsValid bool `json:"is_valid"`
		}
		data string
	)
	data = c.GetString("data")
	byteData = []byte(data)
	if SafeMode {
		byteData, err = c.DecryptData(data)
		if err != nil {
			logs.Error("获取余额失败：", err.Error())
			c.Data["json"] = c.GeneralReturn(FailureCode, "解密失败："+err.Error(), nil)
			c.ServeJSON()
			return
		}
	}
	err = json.Unmarshal(byteData, &objStruct)
	if err != nil {
		c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "传输数据有误", nil)
		c.ServeJSON()
		return
	}
	client := c.TetherClient()
	if client != nil {
		//归集账户
		account := beego.AppConfig.String("usdt_server::union_account")
		contract := beego.AppConfig.String("usdt_server::contract")
		token, err := contract_token.NewTetherToken(common.HexToAddress(contract), client)
		balance, err := token.BalanceOf(nil, common.HexToAddress(account))
		if err != nil {
			logs.Error("获取余额失败：", err.Error())
			c.Data["json"] = c.GeneralReturn(FailureCode, "获取余额失败", nil)
			c.ServeJSON()
			return
		}
		ratio := math.Pow(10, 6)
		amount, err := strconv.ParseFloat(balance.String(), 64)
		realAmount := amount / ratio
		mapData := map[string]interface{}{
			"balance": realAmount}
		c.Data["json"] = c.GeneralReturn(SuccessCode, "获取余额成功", mapData)
		c.ServeJSON()
		return
	}
	c.Data["json"] = c.GeneralReturn(FailureCode, "获取余额失败", nil)
	c.ServeJSON()
}

//账户余额
func (c *LotusController) AccountBalance() {
	var (
		byteData  []byte
		err       error
		objStruct struct {
			Account string `json:"account"`
		}
		data string
	)
	data = c.GetString("data")
	byteData = []byte(data)
	if SafeMode {
		byteData, err = c.DecryptData(data)
		if err != nil {
			logs.Error("获取账户余额失败：", err.Error())
			c.Data["json"] = c.GeneralReturn(FailureCode, "解密失败："+err.Error(), nil)
			c.ServeJSON()
			return
		}
	}
	err = json.Unmarshal(byteData, &objStruct)
	if err != nil {
		c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "传输数据有误", nil)
		c.ServeJSON()
		return
	}
	client := c.TetherClient()
	if client != nil {
		contract := beego.AppConfig.String("usdt_server::contract")
		token, err := contract_token.NewTetherToken(common.HexToAddress(contract), client)
		balance, err := token.BalanceOf(nil, common.HexToAddress(objStruct.Account))
		if err != nil {
			logs.Error("获取余额失败：", err.Error())
			c.Data["json"] = c.GeneralReturn(FailureCode, "获取账户余额失败：RPC服务错误", nil)
			c.ServeJSON()
			return
		}
		ratio := math.Pow(10, 6)
		amount, err := strconv.ParseFloat(balance.String(), 64)
		realAmount := amount / ratio
		mapData := map[string]interface{}{
			"balance": realAmount}
		c.Data["json"] = c.GeneralReturn(SuccessCode, "获取账户余额成功", mapData)
		c.ServeJSON()
		return
	}
	c.Data["json"] = c.GeneralReturn(FailureCode, "获取账户余额失败", nil)
	c.ServeJSON()
}

//获取账户对应地址（创建新地址）
func (c *LotusController) GetReceiveAddress() {
	var (
		byteData  []byte
		err       error
		objStruct struct {
			Account string `json:"account"`
		}
		data string
	)
	data = c.GetString("data")
	byteData = []byte(data)
	if SafeMode {
		byteData, err = c.DecryptData(data)
		if err != nil {
			logs.Error("获取地址失败：", err.Error())
			c.Data["json"] = c.GeneralReturn(FailureCode, "解密失败："+err.Error(), nil)
			c.ServeJSON()
			return
		}
	}
	err = json.Unmarshal(byteData, &objStruct)
	if err != nil {
		c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "传输数据有误", nil)
		c.ServeJSON()
		return
	}
	if strings.Trim(objStruct.Account, " ") == "" {
		logs.Error("获取地址失败：", err.Error())
		c.Data["json"] = c.GeneralReturn(FailureCode, "获取地址失败", nil)
		c.ServeJSON()
		return
	}
	address := c.NewAccount(objStruct.Account)
	if address != "" {
		//把账户当做密码
		msg := fmt.Sprintf("账户：%v，地址：%v，密码：%v", objStruct.Account, address, objStruct.Account)
		logs.Info("获取地址信息：", msg)
		mapData := map[string]interface{}{
			"account": objStruct.Account, "address": address, "password": objStruct.Account}
		c.Data["json"] = c.GeneralReturn(SuccessCode, "获取地址成功", mapData)
		c.ServeJSON()
		return
	}
	c.Data["json"] = c.GeneralReturn(FailureCode, "获取地址失败", nil)
	c.ServeJSON()
}

//提币
func (c *LotusController) SendToAddress() {
	var (
		byteData  []byte
		err       error
		limit     bool
		objStruct struct {
			WithdrawId   int     `json:"withdraw_id"`
			Address      string  `json:"address"`
			OriginAmount float64 `json:"origin_amount"`
			Amount       float64 `json:"amount"`
			Auto         bool    `json:"auto"`
		}
		data string
	)
	data = c.GetString("data")
	byteData = []byte(data)
	if SafeMode {
		byteData, err = c.DecryptData(data)
		if err != nil {
			logs.Error("转账失败：", err.Error())
			c.Data["json"] = c.GeneralReturn(FailureCode, "解密失败："+err.Error(), nil)
			c.ServeJSON()
			return
		}
	}
	err = json.Unmarshal(byteData, &objStruct)
	if err != nil {
		c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "传输数据有误", nil)
		c.ServeJSON()
		return
	}
	if !objStruct.Auto {
		limit = frequencyLimit("send_coin", 5)
		if limit {
			c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "请求频率过高", nil)
			c.ServeJSON()
			return
		}
	}
	if strings.TrimSpace(objStruct.Address) == "" {
		c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "提币地址不能为空", nil)
		c.ServeJSON()
		return
	}
	if objStruct.Amount <= 0 {
		c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "提币数额必须大于0", nil)
		c.ServeJSON()
		return
	}
	if objStruct.OriginAmount < objStruct.Amount {
		c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "订单有误，原始金额不能小于实际金额", nil)
		c.ServeJSON()
		return
	}
	if objStruct.WithdrawId < 1 {
		c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "提币订单ID有误", nil)
		c.ServeJSON()
		return
	}
	//提币配置
	configMap := coin.CoinConfig(CoinId)
	if len(configMap) == 0 {
		c.Data["json"] = c.GeneralReturn(FailureCode, "币种未配置提币限制信息，禁止提币", nil)
		c.ServeJSON()
		return
	}
	min := configMap["min_withdraw"]
	max := configMap["max_withdraw"]
	if objStruct.OriginAmount < min || objStruct.OriginAmount > max {
		msg := fmt.Sprintf("提币量应在%v-%v之间", min, max)
		c.Data["json"] = c.GeneralReturn(FailureCode, msg, nil)
		c.ServeJSON()
		return
	}

	//内部地址过滤
	//addressList := coin.GetWalletAddresses(CoinId)
	//for _, v := range addressList {
	//	address := v["address"].(string)
	//	if strings.EqualFold(address, objStruct.Address) {
	//		c.Data["json"] = c.GeneralReturn(FailureCode, "无效提币：提币地址不能为节点上地址", nil)
	//		c.ServeJSON()
	//		return
	//	}
	//}
	//判断是否已经提币
	order := coin.GetWithdrawInfo(objStruct.WithdrawId)
	if order != nil {
		mapData := map[string]interface{}{
			"tx_id": order["tx_id"].(string)}
		c.Data["json"] = c.GeneralReturn(SuccessCode, "此订单已提币，不能重复提币", mapData)
		c.ServeJSON()
		return
	}
	if objStruct.Address[0:2] != "0x" {
		c.Data["json"] = c.GeneralReturn(FailureCode, "错误的提币地址", nil)
		c.ServeJSON()
		return
	}
	if len(objStruct.Address) != 42 {
		c.Data["json"] = c.GeneralReturn(FailureCode, "提币地址应为42位", nil)
		c.ServeJSON()
		return
	}
	//提币订单币种确认
	client := c.TetherClient()
	if client != nil {
		//要转账的USDT数量基本单位
		baseU := objStruct.Amount * math.Pow(10, float64(6))
		numStr := strconv.FormatFloat(baseU, 'f', 0, 64)
		transAmount := new(big.Int)
		transAmount, ok := transAmount.SetString(numStr, 10)
		//fmt.Println(transAmount)
		if ok {
			keystorePath := beego.AppConfig.String("usdt_server::keystore_path")
			account := beego.AppConfig.String("usdt_server::union_account")
			contract := beego.AppConfig.String("usdt_server::contract")
			password := beego.AppConfig.String("usdt_server::password")
			privateKey := eth.KeyStoreJson(keystorePath, account, password)
			objStruct.Address = strings.TrimSpace(objStruct.Address)
			if strings.EqualFold(account, objStruct.Address) {
				c.Data["json"] = c.GeneralReturn(FailureCode, "提币地址不合法", nil)
				c.ServeJSON()
				return
			}
			toAddress := common.HexToAddress(objStruct.Address)
			auth, err := bind.NewTransactor(strings.NewReader(privateKey), password)
			if err != nil {
				logs.Error("[Auth] 私钥Json解析失败，地址：" + account + " 密码：" + password)
				c.Data["json"] = c.GeneralReturn(FailureCode, "认证失败", nil)
				c.ServeJSON()
				return
			}
			token, err := contract_token.NewTetherToken(common.HexToAddress(contract), client)
			//设置成预估的限制
			auth.GasLimit = 0
			tx, err := token.Transfer(auth, toAddress, transAmount)
			if err != nil {
				c.Data["json"] = c.GeneralReturn(FailureCode, "转账失败："+err.Error(), nil)
				c.ServeJSON()
				return
			}
			txId := tx.Hash().String()
			msg := fmt.Sprintf("WithdrawId：%v，TxId：%v，原始金额：%v，金额：%v，转账地址：%v，提币地址：%v", objStruct.WithdrawId, txId, objStruct.OriginAmount, objStruct.Amount, account, objStruct.Address)
			logs.Info("提币信息：", msg)
			result := coin.WithdrawInfo(objStruct.WithdrawId, CoinId, objStruct.Amount, txId, msg)
			if result < 1 {
				logs.Error("提币记录写入失败：", msg)
			}
			mapData := map[string]interface{}{
				"tx_id": txId}
			c.Data["json"] = c.GeneralReturn(SuccessCode, "提币成功，请耐心等待......", mapData)
			c.ServeJSON()
			return
		}

	}
	c.Data["json"] = c.GeneralReturn(FailureCode, "ETH RPC服务连接有误", nil)
	c.ServeJSON()
}

//人工入账
func (c *LotusController) Recharge() {
	var (
		byteData  []byte
		err       error
		objStruct struct {
			TxId string `json:"tx_id"`
		}
		data   string
		result int64
		limit  bool
	)
	limit = frequencyLimit("manual_recharge", 5)
	if limit {
		c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "请求频率过高", nil)
		c.ServeJSON()
		return
	}
	data = c.GetString("data")
	byteData = []byte(data)
	if SafeMode {
		byteData, err = c.DecryptData(data)
		if err != nil {
			logs.Error("入账数据解析失败：", err.Error())
			c.Data["json"] = c.GeneralReturn(FailureCode, "解密失败："+err.Error(), nil)
			c.ServeJSON()
			return
		}
	}
	err = json.Unmarshal(byteData, &objStruct)
	if err != nil {
		c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "传输数据有误", nil)
		c.ServeJSON()
		return
	}
	client := c.TetherClient()
	ctx := context.Background()
	objStruct.TxId = strings.TrimSpace(objStruct.TxId)
	realHash := common.HexToHash(objStruct.TxId)
	trans, isPending, err := client.TransactionByHash(ctx, realHash)
	if err != nil {
		logs.Error("入账获取交易信息失败：", "txid："+objStruct.TxId+"错误信息："+err.Error())
		c.Data["json"] = c.GeneralReturn(FailureCode, "获取交易信息失败", nil)
		c.ServeJSON()
		return
	}
	if isPending {
		c.Data["json"] = c.GeneralReturn(FailureCode, "交易还在等待中", nil)
		c.ServeJSON()
		return
	}
	//归集账户
	unionAccount := beego.AppConfig.String("usdt_server::union_account")
	contract := beego.AppConfig.String("usdt_server::contract")
	transTo := trans.To()
	if transTo == nil {
		c.Data["json"] = c.GeneralReturn(FailureCode, "此交易为归集交易，入账失败", nil)
		c.ServeJSON()
		return
	}
	toContract := trans.To().String()
	if !strings.EqualFold(toContract, contract) {
		fmt.Println("大小写 比较错误合约：", toContract, "配置合约：", contract)
		c.Data["json"] = c.GeneralReturn(FailureCode, "此交易不属于TE Tether合约", nil)
		c.ServeJSON()
		return
	}
	//解析
	inputData := hexutil.Bytes(trans.Data())
	//fmt.Println(inputData)
	inputStr := inputData.String()
	//fmt.Println(inputStr)
	unpackData, err := eth.UnpackInput(inputStr)
	if err == nil {
		transAmount := unpackData.Value.String()
		toAddress := unpackData.To.String()
		amount, err := strconv.ParseFloat(transAmount, 64)
		if err != nil || amount <= 0 {
			c.Data["json"] = c.GeneralReturn(FailureCode, "金额有误", nil)
			c.ServeJSON()
			return
		}
		if strings.EqualFold(toAddress, unionAccount) {
			c.Data["json"] = c.GeneralReturn(FailureCode, "入账失败：此交易为归集交易", nil)
			c.ServeJSON()
			return
		}
		ratio := math.Pow(10, 6)
		realAmount := amount / ratio
		receiptInfo, err := client.TransactionReceipt(context.Background(), realHash)
		if err == nil && receiptInfo.Status == 1 {
			result = coin.Recharge(CoinId, realAmount, toAddress, objStruct.TxId)
			if result == 1 {
				c.Data["json"] = c.GeneralReturn(SuccessCode, "入账成功", nil)
				c.ServeJSON()
				return
			}
		}
	}

	msg := "入账失败：不存在钱包信息或此交易已入账"
	if result == -1 {
		msg = "此交易地址不存在钱包信息"
	} else if result == -2 {
		msg = "此交易已入账"
	}

	c.Data["json"] = c.GeneralReturn(FailureCode, msg, nil)
	c.ServeJSON()
}

//频率限制
func frequencyLimit(tag string, limitSeconds int64) bool {
	cache, err := cache2.NewCache("file", `{"CachePath":"./cache","FileSuffix":".cache","DirectoryLevel":"2","EmbedExpiry":"0"}`)
	if err == nil {
		currentTime := time.Now()
		ts := currentTime.Unix() //当前时间戳
		tagTime := cache.Get(tag)
		if tagTime == "" {
			//设置时间戳
			cache.Put(tag, ts, 0)
			return false
		}
		oldTime := tagTime.(int64)
		if ts-oldTime > limitSeconds {
			//超过限制时间，合法，重置时间戳
			cache.Put(tag, ts, 0)
			return false
		}
		return true
	}
	return true
}

//覆盖转账
func (c *LotusController) RecoverTx() {
	var (
		byteData  []byte
		err       error
		objStruct struct {
			TxId string `json:"tx_id"`
		}
		chainId  *big.Int
		gasPrice *big.Int
		transMsg types.Message
		data     string
		limit    bool
	)
	limit = frequencyLimit("manual_recover", 1)
	if limit {
		c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "请求频率过高", nil)
		c.ServeJSON()
		return
	}
	data = c.GetString("data")
	byteData = []byte(data)
	if SafeMode {
		byteData, err = c.DecryptData(data)
		if err != nil {
			logs.Error("覆盖交易数据解析失败：", err.Error())
			c.Data["json"] = c.GeneralReturn(FailureCode, "解密失败："+err.Error(), nil)
			c.ServeJSON()
			return
		}
	}
	err = json.Unmarshal(byteData, &objStruct)
	if err != nil {
		c.Data["json"] = c.GeneralReturn(UnprocessableEntityCode, "传输数据有误", nil)
		c.ServeJSON()
		return
	}
	client := c.TetherClient()
	ctx := context.Background()
	objStruct.TxId = strings.TrimSpace(objStruct.TxId)
	realHash := common.HexToHash(objStruct.TxId)
	trans, isPending, err := client.TransactionByHash(ctx, realHash)
	if err != nil {
		logs.Error("入账获取交易信息失败：", "txid："+objStruct.TxId+"错误信息："+err.Error())
		c.Data["json"] = c.GeneralReturn(FailureCode, "获取交易信息失败", nil)
		c.ServeJSON()
		return
	}
	if isPending {
		//归集账户
		unionAccount := beego.AppConfig.String("usdt_server::union_account")
		//合约判断
		contract := beego.AppConfig.String("usdt_server::contract")
		keystorePath := beego.AppConfig.String("usdt_server::keystore_path")
		account := beego.AppConfig.String("usdt_server::union_account")
		password := beego.AppConfig.String("usdt_server::password")
		privateKey := eth.KeyStoreJson(keystorePath, account, password)
		toContract := trans.To().String()
		if !strings.EqualFold(toContract, contract) {
			c.Data["json"] = c.GeneralReturn(FailureCode, "此交易不属于Tether合约", nil)
			c.ServeJSON()
			return
		}
		chainId, err = client.NetworkID(context.Background())
		if err != nil {
			c.Data["json"] = c.GeneralReturn(FailureCode, "获取Chain失败", nil)
			c.ServeJSON()
			return
		}
		transMsg, err = trans.AsMessage(types.NewEIP155Signer(chainId))
		if err != nil {
			c.Data["json"] = c.GeneralReturn(FailureCode, "获取交易信息失败", nil)
			c.ServeJSON()
			return
		}
		//来源地址
		fromAddress := transMsg.From().String()
		if !strings.EqualFold(unionAccount, fromAddress) {
			c.Data["json"] = c.GeneralReturn(FailureCode, "此交易来源非归集地址", nil)
			c.ServeJSON()
			return
		}
		//解析
		inputData := hexutil.Bytes(trans.Data())
		inputStr := inputData.String()
		unpackData, err := eth.UnpackInput(inputStr)
		if err != nil {
			c.Data["json"] = c.GeneralReturn(FailureCode, "解析数据失败", nil)
			c.ServeJSON()
			return
		}
		//转账数量
		transAmount := unpackData.Value
		//收款地址
		toAddress := unpackData.To
		originGasPrice := trans.GasPrice()
		originGas := trans.Gas()
		nonce := trans.Nonce()
		fmt.Println("转出地址：", fromAddress, "转账数量：", transAmount, "接收地址：", toAddress.String(), "原始GasPrice：", originGasPrice, "使用Gas", originGas, "Nonce：", nonce)
		//重新转账替换
		auth, err := bind.NewTransactor(strings.NewReader(privateKey), password)
		if err != nil {
			logs.Error("[Auth] 私钥Json解析失败，地址：" + account + " 密码：" + password)
			c.Data["json"] = c.GeneralReturn(FailureCode, "认证失败", nil)
			c.ServeJSON()
			return
		}
		//代币
		token, err := contract_token.NewTetherToken(common.HexToAddress(contract), client)
		if err != nil {
			c.Data["json"] = c.GeneralReturn(FailureCode, "Token有误："+err.Error(), nil)
			c.ServeJSON()
			return
		}
		//gas限制 0估算值
		gasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			c.Data["json"] = c.GeneralReturn(FailureCode, "GasPrice有误："+err.Error(), nil)
			c.ServeJSON()
			return
		}
		//覆盖的最少GasPrice建议原来GasPrice加10%
		minGasPrice := originGasPrice.Mul(originGasPrice, big.NewInt(11)).Div(originGasPrice, big.NewInt(10))
		if gasPrice.Cmp(minGasPrice) == -1 {
			gasPrice = minGasPrice
		}
		//取消操作
		//auth.Value = big.NewInt(0)
		//transAmount = big.NewInt(0)
		auth.GasLimit = 210000
		auth.GasPrice = gasPrice //originGasPrice.Add(originGasPrice, addGas)
		auth.Nonce = new(big.Int).SetInt64(int64(nonce))
		fmt.Println("新的GasPrice：", auth.GasPrice)
		tx, err := token.Transfer(auth, toAddress, transAmount)
		if err != nil {
			c.Data["json"] = c.GeneralReturn(FailureCode, "转账失败："+err.Error(), nil)
			c.ServeJSON()
			return
		}
		txId := tx.Hash().String()
		mapData := map[string]interface{}{
			"tx_id": txId}
		c.Data["json"] = c.GeneralReturn(SuccessCode, "覆盖交易成功", mapData)
		c.ServeJSON()
		return
	}

	c.Data["json"] = c.GeneralReturn(FailureCode, "非pending状态交易不能覆盖", nil)
	c.ServeJSON()
}
