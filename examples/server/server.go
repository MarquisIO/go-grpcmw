package serverpb

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/MarquisIO/go-grpcmw/examples/proto"
)

// Example implements the `pb.Example` service
type Example struct{}

// Method prints:
// "Received : <message>"
func (e *Example) Method(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	fmt.Printf("Received : %s\n", msg.Msg)
	return msg, nil
}
