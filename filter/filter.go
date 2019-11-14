package filter

import (
       "crypto/rsa"
       "crypto/x509"
       "encoding/pem"
       "github.com/astaxie/beego"
       "io/ioutil"
)

var IsFilterVerify =beego.AppConfig.String("filter")

type ValidateRequest struct {
       Data string
       DataHash string
       TimeStamp int64
       Sign string
}

//获取公钥
var PubKey *rsa.PublicKey
func Init()  {
       if IsFilterVerify !="false"{
               beego.InsertFilter("*",beego.BeforeExec,UserFilter)
       }else {
               return
       }
       p,err:=ioutil.ReadFile("./conf/cert.pem")
       if err != nil {
               panic("can not get pubkey: "+err.Error())
       }
       block,_:= pem.Decode(p)
       pub,err:=x509.ParseCertificate(block.Bytes)
       if err != nil {
               panic("can not get pubkey: "+err.Error())
       }
       PubKey=pub.PublicKey.(*rsa.PublicKey)
}