package redismodel

import "strconv"

var Client =GetRedisClient()

type Tx struct {
	Num int
	Req string
}




func (s *Tx) Do() {
	Client.Put(strconv.Itoa(s.Num),s.Req)


	//time.Sleep(10 * time.Millisecond)
}