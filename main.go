/**
https://developers.google.com/protocol-buffers/docs/encoding
 */
package protoc

import (
	"errors"
	"fmt"
	"github.com/pyting/protobuf/proto"
	"strconv"
)

type Message struct {
	Tag            uint64
	WireType       uint64
	WireTypeString string
	Len            int
	Value          interface{}
}

func DecodeProtocolBuf(data []byte) (message []*Message, out string, err error) {
	return tryDecode(data, 0)
}

func tryDecode(data []byte, depth int) (message []*Message, out string, err error) {
	var u uint64
	var str string
	msg := make([]*Message, 0)

	out += "{\n"
	p := proto.NewBuffer(data)
	lenDate := len(data)

	for {
		if p.Index() == lenDate {
			for i := 0; i < depth-1; i++ {
				out += "	"
			}
			if depth > 0 {
				out += "  "
			}
			out += "}\n"
			break
		}

		for i := 0; i < depth; i++ {
			out += "	"
		}

		u, err = p.DecodeVarint()
		if err != nil {
			return
		}

		tag := u >> 3
		wire := u & 7

		switch wire {
		case proto.WireVarint:
			// int32, int64, uint32, uint64, sint32, sint64, bool, enum
			u, err = p.DecodeVarint()
			if err != nil {
				return
			}
			message = append(message, &Message{Tag: tag, WireType: wire, WireTypeString: "varInt", Value: u})
			out += strconv.FormatUint(tag, 10) + ":varInt" + "=" + strconv.FormatUint(u, 10) + "\n"
		case proto.WireFixed64:
			// fixed64, sfixed64, double
			u, err = p.DecodeFixed64()
			if err != nil {
				return
			}
			message = append(message, &Message{Tag: tag, WireType: wire, WireTypeString: "fixed64", Value: u})
			out += strconv.FormatUint(tag, 10) + ":fixed64" + "=" + fmt.Sprintf("%#x", u) + "\n"
		case proto.WireBytes:
			// string, bytes, embedded messages, packed repeated fields
			var r []byte
			r, err = p.DecodeRawBytes(false)
			if err != nil {
				return
			}
			if msg, str, err = tryDecode(r, depth+1); err != nil {
				message = append(message, &Message{Tag: tag, WireType: wire, WireTypeString: "bytes", Len: len(r), Value: r})
				out += strconv.FormatUint(tag, 10) + ":bytes len(" + strconv.Itoa(len(r)) + ")" +
					"=\"" + string(r) + "\"\n"
			} else {
				message = append(message, &Message{Tag: tag, WireType: wire, WireTypeString: "message", Len: len(r), Value: msg})
				out += strconv.FormatUint(tag, 10) + ":" + str
			}
		case proto.WireStartGroup:
			err = errors.New("deprecated")
			return
		case proto.WireEndGroup:
			err = errors.New("deprecated")
			return
		case proto.WireFixed32:
			// fixed32, sfixed32, float
			u, err = p.DecodeFixed32()
			if err != nil {
				return
			}
			message = append(message, &Message{Tag: tag, WireType: wire, WireTypeString: "fixed32", Value: u})
			out += strconv.FormatUint(tag, 10) + ":fixed32 " + "=" + fmt.Sprintf("%#x", u) + "\n"
		default:
			err = errors.New("error wire type")
			return
		}
	}

	err = nil
	return
}
