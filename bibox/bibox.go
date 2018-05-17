package bibox

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	. "github.com/marstau/GoEx"
	"net/http"
	"net/url"
)

type Bibox struct {
	httpClient *http.Client
	accountId,
	baseUrl,
	accessKey,
	secretKey string
}

type response struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data"`
	Errmsg  string          `json:"err-msg"`
	Errcode string          `json:"err-code"`
}

type StBody struct {
	Pair string  `json:"pair"`
}

type Cmds struct {
	Cmd string `json:"cmd"`
	Body StBody `json:"body"`
}

func NewBibox(httpClient *http.Client, accessKey, secretKey, clientId string) *Bibox {
	return &Bibox{httpClient, clientId, "https://api.bibox.com", accessKey, secretKey}
}

func (bb *Bibox) GetAccountId() (string, error) {
	// path := "/v1/account/accounts"
	// params := &url.Values{}
	// bb.buildPostForm("GET", path, params)

	// //log.Println(bb.baseUrl + path + "?" + params.Encode())

	// respmap, err := HttpGet(bb.httpClient, bb.baseUrl+path+"?"+params.Encode())
	// if err != nil {
	// 	return "", err
	// }
	// //log.Println(respmap)
	// if respmap["status"].(string) != "ok" {
	// 	return "", errors.New(respmap["err-code"].(string))
	// }

	// data := respmap["data"].([]interface{})
	// accountIdMap := data[0].(map[string]interface{})
	// bb.accountId = fmt.Sprintf("%.f", accountIdMap["id"].(float64))

	// //log.Println(respmap)
	// return bb.accountId, nil
	path := fmt.Sprintf("/v1/user")
	cmdlist := fmt.Sprintf("user/userInfo")
	emptyBody := StBody{}
	cmdsJ := Cmds{Cmd : cmdlist, Body : emptyBody }
	params := url.Values{}

	c, err := json.Marshal(cmdsJ)
	if err != nil {
		log.Println("error:", err)
		return "", err
	}

	log.Println("body=" + string(c) + ", "+ string(c[:]))
	bb.buildPostForm(&params, string(c))

	urlStr := bb.baseUrl + path

	log.Println("urlStr=",urlStr)

	bodyData, err := HttpPostForm(bb.httpClient, urlStr, params)
	if err != nil {
		log.Println("error:", err)
		return "", err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(bodyData, &bodyDataMap)
	if err != nil {
		println(string(bodyData))
		return "", err
	}

	log.Println(bodyDataMap)
	
	if bodyDataMap["code"] != nil {
		return "", errors.New(string(bodyData))
	}


	return "", nil
}

// [balance:0.11237468
// confirm_count:3
// forbid_info:
// supply_time:2009/1/3
// totalBalance:0.24292026
// id:42
// icon_url:/appimg/BTC_icon.png
// is_erc20:0
// BTCValue:0.24292026
// freeze:0.13054557
// enable_deposit:1
// describe_summary:[{"lang":"zh-cn","text":"Bitcoin 比特币的概念最初由中本聪在2009年提出，是点对点的基于 SHA-256 算法的一种P2P形式的数字货币，点对点的传输意味着一个去中心化的支付系统。"},{"lang":"en-ww","text":"Bitcoin is a digital asset and a payment system invented by Satoshi Nakamoto who published a related paper in 2008 and released it as open-source software in 2009. The system featured as peer-to-peer; users can transact directly without an intermediary."}]
// price:--
// describe_url:[{"lang":"zh-cn","link":"https://bibox.zendesk.com/hc/zh-cn/articles/115004798214"},{"lang":"en-ww","link":"https://bibox.zendesk.com/hc/en-us/articles/115004798214"}]
// name:Bitcoin
// enable_withdraw:1
// total_amount:2.1e+07
// supply_amount:1.6818187e+07
// CNYValue:12927.64793111
// USDValue:2027.44302504
// symbol:BTC]

func (bb *Bibox) GetAccount() (*Account, error) {
	log.Println("GetAccount")
	path := fmt.Sprintf("/v1/transfer")
	cmdlist := fmt.Sprintf("transfer/coinList")
	emptyBody := StBody{}
	cmdsJ := Cmds{Cmd : cmdlist, Body : emptyBody }
	params := url.Values{}

	c, err := json.Marshal(cmdsJ)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	log.Println("body=" + string(c) + ", "+ string(c[:]))
	bb.buildPostForm(&params, string(c))

	urlStr := bb.baseUrl + path

	log.Println("urlStr=",urlStr)

	bodyData, err := HttpPostForm(bb.httpClient, urlStr, params)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(bodyData, &bodyDataMap)
	if err != nil {
		println(string(bodyData))
		return nil, err
	}

	resultMap := bodyDataMap["result"].([]interface{})[0].(map[string]interface{})

	list := resultMap["result"].([]interface{})
	// log.Println(list)
	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount, 6)
	acc.Exchange = bb.GetExchangeName()

	subAccMap := make(map[Currency]*SubAccount)

	for _, v := range list {
		balancemap := v.(map[string]interface{})
		currencySymbol := balancemap["symbol"].(string)
		currency := NewCurrency(currencySymbol, "")
		if subAccMap[currency] == nil {
			subAccMap[currency] = new(SubAccount)
		}
		subAccMap[currency].Currency = currency
		subAccMap[currency].Amount = ToFloat64(balancemap["balance"])
		subAccMap[currency].ForzenAmount = ToFloat64(balancemap["freeze"])
	}

	for k, v := range subAccMap {
		acc.SubAccounts[k] = *v
	}

	return acc, nil
}

