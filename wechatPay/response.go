package wechatPay

import (
    "encoding/xml"
    "errors"
    "fmt"
    "github.com/reaburoa/wechatPaySDK/wechatPay/request"
)

type Response string

func (r *Response) ToMap() (map[string]interface{}, error) {
    if *r == "" {
        return nil, errors.New("Response Is Empty")
    }
    var mapResp interface{}
    fmt.Println(*r)
    err := xml.Unmarshal([]byte(*r), &mapResp)
    fmt.Println(mapResp)
    if err != nil {
        return nil, err
    }
    
    return nil, nil
}

func (r *Response) GetResponse(req request.Requester, client *WechatPayClient) (map[string]interface{}, error) {
    return nil, nil
}
