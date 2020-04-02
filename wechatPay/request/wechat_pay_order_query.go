package request

type WechatPayOrderQuery struct {
    SubAppId      string `json:"sub_appid"`
    SubMchId      string `json:"sub_mch_id"`
    TransactionId string `json:"transaction_id"`
    OutTradeNo    string `json:"out_trade_no"`
}

func (w *WechatPayOrderQuery) GetApiUrl() string {
    return "https://api.mch.weixin.qq.com/pay/orderquery"
}

func (w *WechatPayOrderQuery) SetNotifyUrl(str string) {

}

func (w *WechatPayOrderQuery) GetNotifyUrl() string {
    return ""
}

func (w *WechatPayOrderQuery) GetSignType() string {
    return "MD5"
}

func (w *WechatPayOrderQuery) GetClientIp() string {
    return "127.0.0.1"
}

func (w *WechatPayOrderQuery) GetParams() Requester {
    return w
}

func (w *WechatPayOrderQuery) GetRequestDataType() string {
    return RequestDataXML
}