func (bb *Bibox) placeOrder(amount, price string, pair CurrencyPair, orderType string) (string, error) {
	// path := "/v1/order/orders/place"
	// params := url.Values{}
	// params.Set("account-id", bb.accountId)
	// params.Set("amount", amount)
	// params.Set("symbol", strings.ToLower(pair.ToSymbol("")))
	// params.Set("type", orderType)

	// switch orderType {
	// case "buy-limit", "sell-limit":
	// 	params.Set("price", price)
	// }

	// bb.buildPostForm("POST", path, &params)

	// resp, err := HttpPostForm3(bb.httpClient, bb.baseUrl+path+"?"+params.Encode(), bb.toJson(params),
	// 	map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})
	// if err != nil {
	// 	return "", err
	// }

	// respmap := make(map[string]interface{})
	// err = json.Unmarshal(resp, &respmap)
	// if err != nil {
	// 	return "", err
	// }

	// if respmap["status"].(string) != "ok" {
	// 	return "", errors.New(respmap["err-code"].(string))
	// }

	// return respmap["data"].(string), nil
	return "", nil
}

func (bb *Bibox) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := bb.placeOrder(amount, price, currency, "buy-limit")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		OrderID2: orderId,
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     BUY}, nil
}

func (bb *Bibox) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := bb.placeOrder(amount, price, currency, "sell-limit")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		OrderID2: orderId,
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     SELL}, nil
}

func (bb *Bibox) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := bb.placeOrder(amount, price, currency, "buy-market")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		OrderID2: orderId,
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     BUY_MARKET}, nil
}

func (bb *Bibox) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := bb.placeOrder(amount, price, currency, "sell-market")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		OrderID2: orderId,
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     SELL_MARKET}, nil
}

func (bb *Bibox) parseOrder(ordmap map[string]interface{}) Order {
	ord := Order{
		OrderID:    ToInt(ordmap["id"]),
		OrderID2:   fmt.Sprint(ToInt(ordmap["id"])),
		Amount:     ToFloat64(ordmap["amount"]),
		Price:      ToFloat64(ordmap["price"]),
		DealAmount: ToFloat64(ordmap["field-amount"]),
		Fee:        ToFloat64(ordmap["field-fees"]),
		OrderTime:  ToInt(ordmap["created-at"]),
	}

	state := ordmap["state"].(string)
	switch state {
	case "submitted","pre-submitted":
		ord.Status = ORDER_UNFINISH
	case "filled":
		ord.Status = ORDER_FINISH
	case "partial-filled":
		ord.Status = ORDER_PART_FINISH
	case "canceled", "partial-canceled":
		ord.Status = ORDER_CANCEL
	default:
		ord.Status = ORDER_UNFINISH
	}

	if ord.DealAmount > 0.0 {
		ord.AvgPrice = ToFloat64(ordmap["field-cash-amount"]) / ord.DealAmount
	}

	typeS := ordmap["type"].(string)
	switch typeS {
	case "buy-limit":
		ord.Side = BUY
	case "buy-market":
		ord.Side = BUY_MARKET
	case "sell-limit":
		ord.Side = SELL
	case "sell-market":
		ord.Side = SELL_MARKET
	}
	return ord
}

