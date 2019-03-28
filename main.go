/**
https://developers.google.com/protocol-buffers/docs/encoding
 */
package protoc

import (
	"errors"
	"github.com/pyting/protobuf/proto"
	"strconv"
)

func DecodeProtocolBuf(data []byte) (out string, err error) {
	return tryDecode(data, 0)
}

func tryDecode(data []byte, depth int) (out string, err error) {
	var u uint64
	var str string

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
			//fmt.Printf("%3d: fetching op err %v\n", p.Index(), err)
			return
		}

		tag := u >> 3
		wire := u & 7

		switch wire {
		case proto.WireVarint:
			u, err = p.DecodeVarint()
			if err != nil {
				//fmt.Printf("%3d: t=%3d varint err %v\n", p.Index(), tag, err)
				return
			}
			out += strconv.FormatUint(tag, 10) + ":varInt" + "=" + strconv.FormatUint(u, 10) + "\n"
		case proto.WireFixed64:
			u, err = p.DecodeFixed64()
			if err != nil {
				//fmt.Printf("%3d: t=%3d fixed64 err %v\n", p.Index(), tag, err)
				return
			}
			out += strconv.FormatUint(tag, 10) + ":fixed64" + "=" + strconv.FormatUint(u, 10) + "\n"
		case proto.WireBytes:
			var r []byte
			r, err = p.DecodeRawBytes(false)
			if err != nil {
				//fmt.Printf("%3d: t=%3d bytes err %v\n", p.Index(), tag, err)
				return
			}
			if str, err = tryDecode(r, depth+1); err != nil {
				out += strconv.FormatUint(tag, 10) + ":bytes len(" + strconv.Itoa(len(r)) + ")" +
					"=\"" + string(r) + "\"\n"
			} else {
				out += strconv.FormatUint(tag, 10) + ":" + str
			}
		case proto.WireStartGroup:
			err = errors.New("TODO WireStartGroup")
			//fmt.Printf("%3d: t=%3d bytes err %v\n", p.Index(), tag, err)
			return
		case proto.WireEndGroup:
			err = errors.New("TODO WireEndGroup")
			//fmt.Printf("%3d: t=%3d bytes err %v\n", p.Index(), tag, err)
			return
		case proto.WireFixed32:
			u, err = p.DecodeFixed32()
			if err != nil {
				//fmt.Printf("%3d: t=%3d fixed32 err %v\n", p.Index(), tag, err)
				return
			}
			out += strconv.FormatUint(tag, 10) + ":fixed32 " + "=" + strconv.FormatUint(u, 10) + "\n"
		default:
			err = errors.New("error wire type")
			return
		}
	}

	err = nil
	return
}
