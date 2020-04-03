package request

var (
    RequestDataXML = "XML"
    RequestDataJSON = "JSON"
)

type Requester interface {
    GetApiUrl() string
    GetNotifyUrl() string
    SetNotifyUrl(url string)
    SetClientIp(str string)
    GetClientIp() string
    GetParams() Requester
    GetRequestDataType() string
    SetSignType(signType string)
    GetSignType() string
}
