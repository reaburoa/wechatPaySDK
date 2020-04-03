package request

type WechatPayRefundQuery struct {
    SubAppId      string `json:"sub_appid"`
    SubMchId      string `json:"sub_mch_id"`
    TransactionId string `json:"transaction_id"`
    OutTradeNo    string `json:"out_trade_no"`
    OutRefundNo   string `json:"out_refund_no"`
    RefundId      string `json:"refund_id"`
    Offset        int    `json:"offset"`
}

func (w *WechatPayRefundQuery) GetApiUrl() string {
    return "https://api.mch.weixin.qq.com/pay/refundquery"
}

func (w *WechatPayRefundQuery) SetNotifyUrl(url string) {

}

func (w *WechatPayRefundQuery) GetNotifyUrl() string {
    return ""
}

func (w *WechatPayRefundQuery) SetClientIp(str string) {

}

func (w *WechatPayRefundQuery) GetClientIp() string {
    return ""
}

func (w *WechatPayRefundQuery) GetParams() Requester {
    return w
}

func (w *WechatPayRefundQuery) GetRequestDataType() string {
    return RequestDataXML
}
