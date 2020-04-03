package wechatPay

import (
    "encoding/json"
    "encoding/xml"
    "errors"
    "github.com/reaburoa/wechatPaySDK/wechatPay/request"
    "strings"
)

type Response string

func (r *Response) ToMap(req request.Requester) (map[string]interface{}, error) {
    switch req.GetRequestDataType() {
    case request.RequestDataXML:
        return r.XmlToMap()
    case request.RequestDataJSON:
        return r.JsonToMap()
    }
    
    return nil, nil
}

func (r *Response) XmlToMap() (map[string]interface{}, error) {
    if string(*r) == "" {
        return nil, errors.New("Response Is Empty")
    }
    mm := make(map[string]interface{})
    key, val := "", ""
    decoder := xml.NewDecoder(strings.NewReader(string(*r)))
    for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
        switch token := t.(type) {
        case xml.StartElement: // 处理元素开始（标签）
            key = token.Name.Local
        case xml.EndElement: // 处理元素结束（标签）
        case xml.CharData: // 处理字符数据（这里就是元素的文本）
            content := string([]byte(token))
            if content == "\n" || content == "\r" || content == "\r\n" {
                continue
            }
            val = content
        }
        if key != "xml" {
            mm[key] = val
        }
    }
    return mm, nil
}

func (r *Response) JsonToMap() (map[string]interface{}, error) {
    mm := make(map[string]interface{})
    err := json.Unmarshal([]byte(*r), &mm)
    if err != nil {
        return nil, errors.New("Json Decode Failed")
    }
    
    return mm, nil
}

func (r *Response) OriginResp() string {
    return string(*r)
}
