package wechatPay

import (
    "encoding/json"
    "errors"
    "github.com/reaburoa/wechatPaySDK/wechatPay/request"
)

type Response string

func (r *Response) ToMap() (map[string]interface{}, error) {
    if *r == "" {
        return nil, errors.New("Response Is Empty")
    }
    var mapResp = make(map[string]interface{})
    err := json.Unmarshal([]byte(*r), &mapResp)
    if err != nil {
        return nil, err
    }
    
    return mapResp, nil
}

func (r *Response) GetResponse(req request.Requester, client *WechatPayClient) (map[string]interface{}, error) {
    return nil, nil
}
