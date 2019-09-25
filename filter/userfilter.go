package filter

import (
	"apiserver/controllers/tool"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"io/ioutil"
	"strconv"
	"time"
)
const(
	timeOut = 30
)

type validate struct {
	Data string
	DataHash string
	TimeStamp int64
	Sign string
}

//获取公钥
var PubKey *rsa.PublicKey
func init()  {
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
//验证签名信息
var UserFilter beego.FilterFunc =   func (ctx *context.Context){
	data:=ctx.Input.RequestBody
	var v = validate{}
	err:=json.Unmarshal(data,&v)
	if err != nil {
		FilterResError(ctx,"request data error")
	}
	hash:=sha256.Sum256([]byte(v.Data))
	h:=hex.EncodeToString(hash[:])
	if h!=v.DataHash {
		FilterResError(ctx,"Datahash error")
	}
	now:=time.Now().Unix()
	beego.Debug("now:",now , "timestamp:",v.TimeStamp ,"outtime",v.TimeStamp+timeOut)
	if !(now>v.TimeStamp-timeOut &&  now < v.TimeStamp+timeOut){
		FilterResError(ctx,"Time validate out")
	}
	sign ,_:=base64.StdEncoding.DecodeString(v.Sign)
	hashed:=sha256.Sum256([]byte(v.DataHash+strconv.FormatInt(v.TimeStamp,10)))
	err=rsa.VerifyPKCS1v15(PubKey,crypto.SHA256,hashed[:],sign)
	if err != nil {
		FilterResError(ctx,"Time validate error")
	}
}
//过滤错误返回
func FilterResError(ctx *context.Context,errMsg string) {
	var res = tool.Result{
		Code: 400,
		Msg:  errMsg,
	}
	msg, _ := json.Marshal(res)
	_ = ctx.Output.Body([]byte(msg))
}
