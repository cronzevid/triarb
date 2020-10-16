package hitbtc

import ("encoding/json"
        "net/http"
        "net/url"
        "io/ioutil"
        "os"
        //"time"
        "strconv"
        "strings"
        "fmt")

const api = "api_key"
const secret = "sec_key"

type Symbol struct {
  Id string
  BaseCurrency string
  QuoteCurrency string
  FeeCurrency string
  TakeLiquidityRate float64 `json:",string"`
  ProvideLiquidityRate float64 `json:",string"`
  QuantityIncrement float64  `json:",string"`
  TickSize float64 `json:",string"`
}
type Balance struct {
  Currency string
  Available float64 `json:",string"`
  Reserved float64  `json:",string"`
}
type OrderId struct {
  Id int
  ClientOrderId string
  Symbol string
  Side string
  Status string
  Type string
  TimeInForce string
  Quantity float64 `json:",string"`
  Price float64 `json:",string"`
  CumQuantity float64 `json:",string"`
  PostOnly bool `json:",bool"`
  CreatedAt string
  UpdatedAt string
}
type ApiError struct {
  Error ErrorText
}
type ErrorText struct {
  Code int
  Message string
  Description string
}

func DataRequest(url string) []byte {
    //defer timeTrack(time.Now(), "api call " + url)

    resp, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    return body
}

func BalanceRequest(asset string) Balance {
    //defer timeTrack(time.Now(), "api balance call")
    var resAsset Balance

    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://api.hitbtc.com/api/2/trading/balance", nil)
    req.SetBasicAuth(api, secret)
    resp, err := client.Do(req)
    if err != nil{
        fmt.Println(err)
        os.Exit(1)
    }
    body, err := ioutil.ReadAll(resp.Body)

    var balanceJSON []Balance
    json.Unmarshal(body, &balanceJSON)

    //if strings.Contains(string(balanceJSON), "error") {
    //    fmt.Println(balanceJSON)
    //    os.Exit(1)
    //}

    for _, element := range balanceJSON {
        if element.Currency == asset {
            resAsset = element
            break
        }
    }

    return resAsset
}

func OrderHistory(id string) OrderId {
    //defer timeTrack(time.Now(), "api order history call")

    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://api.hitbtc.com/api/2/history/order", nil)
    req.SetBasicAuth(api, secret)
    resp, err := client.Do(req)
    if err != nil{
        fmt.Println(err)
        os.Exit(1)
    }
    body, err := ioutil.ReadAll(resp.Body)

    //fmt.Println(string(body))

    var res OrderId
    var history []OrderId
    json.Unmarshal(body, &history)

    for _, element := range history {
        if element.ClientOrderId == id {
            res = element
            break
        }
    }

    //if strings.Contains(string(curOrder), "error") {
    //    fmt.Println(curOrder)
    //    os.Exit(1)
    //}

    return res
}

func OrderStatus(id string) OrderId {
    //defer timeTrack(time.Now(), "api order status call")

    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://api.hitbtc.com/api/2/order/" + id, nil)
    req.SetBasicAuth(api, secret)
    resp, err := client.Do(req)
    if err != nil{
        fmt.Println(err)
        os.Exit(1)
    }
    body, err := ioutil.ReadAll(resp.Body)

    //fmt.Println(string(body))

    var curOrder OrderId
    json.Unmarshal(body, &curOrder)

    //if strings.Contains(string(curOrder), "error") {
    //    fmt.Println(curOrder)
    //    os.Exit(1)
    //}

    return curOrder
}

func Order(symbol string, side string, quantity float64, price float64) OrderId {
    //defer timeTrack(time.Now(), "api order call")
    quantity_conv := strconv.FormatFloat(quantity, 'f', -1, 64)
    price_conv := strconv.FormatFloat(price, 'f', -1, 64)

    //fmt.Println(quantity_conv, price_conv)

    params := url.Values{
	"symbol": {symbol},
	"side": {side},
	//"timeInForce": "FOK",
	"quantity": {quantity_conv},
	"price": {price_conv},
    }

    client := &http.Client{}
    req, err := http.NewRequest("POST", "https://api.hitbtc.com/api/2/order", strings.NewReader(params.Encode()))
    req.SetBasicAuth(api, secret)

    resp, err := client.Do(req)
    if err != nil{
        fmt.Println(err)
        os.Exit(1)
    }
    body, err := ioutil.ReadAll(resp.Body)

    //fmt.Println(string(body))

    var Order OrderId
    json.Unmarshal(body, &Order)

    //if strings.Contains(string(Order), "error") {
    //    fmt.Println(Order)
    //    os.Exit(1)
    //}

    return Order
}

func CancelOrder(id string) {
    //defer timeTrack(time.Now(), "api cancel order call")

    client := &http.Client{}
    req, _ := http.NewRequest("DELETE", "https://api.hitbtc.com/api/2/order/" + id, nil)
    req.SetBasicAuth(api, secret)

    client.Do(req)
    //resp, err := client.Do(req)
    //if err != nil{
    //    fmt.Println(err)
    //    os.Exit(1)
    //}
    //body, err := ioutil.ReadAll(resp.Body)

    //fmt.Println(string(body))

    //if strings.Contains(string(Order), "error") {
    //    fmt.Println(Order)
    //    os.Exit(1)
    //}

}
