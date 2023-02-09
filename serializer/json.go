package serializer

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func ProtobufToJson(message proto.Message) (string, error) {
	marshaller := jsonpb.Marshaler{
		EnumsAsInts:  true,
		OrigName:     true,
		EmitDefaults: true,
		Indent:       "    ",
	}

	return marshaller.MarshalToString(message)
}
