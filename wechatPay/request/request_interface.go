package request

type Requester interface {
    GetApiUrl() string
    GetNotifyUrl() string
    SetNotifyUrl(url string)
    GetSignType() string
    GetClientIp() string
    GetParams() Requester
}
