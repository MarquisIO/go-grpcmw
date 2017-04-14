package descriptor

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	annotations "github.com/MarquisIO/BKND-gRPCMiddleware/proto"
)

type Interceptors struct {
	Symbols []string
}

func GetInterceptors(pb proto.Message, desc *proto.ExtensionDesc) (*Interceptors, error) {
	if !proto.HasExtension(pb, desc) {
		return nil, nil
	}
	ext, err := proto.GetExtension(pb, desc)
	if err != nil {
		return nil, err
	}
	interceptors, ok := ext.(*annotations.Interceptors)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an Interceptors", ext)
	}
	return &Interceptors{
		Symbols: interceptors.GetSymbols(),
	}, nil
}
