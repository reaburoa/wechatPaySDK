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
    iValues := reflect.ValueOf(iParams).Elem()
    iTypes := iValues.Type()
    for pI := 0; pI < iTypes.NumField(); pI++ {
        if iValues.Field(pI).IsZero() {
            continue
        }
        if iValues.Field(pI).Kind() == reflect.Struct {
            mm := make(map[string]interface{})
            mmType := iValues.Field(pI).Type()
            for si := 0; si < iValues.Field(pI).NumField(); si ++ {
                mm[mmType.Field(si).Tag.Get("json")] = iValues.Field(pI).Field(si).Interface()
            }
            b, e := json.Marshal(mm)
            if e != nil {
                panic("Can Not Json Encode")
            }
            m[iTypes.Field(pI).Tag.Get("json")] = string(b)
            continue
        }
        m[iTypes.Field(pI).Tag.Get("json")] = iValues.Field(pI).Interface()
    }
    return m
}
