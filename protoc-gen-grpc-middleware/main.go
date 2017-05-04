package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/MarquisIO/go-grpcmw/protoc-gen-grpc-middleware/descriptor"
	"github.com/MarquisIO/go-grpcmw/protoc-gen-grpc-middleware/template"
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
	var res *plugin.CodeGeneratorResponse
	if req, err := parseRequest(os.Stdin); err != nil {
		res = getResponseFromError(err)
	} else if pkgs, err := descriptor.Parse(req); err != nil {
		res = getResponseFromError(err)
	} else if res, err = template.Apply(pkgs); err != nil {
		res = getResponseFromError(err)
	}
	if buf, err := proto.Marshal(res); err != nil {
		log.Fatalf("Could not marshal response: %v", err)
	} else if _, err = os.Stdout.Write(buf); err != nil {
		log.Fatalf("Could not write response to stdout: %v", err)
	}
}
