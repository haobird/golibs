# MD5 sign 签名校验

## 默认字段声明

KeyNameTimeStamp    = "timestamp"
KeyNameNonceStr     = "nonce"
KeyNameAppID        = "appid"
KeyNameSign         = "sign"
KeyNameAppSecret    = "appsecret"

## 功能

* 根据传入参数生成签名
* 获得签名过程中的一系列阶段
* 自动生成时间戳和随机字符串（缺少的自动补全）
* 指定校验字段或者不包含字段
* 校验时间、校验随机字符串
* 修改默认的约定字段

## 其中signer是生成的功能

## 校验功能，直接获取到相关字段，然后校验即可

## 配置说明

* 默认字段的配置（重置）
* 超时时间的配置（校验使用）
* 加密函数的指定
* 前后缀、secret的值

## signer需要的参数

* 默认字段的指定
* 前后缀的指定
* secretVal


