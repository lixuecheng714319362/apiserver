### 1、 编译

```shell
#项目目录下执行
go build
#编译出apiserver可执行文件
```

### 2、启动API服务

启动服务只依赖conf目录下app.conf 和apiserver.conf文件

##### 2、1 配置 `./conf/app.conf`文件

```txt
appname = apiserver 	#服务名称
httpport = 8088			#监听端口号
runmode = dev			#运行模式
autorender = false		#自动渲染
copyrequestbody = true	#获取请求原始数据
include "apiserver.conf" #导入其他配置文件
```

##### 2、2 配置 `./conf/apiserver.conf`文件

```txt
loglevel = info			#日志级别

#redis 配置在当前版本不使用
redis = false
redishost = 127.0.0.1
redisport = 6379
redismaxidle = 30
redismaxactive = 1000
redisidletimeout = 30

```

##### 2、3 启动apiserver服务

```shell
nohup ./apiserver &
```

### 3、接口

##### 3、1  获取指channel内账本信息

Method : `POST`  	Router : `/api/v1/ledger/info` 

参数:

```json
{
    Data:"",		#json序列化字符串
    DataHash:"hash",	#Data请求数据hash，16进制格式
    TimeStamp: 123123123,	#请求生成时间戳
    Sign:"" 			#base64格式，datahash+timestamp的签名
}


Data：{
   "ConfigPath":"",		#yaml格式SDK配置文件内容 
    "ChannelID":"mychannel",		#指定channelID
    "UserName":"Admin",				#指定用户名称
    "OrgName":"ShuQinOrgOne"		#指定组织MSP
}

```

返回：

```jso
{
  "code": 200,			#请求状态
  "msg": "success",		#状态信息
  "data": {
    "height": 250653,	#最新区块总数
    "currentBlockHash": "wiKw3km2bZbQLGcitAIe0QmlSqf4Q2eLV3mrp2ULfGo=",	#当前块hash
    "previousBlockHash": "DB+FBTzKivrbLYoJuqIzH6EU2PpDcOyOb2RQRUFOlLo="	#上个区块hash
  }
}
```

##### 3、2  获取指定区间区块信息

Method : `POST`  	Router : `/api/v1/ledger/range` 

参数:

```json
{
    Data:"",		#json序列化字符串
    DataHash:"hash",	#Data请求数据hash，16进制格式
    TimeStamp: 123123123,	#请求生成时间戳
    Sign:"" 			#base64格式，datahash+timestamp的签名
}

#Data序列化前结构
Data：{
   "ConfigPath":"",		#yaml格式SDK配置文件内容 
    "ChannelID":"mychannel",		#指定channelID
    "UserName":"Admin",				#指定用户名称
    "OrgName":"ShuQinOrgOne"，		#指定组织MSP
    "Start": 0 ,					#起始区块编号
    "End": 1						#结束区块编号
}

```

返回：

```jso
{
  "code": 200,				#请求状态
  "msg": "success",			#状态信息
  "data": [					#区块数组
    {
      "Number": 3,			#区块编号
      "CurrentBlockHash": "vWNRCDegPQAFuW4Qp93/Jx5j2CvlrUW/UP7OylfHRa0=",	#当前块hash
      "PreviousHash": "zc4Illu+XI+I0td7m7AfAqr9XG0ZnVTJImhKX50abuk=",	#前区块hash
      "DataHash": "84eQPFS2k5phD6HC+GoeDQd8jocwGHc0lVSd8OrLhsA=",		#区块数据hash
      "TransactionNumber": 1,											#区块包含交易数
      "Transactions": [													#交易数组
        {
          "CreatorMSPID": "ShuQinOrgOne",				#发起交易者所属组织
          "CreateID": "Admin@shuqinorgone.com",			#发起交易用户id
          "Type": "ENDORSER_TRANSACTION",				#交易类型
          "Timestamp": 1568080488,						#交易时间戳
          "Nanos": 452109892,							#纳秒级时间戳
          "ChannelId": "mychannel",						#交易所属channel
          "TxId": "f06a275a219bf8f462ce3c1e2f1ef6a02458de195570eeaa1bec06d51c000ef3",
          #交易ID
          "Actions": [		#交易执行内容
            {
              "CCID": "mycc",	#执行交易智能合约ID
              "TxArgs": "[\"invoke\",\"test7858284073146340334\",\"o1p0\"]",#稚嫩合约传参
              "NsRwSets": "",		#交易读写集
              "ReponseStatus": 200	#交易状态
            }
          ]，
          "TransactionFilter": 0		#交易验证结果  （0 有效交易 非0 无效交易  404 未验证）
        }
      ]
    }
  ]
}

```