func (bb *Bibox) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	// path := "/v1/order/orders/" + orderId
	// params := url.Values{}
	// bb.buildPostForm("GET", path, &params)
	// respmap, err := HttpGet(bb.httpClient, bb.baseUrl+path+"?"+params.Encode())
	// if err != nil {
	// 	return nil, err
	// }

	// if respmap["status"].(string) != "ok" {
	// 	return nil, errors.New(respmap["err-code"].(string))
	// }

	// datamap := respmap["data"].(map[string]interface{})
	// order := bb.parseOrder(datamap)
	// order.Currency = currency
	// //log.Println(respmap)
	// return &order, nil
	return nil, nil
}

func (bb *Bibox) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	return bb.getOrders(queryOrdersParams{
		pair:   currency,
		states: "pre-submitted,submitted,partial-filled",
		size:   100,
		//direct:""
	})
}

func (bb *Bibox) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	// path := fmt.Sprintf("/v1/order/orders/%s/submitcancel", orderId)
	// params := url.Values{}
	// bb.buildPostForm("POST", path, &params)
	// resp, err := HttpPostForm3(bb.httpClient, bb.baseUrl+path+"?"+params.Encode(), bb.toJson(params),
	// 	map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})
	// if err != nil {
	// 	return false, err
	// }

	// var respmap map[string]interface{}
	// err = json.Unmarshal(resp, &respmap)
	// if err != nil {
	// 	return false, err
	// }

	// if respmap["status"].(string) != "ok" {
	// 	return false, errors.New(string(resp))
	// }

	return true, nil
}

func (bb *Bibox) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return bb.getOrders(queryOrdersParams{
		pair:   currency,
		size:   pageSize,
		states: "partial-canceled,filled",
		direct: "next",
	})
}

type queryOrdersParams struct {
	types,
	startDate,
	endDate,
	states,
	from,
	direct string
	size int
	pair CurrencyPair
}

func (bb *Bibox) getOrders(queryparams queryOrdersParams) ([]Order, error) {
	// path := "/v1/order/orders"
	// params := url.Values{}
	// params.Set("symbol", strings.ToLower(queryparams.pair.ToSymbol("")))
	// params.Set("states", queryparams.states)

	// if queryparams.direct != "" {
	// 	params.Set("direct", queryparams.direct)
	// }

	// if queryparams.size > 0 {
	// 	params.Set("size", fmt.Sprint(queryparams.size))
	// }

	// bb.buildPostForm("GET", path, &params)
	// respmap, err := HttpGet(bb.httpClient, fmt.Sprintf("%s%s?%s", bb.baseUrl, path, params.Encode()))
	// if err != nil {
	// 	return nil, err
	// }

	// if respmap["status"].(string) != "ok" {
	// 	return nil, errors.New(respmap["err-code"].(string))
	// }

	// datamap := respmap["data"].([]interface{})
	// var orders []Order
	// for _, v := range datamap {
	// 	ordmap := v.(map[string]interface{})
	// 	ord := bb.parseOrder(ordmap)
	// 	ord.Currency = queryparams.pair
	// 	orders = append(orders, ord)
	// }

	// return orders, nil
	return nil, nil
}

func (bb *Bibox) GetExchangeName() string {
	return "bibox.com"
}

// map[
// last:0.00000393
// last_usd:0.03
// last_cny:0.19
// high:0.00000423
// sell:0.00000397
// sell_amount:558.0848
// vol:1304725
// pair:MT_BTC
// timestamp:1.526549490466e+12
// percent:+1.29%
// buy:0.00000393
// buy_amount:3917.4615
// low:0.00000379
// ]
func (bb *Bibox) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	path := fmt.Sprintf("/v1/mdata")
	cmdlist := fmt.Sprintf("api/ticker")
	body := StBody{Pair : "MT_BTC"}
	cmdsJ := Cmds{Cmd : cmdlist, Body : body }
	params := url.Values{}

	c, err := json.Marshal(cmdsJ)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	log.Println("body=" + string(c) + ", "+ string(c[:]))
	bb.buildPostForm(&params, string(c))

	urlStr := bb.baseUrl + path

	log.Println("urlStr=",urlStr)

	bodyData, err := HttpPostForm(bb.httpClient, urlStr, params)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(bodyData, &bodyDataMap)
	if err != nil {
		println(string(bodyData))
		return nil, err
	}

	resultMap := bodyDataMap["result"].([]interface{})[0].(map[string]interface{})

	tickmap, ok := resultMap["result"].(map[string]interface{})
	if !ok {
		return nil, errors.New("tick assert error")
	}

	ticker := new(Ticker)
	ticker.Vol = ToFloat64(tickmap["vol"])
	ticker.Low = ToFloat64(tickmap["low"])
	ticker.High = ToFloat64(tickmap["high"])
	ticker.Buy = ToFloat64(tickmap["buy"])
	ticker.Sell = ToFloat64(tickmap["sell"])
	ticker.Last = ToFloat64(tickmap["last"])
	ticker.Date = ToUint64(respmap["timestamp"])

	return nil, ticker
}

