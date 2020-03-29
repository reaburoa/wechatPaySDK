package request

type WechatPayMicropay struct {
    SubMchId       string `json:"sub_mch_id"`
    Body           string `json:"body"`
    OutTradeNo     string `json:"out_trade_no"`
    TotalFee       string `json:"total_fee"`
    AuthCode       string `json:"auth_code"`
    SpbillCreateIp string `json:"spbill_create_ip"`
}

func (w *WechatPayMicropay) GetApiUrl() string {
    return "https://api.mch.weixin.qq.com/pay/micropay"
}

func (w *WechatPayMicropay) SetNotifyUrl(str string) {

}

func (w *WechatPayMicropay) GetNotifyUrl() string {
    return ""
}

func (w *WechatPayMicropay) GetSignType() string {
    return "md5"
}

func (w *WechatPayMicropay) GetClientIp() string {
    return "127.0.0.1"
}

func (w *WechatPayMicropay) GetParams() Requester {
    return w
}
