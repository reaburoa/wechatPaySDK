# 请求微信接口封装类
### 目的在于在应用中更加快速、简便的使用微信接口以及新增接口。使得整个微信的接口使用全部面向对象化。

### 安装
```
go get -u github.com/reaburoa/wechatSDK
```

### 使用
在使用中只需要初始化相关微信接口类，即刻快速在应用中使用微信收款等一系列操作。接口返回数据可以根据需要获得不同类型的数据。

```go
client := wechatPay.NewClient(
    "wxbf***",
    "mch_id",
    "secretKey",
    "cert",
    "appSecret",
    "HMAC-SHA256",
)

now := time.Now()
start := now.Format("20060102150405")
req := &request.WechatPayMicropay{
  SubMchId: "mch_id",
  Body: "testorder323",
  OutTradeNo: "testorder1223333",
  TotalFee: "1",
  AuthCode: "13516",
  TimeStart: start,
  TimeExpire: now.Add(2*time.Minute).Format("20060102150405"),
  SceneInfo: request.SceneInfoData{
      Id:       "11",
      Name:     "2233",
      AreaCode: "332",
      Address:  "sddd",
  },
}
req.SetClientIp("112.114.123.47")
resp, _ := client.Execute(req, "POST")
fmt.Println(resp.XmlToMap())
```
