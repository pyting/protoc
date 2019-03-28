package protoc

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/pyting/protoc/protoc2"
	"math"
	"testing"
)

func TestDecodeProtocolBuf(t *testing.T) {
	m0 := &protoc2.M0{
		M1: []*protoc2.M1{
			{
				Str_1:      proto.String("pyting"),                    // string bytes 都会被解析成bytes
				Bytes_1:    []byte{0x8, 0x10, 0x10, 0x11, 0x18, 0x12}, // bytes 会被解析成 message
				Int32_1:    proto.Int(math.MaxInt32),
				Int32_2:    proto.Int(math.MinInt32),
				Uint32_1:   proto.Uint32(math.MaxUint32),
				Sint32_1:   proto.Int(math.MaxInt32),     // 会使用zigzag编码
				Sint32_2:   proto.Int(math.MinInt32),     // 会使用zigzag编码
				Fixed32_1:  proto.Uint32(math.MaxUint32), // 32无符号整形
				Sfixed32_1: proto.Int(math.MaxInt32),     // 会使用zigzag编码
				Sfixed32_2: proto.Int(math.MinInt32),     // 会使用zigzag编码
				Int64_1:    proto.Int64(math.MaxInt64),
				Int64_2:    proto.Int64(math.MinInt64),
				Uint64_1:   proto.Uint64(math.MaxUint64),
				Sint64_1:   proto.Int64(math.MaxInt64),   // 会使用zigzag编码
				Sint64_2:   proto.Int64(math.MinInt64),   // 会使用zigzag编码
				Fixed64_1:  proto.Uint64(math.MaxUint64), // 有符号类型,会使用zigzag编码
				Sfixed64_1: proto.Int64(math.MaxInt64),   // 有符号类型,会使用zigzag编码
				Sfixed64_2: proto.Int64(math.MinInt64),   // 有符号类型,会使用zigzag编码
				Bool_1:     proto.Bool(false),
				Bool_2:     proto.Bool(true),
				Float_1:    proto.Float32(math.MaxFloat32), // 有符号类型,会使用zigzag编码
				Double_1:   proto.Float64(math.MaxFloat64), // 有符号类型,会使用zigzag编码
				Enum_1: []protoc2.M1_Enum_1{
					protoc2.M1_Enum_1(protoc2.M1_Enum_1_1),
				},
			},
		},
	}

	b, err := proto.Marshal(m0)
	if err != nil {
		t.Error(err)
	}

	msg, out, err := DecodeProtocolBuf(b)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(out)
	for i, v := range msg {
		// TODO
		fmt.Println(i, v)
	}
}
