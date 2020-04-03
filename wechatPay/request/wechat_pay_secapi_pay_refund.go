package request

type WechatPaySecapiPayRefund struct {
    SubAppId      string `json:"sub_appid"`
    SubMchId      string `json:"sub_mch_id"`
    SignType      string `json:"sign_type"`
    TransactionId string `json:"transaction_id"`
    OutTradeNo    string `json:"out_trade_no"`
    OutRefundNo   string `json:"out_refund_no"`
    TotalFee      string `json:"total_fee"`
    RefundFee     string `json:"refund_fee"`
    RefundFeeType string `json:"refund_fee_type"`
    RefundDesc    string `json:"refund_desc"`
    RefundAccount string `json:"refund_account"`
    NotifyUrl     string `json:"notify_url"`
}

func (w *WechatPaySecapiPayRefund) GetApiUrl() string {
    return "https://api.mch.weixin.qq.com/secapi/pay/refund"
}

func (w *WechatPaySecapiPayRefund) SetNotifyUrl(url string) {
    w.NotifyUrl = url
}

func (w *WechatPaySecapiPayRefund) GetNotifyUrl() string {
    return w.NotifyUrl
}

func (w *WechatPaySecapiPayRefund) SetClientIp(str string) {

}

func (w *WechatPaySecapiPayRefund) GetClientIp() string {
    return ""
}

func (w *WechatPaySecapiPayRefund) GetParams() Requester {
    return w
}

func (w *WechatPaySecapiPayRefund) GetRequestDataType() string {
    return RequestDataXML
}

func (w *WechatPaySecapiPayRefund) GetSignType() string {
    return w.SignType
}

func (w *WechatPaySecapiPayRefund) SetSignType(signType string) {
    w.SignType = signType
}
