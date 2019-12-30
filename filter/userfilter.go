package filter

import (
	"apiserver/controllers/tool"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strconv"
	"time"
)
const(
	timeOut = 30
)




//验证签名信息
var UserFilter beego.FilterFunc =   func (ctx *context.Context){
	beego.Debug("start userfilter",time.Now().Unix())
	data:=ctx.Input.RequestBody
	var v = ValidateRequest{}
	err:=json.Unmarshal(data,&v)
	if err != nil {
		FilterResError(ctx,"request data unmarshal error")
		return
	}
	if v.Sign==""||v.Data==""||v.DataHash==""{
		FilterResError(ctx,"request data empty ")
		return
	}
	hash:=sha256.Sum256([]byte(v.Data))
	h:=hex.EncodeToString(hash[:])
	if h!=v.DataHash {
		FilterResError(ctx,"Datahash error")
		return
	}
	now:=time.Now().Unix()
	beego.Debug("now:",now , "timestamp:",v.TimeStamp ,"outtime",v.TimeStamp+timeOut)
	if !(now>v.TimeStamp-timeOut &&  now < v.TimeStamp+timeOut){
		FilterResError(ctx,"Time ValidateRequest error")
		return
	}
	sign ,_:=base64.StdEncoding.DecodeString(v.Sign)
	hashed:=sha256.Sum256([]byte(v.DataHash+strconv.FormatInt(v.TimeStamp,10)))
	err=rsa.VerifyPKCS1v15(PubKey,crypto.SHA256,hashed[:],sign)
	if err != nil {
		FilterResError(ctx,"Time validate error")
		return
	}
	beego.Debug("end userfilter",time.Now().Unix())

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
