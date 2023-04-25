package util

import (
	"github.com/topfreegames/pitaya/protos"
	"github.com/topfreegames/pitaya/serialize/protobuf"
)

func PbMarshal(v interface{}) []byte {
	data, err := protobuf.NewSerializer().Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func PbUnmarshal(data []byte, i interface{}) error {
	err := protobuf.NewSerializer().Unmarshal(data, i)
	if err != nil {
		panic(err)
	}
	return nil
}

func PbUnmarshalErr(data []byte) *protos.Error {
	i := new(protos.Error)
	err := protobuf.NewSerializer().Unmarshal(data, i)
	if err != nil {
		panic(err)
	}
	return i
}
