package wechatPay

import (
    "crypto/tls"
    "encoding/json"
    "encoding/pem"
    "errors"
    "fmt"
    "github.com/reaburoa/elec-signature/signature"
    "github.com/reaburoa/wechatPaySDK/wechatPay/request"
    "golang.org/x/crypto/pkcs12"
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "sort"
    "strings"
    "time"
)

var (
    signTypeMd5    = "MD5"
    signTypeHmac   = "HMAC-SHA256"
    nonceStrLength = 16
)

type WechatPayClient struct {
    AppId     string
    MchId     string
    SecretKey string
    AppSecret string
    SignType  string
    CertData  []byte
    Client    *http.Client
}

func NewClient(appId, mchId, secretKey, cert, appSecret, signType string) *WechatPayClient {
    certData, err := ioutil.ReadFile(cert)
    if err != nil {
        panic("Cert File Error")
    }
    return &WechatPayClient{
        AppId:     appId,
        MchId:     mchId,
        SecretKey: secretKey,
        AppSecret: appSecret,
        CertData:  certData,
        SignType:  signType,
        Client:    http.DefaultClient,
    }
}

func (w *WechatPayClient) sortContentByKeys(data map[string]interface{}) []string {
    keys := []string{}
    for k, _ := range data {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    return keys
}

func (w *WechatPayClient) genSignContent(data map[string]interface{}) string {
    sortedKeys := w.sortContentByKeys(data)
    toSignData := []string{}
    for _, key := range sortedKeys {
        value := data[key]
        if value == nil || strings.Trim(value.(string), "") == "" {
            continue
        }
        toSignData = append(toSignData, fmt.Sprintf("%s=%v", key, value))
    }
    return strings.Join(toSignData, "&")
}

func (w *WechatPayClient) genRand(min, max int) int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(max-min) + min
}

func (w *WechatPayClient) GenNonceStr() string {
    salt := make([]rune, nonceStrLength)
    mid := nonceStrLength / 2
    for i := 0; i < mid; i++ {
        salt[i] = rune(w.genRand(48, 90))
    }
    for j := mid; j < nonceStrLength; j++ {
        salt[j] = rune(w.genRand(97, 122))
    }
    
    return string(salt)
}

func (w *WechatPayClient) genSignByMd5(data string) string {
    str := fmt.Sprintf("%s&key=%s", data, w.SecretKey)
    return strings.ToUpper(signature.Md5(str))
}

func (w *WechatPayClient) gentSignByHmacSHA256(data string) string {
    str := fmt.Sprintf("%s&key=%s", data, w.SecretKey)
    return strings.ToUpper(signature.Hmac(str, w.SecretKey, "SHA-256"))
}

func (w *WechatPayClient) checkSign(req request.Requester, resp Response) bool {
    mm, err := w.getResponseMap(req, resp)
    if err != nil {
        return false
    }
    sign := mm["sign"]
    delete(mm, "sign")
    str := w.genSignContent(mm)
    toSign := ""
    switch w.SignType {
    case signTypeMd5:
        toSign = w.genSignByMd5(str)
    case signTypeHmac:
        toSign = w.gentSignByHmacSHA256(str)
    }
    return toSign == sign
}

func (w *WechatPayClient) toXml(req request.Requester) string {
    reqMap := w.genReqData(req)
    xml := []string{"<xml>"}
    for k, v := range reqMap {
        item := ""
        if v.(string) != "" {
            item = fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v.(string), k)
        }
        xml = append(xml, item)
    }
    xml = append(xml, "</xml>")
    return strings.Join(xml, "")
}

func (w *WechatPayClient) toJson(req request.Requester) string {
    reqMap := w.genReqData(req)
    bt, err := json.Marshal(reqMap)
    if err != nil {
        return ""
    }
    return string(bt)
}

func (w *WechatPayClient) genReqData(req request.Requester) map[string]interface{} {
    commonReq := CommonRequest{
        AppId:     w.AppId,
        MchId:     w.MchId,
        NonceStr:  w.GenNonceStr(),
        Requester: req,
    }
    clientMap := commonReq.toMap()
    if req.GetSignType() == "" {
        w.SignType = signTypeMd5
    } else {
        w.SignType = req.GetSignType()
    }
    switch w.SignType {
    case signTypeMd5:
        commonReq.Sign = w.genSignByMd5(w.genSignContent(clientMap))
    case signTypeHmac:
        commonReq.Sign = w.gentSignByHmacSHA256(w.genSignContent(clientMap))
    }
    clientMap["sign"] = commonReq.Sign
    return clientMap
}

func (w *WechatPayClient) pkcs12ToPem() tls.Certificate {
    blocks, err := pkcs12.ToPEM(w.CertData, w.MchId)
    defer func() {
        if x := recover(); x != nil {
            log.Print(x)
        }
    }()
    if err != nil {
        panic(err)
    }
    var pemData []byte
    for _, b := range blocks {
        pemData = append(pemData, pem.EncodeToMemory(b)...)
    }
    cert, err := tls.X509KeyPair(pemData, pemData)
    if err != nil {
        panic(err)
    }
    return cert
}

func (w *WechatPayClient) getRequestData(req request.Requester) interface{} {
    switch req.GetRequestDataType() {
    case request.RequestDataXML:
        return w.toXml(req)
    case request.RequestDataJSON:
        return w.toJson(req)
    }
    
    return nil
}

func (w *WechatPayClient) getResponseMap(req request.Requester, resp Response) (map[string]interface{}, error) {
    switch req.GetRequestDataType() {
    case request.RequestDataXML:
        return resp.XmlToMap()
    case request.RequestDataJSON:
        return resp.JsonToMap()
    }
    
    return nil, nil
}

func (w *WechatPayClient) ExecuteWithCert(req request.Requester, method string) (Response, error) {
    xmlStr := w.getRequestData(req)
    buf := strings.NewReader(xmlStr.(string))
    cert := w.pkcs12ToPem()
    config := &tls.Config{
        Certificates: []tls.Certificate{cert},
    }
    transport := &http.Transport{
        TLSClientConfig:    config,
        DisableCompression: true,
    }
    h := &http.Client{Transport: transport}
    resp, err := h.Post(req.GetApiUrl(), "application/xml; charset=utf-8", buf)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    checkRet := w.checkSign(req, Response(body))
    if checkRet != true {
        return "", errors.New("Check Sign Error")
    }
    return Response(body), nil
}

func (w *WechatPayClient) Execute(req request.Requester, method string) (Response, error) {
    xmlStr := w.getRequestData(req)
    buf := strings.NewReader(xmlStr.(string))
    reqes, err := http.NewRequest(method, req.GetApiUrl(), buf)
    if err != nil {
        return "", err
    }
    resp, err := w.Client.Do(reqes)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    checkRet := w.checkSign(req, Response(body))
    if checkRet != true {
        return "", errors.New("Check Sign Error")
    }
    
    return Response(body), nil
}
