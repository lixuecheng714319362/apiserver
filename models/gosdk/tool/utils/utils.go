package utils

import (
	"encoding/asn1"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/common/util"
	"math"
)

type CollectionConfig struct {
	Name		      	string
	MemberOrgsPolicy  	string
	RequiredPeerCount 	int32
	MaximumPeerCount  	int32
	BlockToLive		  	uint64
	MemberOnlyRead		bool
}
func NewCollectionConfig(colName, policy string, reqPeerCount, maxPeerCount int32, blockToLive uint64,memberOnlyRead bool) (*common.CollectionConfig, error) {
	p, err := cauthdsl.FromString(policy)
	if err != nil {
		return nil, err
	}
	cpc := &common.CollectionPolicyConfig{
		Payload: &common.CollectionPolicyConfig_SignaturePolicy{
			SignaturePolicy: p,
		},
	}
	return &common.CollectionConfig{
		Payload: &common.CollectionConfig_StaticCollectionConfig{
			StaticCollectionConfig: &common.StaticCollectionConfig{
				Name:              colName,
				MemberOrgsPolicy:  cpc,
				RequiredPeerCount: reqPeerCount,
				MaximumPeerCount:  maxPeerCount,
				BlockToLive:       blockToLive,
				MemberOnlyRead: memberOnlyRead,
			},
		},
	}, nil
}

//计算当前区块hash
type asn1Header struct {
	Number       int64
	PreviousHash []byte
	DataHash     []byte
}

func GetCurrentBlockHash(block *common.Block)([]byte,error) {
	asn1Header := asn1Header{
		PreviousHash: block.Header.PreviousHash,
		DataHash:     block.Header.DataHash,
	}
	if block.Header.Number > uint64(math.MaxInt64) {
		return nil,fmt.Errorf("Golang does not currently support encoding uint64 to asn1")
	} else {
		asn1Header.Number = int64(block.Header.Number)
	}
	resault,err:= asn1.Marshal(asn1Header)
	if err != nil {
		return nil,err
	}
	hash:=util.ComputeSHA256(resault)
	return hash,nil
}

