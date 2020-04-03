package request

type WechatPayMicropay struct {
    SubAppId       string        `json:"sub_appid"`
    SubMchId       string        `json:"sub_mch_id"`
    DeviceInfo     string        `json:"device_info"`
    Body           string        `json:"body"`
    Detail         string        `json:"detail"`
    Attach         string        `json:"attach"`
    OutTradeNo     string        `json:"out_trade_no"`
    TotalFee       string        `json:"total_fee"`
    FeeType        string        `json:"fee_type"`
    SpbillCreateIp string        `json:"spbill_create_ip"`
    GoodsTag       string        `json:"goods_tag"`
    LimitPay       string        `json:"limit_pay"`
    TimeStart      string        `json:"time_start"`
    TimeExpire     string        `json:"time_expire"`
    AuthCode       string        `json:"auth_code"`
    Receipt        string        `json:"receipt"`
    SceneInfo      SceneInfoData `json:"scene_info"`
}

type SceneInfoData struct {
    Id       string `json:"id"`
    Name     string `json:"name"`
    AreaCode string `json:"area_code"`
    Address  string `json:"address"`
}

func (w *WechatPayMicropay) GetApiUrl() string {
    return "https://api.mch.weixin.qq.com/pay/micropay"
}

func (w *WechatPayMicropay) SetNotifyUrl(str string) {

}

func (w *WechatPayMicropay) GetNotifyUrl() string {
    return ""
}

func (w *WechatPayMicropay) SetClientIp(str string) {
    w.SpbillCreateIp = str
}

func (w *WechatPayMicropay) GetClientIp() string {
    return w.SpbillCreateIp
}

func (w *WechatPayMicropay) GetParams() Requester {
    return w
}

func (w *WechatPayMicropay) GetRequestDataType() string {
    return RequestDataXML
}
