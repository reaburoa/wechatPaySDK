package wechatPay

import (
    "crypto/tls"
    "encoding/pem"
    "fmt"
    "github.com/reaburoa/elec-signature/signature"
    "github.com/reaburoa/wechatPaySDK/wechatPay/request"
    "golang.org/x/crypto/pkcs12"
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "net/url"
    "reflect"
    "sort"
    "strconv"
    "strings"
    "time"
)

var (
    signTypeMd5    = "MD5"
    signTypeRSA2   = "RSA2"
    nonceStrLength = 16
)

type WechatPayClient struct {
    AppId     string
    MchId     string
    SecretKey string
    AppSecret string
    CertData  []byte
    Client    *http.Client
}

func NewClient(appId, mchId, secretKey, cert, appSecret string) *WechatPayClient {
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

func (w *WechatPayClient) number2String(number interface{}) string {
    kStr := reflect.TypeOf(number).Kind()
    switch kStr {
    case reflect.Int64:
        number = strconv.FormatInt(number.(int64), 10)
    case reflect.Int32:
        number = strconv.FormatInt(number.(int64), 10)
    case reflect.Int:
        number = strconv.Itoa(number.(int))
    case reflect.Float64:
        number = strconv.FormatFloat(number.(float64), 'f', -1, 64)
    case reflect.Float32:
        number = strconv.FormatFloat(number.(float64), 'f', -1, 64)
    case reflect.String:
    default:
        number = ""
    }
    
    return number.(string)
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

func (w *WechatPayClient) formatUrlValue(data map[string]interface{}) url.Values {
    var formData = make(url.Values)
    for key, val := range data {
        val = w.number2String(val)
        formData.Set(key, val.(string))
    }
    
    return formData
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

func (w *WechatPayClient) genReqData(req request.Requester) map[string]interface{} {
    commonReq := CommonRequest{
        AppId:     w.AppId,
        MchId:     w.MchId,
        NonceStr:  w.GenNonceStr(),
        Requester: req,
    }
    clientMap := commonReq.toMap()
    signType := req.GetSignType()
    if signType == signTypeMd5 {
        commonReq.Sign = w.genSignByMd5(w.genSignContent(clientMap))
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
    }
    
    return nil
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
    fmt.Println(resp, err)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    /*parsedBody, err := w.parseBody(body, req)
      if err != nil {
          return "", err
      }
      if parsedBody["sign"] != "" {
          checkRet := w.checkSign(parsedBody["sign_data"], parsedBody["sign"], a.SignType)
          if checkRet != true {
              return "", errors.New("Check Sign Error")
          }
      }*/
    return Response(body), nil
}

func (w *WechatPayClient) Execute(req request.Requester, method string) (Response, error) {
    xmlStr := w.getRequestData(req)
    fmt.Println(xmlStr)
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
    /*parsedBody, err := w.parseBody(body, req)
      if err != nil {
          return "", err
      }
      if parsedBody["sign"] != "" {
          checkRet := w.checkSign(parsedBody["sign_data"], parsedBody["sign"], a.SignType)
          if checkRet != true {
              return "", errors.New("Check Sign Error")
          }
      }*/
    return Response(body), nil
}

/*func (w *WechatPayClient) parseBody(body []byte, req request.Requester) (map[string]string, error) {
    bodyStr := string(body)
    responseReg := a.methodNameToResponseName(req)
    if strings.Index(bodyStr, errResponse) > -1 {
        responseReg = errResponse
    }
    mapResp := make(map[string]interface{})
    err := json.Unmarshal(body, &mapResp)
    if err != nil {
        return nil, err
    }
    reg, sign := "", ""
    if strings.Index(bodyStr, signTag) == -1 {
        reg = "{\"" + responseReg + `":\s?{(.*)}`
    } else {
        reg = "{\"" + responseReg + `":\s?{(.*)},`
        sign = mapResp["sign"].(string)
    }
    re, err := regexp.Compile(reg)
    if err != nil {
        return nil, err
    }
    toVerifyStr := re.FindString(bodyStr)
    start := len("{\"" + responseReg + "\":")
    end := len(toVerifyStr) - 1
    return map[string]string{
        "sign_data": strings.Trim(string(toVerifyStr[start:end]), ""),
        "sign":      sign,
    }, nil
}*/
