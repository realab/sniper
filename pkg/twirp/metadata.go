package twirp

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
)

type serverMetadata struct {
	val map[string]map[string]interface{}
}

var ServerMetadata = &serverMetadata{
	val: map[string]map[string]interface{}{},
}

func (m *serverMetadata) Register(in Server) error {
	b, _ := in.ServiceDescriptor()
	zr, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return err
	}
	b2, err := io.ReadAll(zr)
	if err != nil {
		return err
	}
	p := new(descriptorpb.FileDescriptorProto)
	err = proto.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(b2, p)
	if err != nil {
		return err
	}
	fd, err := protodesc.NewFile(p, nil)
	if err != nil {
		return err
	}
	for i := 0; i < fd.Services().Len(); i++ {
		svc := fd.Services().Get(i)
		for j := 0; j < svc.Methods().Len(); j++ {
			method := svc.Methods().Get(j)
			m.val["/twrip/"+string(svc.FullName())+"/"+string(method.Name())] = map[string]interface{}{
				"method":  "POST",
				"methods": []string{"POST"},
			}
		}
	}
	return nil
}

func (m *serverMetadata) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	data := struct {
		Code    int64       `json:"code"`
		Message string      `json:"message"`
		TTL     int64       `json:"ttl"`
		Data    interface{} `json:"data"`
	}{
		Message: "0",
		TTL:     1,
		Code:    0,
		Data:    m.val,
	}
	json.NewEncoder(w).Encode(data)
}
