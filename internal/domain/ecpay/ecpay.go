package ecpay

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

// PaymentData 是一个包含所有支付请求参数的结构体
type PaymentData struct {
	MerchantID        string `json:"MerchantID"`
	MerchantTradeNo   string `json:"MerchantTradeNo"`
	MerchantTradeDate string `json:"MerchantTradeDate"`
	PaymentType       string `json:"PaymentType"`
	TotalAmount       string `json:"TotalAmount"`
	TradeDesc         string `json:"TradeDesc"`
	ItemName          string `json:"ItemName"`
	ReturnURL         string `json:"ReturnURL"`
	ChoosePayment     string `json:"ChoosePayment"`
	CheckMacValue     string `json:"CheckMacValue"`
}

// generateCheckMacValue 生成 CheckMacValue 的函数
func GenerateCheckMacValue(data PaymentData, hashKey, hashIV string) string {
	// 将 PaymentData 转换成 map
	params := map[string]string{
		"MerchantID":        data.MerchantID,
		"MerchantTradeNo":   data.MerchantTradeNo,
		"MerchantTradeDate": data.MerchantTradeDate,
		"PaymentType":       data.PaymentType,
		"TotalAmount":       data.TotalAmount,
		"TradeDesc":         data.TradeDesc,
		"ItemName":          data.ItemName,
		"ReturnURL":         data.ReturnURL,
		"ChoosePayment":     data.ChoosePayment,
	}

	// 按键名排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接字符串
	var sb strings.Builder
	sb.WriteString("HashKey=" + hashKey + "&")
	for _, k := range keys {
		sb.WriteString(k + "=" + url.QueryEscape(params[k]) + "&")
	}
	sb.WriteString("HashIV=" + hashIV)

	// URL编码并转换成小写字母
	encodedStr := strings.ToLower(url.QueryEscape(sb.String()))

	// 计算 MD5 哈希值
	hash := md5.New()
	hash.Write([]byte(encodedStr))
	hashValue := hash.Sum(nil)

	// 将哈希值转换成十六进制字符串
	checkMacValue := hex.EncodeToString(hashValue)
	return checkMacValue
}
