package request

var (
    RequestDataXML = "XML"
    RequestDataJSON = "JSON"
)

type Requester interface {
    GetApiUrl() string
    GetNotifyUrl() string
    SetNotifyUrl(url string)
    GetSignType() string
    GetClientIp() string
    GetParams() Requester
    GetRequestDataType() string
}