func (bb *Bibox) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	path := fmt.Sprintf("/v1/mdata")
	cmdlist := fmt.Sprintf("api/market")
	emptyBody := StBody{}
	cmdsJ := Cmds{Cmd : cmdlist, Body : emptyBody }
	params := url.Values{}

	c, err := json.Marshal(cmdsJ)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	log.Println("body=" + string(c) + ", "+ string(c[:]))
	bb.buildPostForm(&params, string(c))

	urlStr := bb.baseUrl + path

	log.Println("urlStr=",urlStr)

	bodyData, err := HttpPostForm(bb.httpClient, urlStr, params)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(bodyData, &bodyDataMap)
	if err != nil {
		println(string(bodyData))
		return nil, err
	}

	resultMap := bodyDataMap["result"].([]interface{})[0].(map[string]interface{})

	list := resultMap["result"].([]interface{})
	// log.Println(list)
	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount, 6)
	acc.Exchange = bb.GetExchangeName()

	subAccMap := make(map[Currency]*SubAccount)

	for _, v := range list {
		balancemap := v.(map[string]interface{})
		currencySymbol := balancemap["symbol"].(string)
		currency := NewCurrency(currencySymbol, "")
		if subAccMap[currency] == nil {
			subAccMap[currency] = new(SubAccount)
		}
		subAccMap[currency].Currency = currency
		subAccMap[currency].Amount = ToFloat64(balancemap["balance"])
		subAccMap[currency].ForzenAmount = ToFloat64(balancemap["freeze"])
	}

	for k, v := range subAccMap {
		acc.SubAccounts[k] = *v
	}


	// url := bb.baseUrl + "/market/depth?symbol=%s&type=step0"

	// respmap, err := HttpGet(bb.httpClient, fmt.Sprintf(url, strings.ToLower(currency.ToSymbol(""))))
	// if err != nil {
	// 	return nil, err
	// }

	// if "ok" != respmap["status"].(string) {
	// 	return nil, errors.New(respmap["err-msg"].(string))
	// }

	// tick, _ := respmap["tick"].(map[string]interface{})
	// bids, _ := tick["bids"].([]interface{})
	// asks, _ := tick["asks"].([]interface{})

	// depth := new(Depth)
	// _size := size
	// for _, r := range asks {
	// 	var dr DepthRecord
	// 	rr := r.([]interface{})
	// 	dr.Price = ToFloat64(rr[0])
	// 	dr.Amount = ToFloat64(rr[1])
	// 	depth.AskList = append(depth.AskList, dr)

	// 	_size--
	// 	if _size == 0 {
	// 		break
	// 	}
	// }

	// _size = size
	// for _, r := range bids {
	// 	var dr DepthRecord
	// 	rr := r.([]interface{})
	// 	dr.Price = ToFloat64(rr[0])
	// 	dr.Amount = ToFloat64(rr[1])
	// 	depth.BidList = append(depth.BidList, dr)

	// 	_size--
	// 	if _size == 0 {
	// 		break
	// 	}
	// }

	return nil, nil
}

func (bb *Bibox) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (bb *Bibox) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (bb *Bibox) buildPostForm(postForm *url.Values,postCMDS string) error {
	log.Println("accessKey="+bb.accessKey)
	postForm.Set("apikey", bb.accessKey)
	postCMDS = "[" + postCMDS + "]"
	log.Println("cmds="+postCMDS)
	postForm.Set("cmds", postCMDS)
	sign, err := GetParamHmacMD5Sign(bb.secretKey, postCMDS)
	if err != nil {
		return err
	}
	log.Println("sign=",sign)
	postForm.Set("sign", sign)
	return nil
}

func (bb *Bibox) toJson(params url.Values) string {
	parammap := make(map[string]string)
	for k, v := range params {
		parammap[k] = v[0]
	}
	jsonData, _ := json.Marshal(parammap)
	return string(jsonData)
}
