package filter

import (
       "crypto/rsa"
       "crypto/x509"
       "encoding/pem"
       "github.com/astaxie/beego"
       "io/ioutil"
       "net/http"
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

       urlManager()
       verifyFilterRegister()

}
func urlManager()  {
       beego.ErrorHandler("404",urlErr)
}
var urlErr http.HandlerFunc= func(w http.ResponseWriter, r *http.Request) {
       w.Write([]byte("Not Found : 404"))
}

func verifyFilterRegister (){

       if IsFilterVerify =="false"{
              return
       }
       beego.InsertFilter("*",beego.BeforeExec,UserFilter)

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
