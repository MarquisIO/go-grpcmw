package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/MarquisIO/BKND-gRPCMiddleware/protoc-gen-grpc-middleware/generator"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func parseRequest(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	req := new(plugin.CodeGeneratorRequest)
	if err = proto.Unmarshal(input, req); err != nil {
		return nil, err
	}
	return req, nil
}

func getResponseFromError(err error) *plugin.CodeGeneratorResponse {
	ret := err.Error()
	return &plugin.CodeGeneratorResponse{Error: &ret}
}

func main() {
	var (
		res *plugin.CodeGeneratorResponse
		g   = generator.New()
	)
	if req, err := parseRequest(os.Stdin); err != nil {
		res = getResponseFromError(err)
	} else if res, err = g.Generate(req); err != nil {
		res = getResponseFromError(err)
	}
	if buf, err := proto.Marshal(res); err != nil {
		log.Fatalf("Could not marshal response: %v", err)
	} else if _, err = os.Stdout.Write(buf); err != nil {
		log.Fatalf("Could not write response to stdout: %v", err)
	}
}
