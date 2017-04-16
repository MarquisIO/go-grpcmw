package descriptor

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	annotations "github.com/MarquisIO/BKND-gRPCMiddleware/proto"
)

// Interceptors defines interceptors to use.
type Interceptors struct {
	Indexes []string
}

// GetInterceptors extracts the `Interceptors` extension (described by `desc`)
// from `pb`.
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
	} else if len(interceptors.GetIndexes()) == 0 {
		return nil, nil
	}
	return &Interceptors{
		Indexes: interceptors.GetIndexes(),
	}, nil
}
