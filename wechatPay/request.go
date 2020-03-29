package wechatPay

import (
    "encoding/json"
    "github.com/reaburoa/wechatPaySDK/wechatPay/request"
    "reflect"
)

type CommonRequest struct {
    AppId          string `json:"appid"`
    MchId          string `json:"mch_id"`
    NonceStr       string `json:"nonce_str"`
    Sign           string `json:"sign"`
    request.Requester
}

func (r *CommonRequest) toMap() map[string]interface{} {
    m := make(map[string]interface{})
    elemValues := reflect.ValueOf(r).Elem()
    elemTypes := elemValues.Type()
    for i := 0; i < elemTypes.NumField(); i++ {
        if elemValues.Field(i).Kind() != reflect.Interface {
            m[elemTypes.Field(i).Tag.Get("json")] = elemValues.Field(i).Interface()
        }
    }
    iParams := r.GetParams()
    b, e := json.Marshal(iParams)
    if e != nil {
        panic("Can Not Json Encode")
    }
    tmp := make(map[string]interface{})
    _ = json.Unmarshal(b, &tmp)
    for k, v := range tmp {
        m[k] = v
    }
    return m
}
