package grpcmw

import "regexp"

var (
	// TODO: See https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#appendix-a---grpc-for-protobuf
	routeRegexp = regexp.MustCompile("\\/(?:(.+)\\.(?:(.+)\\/(.+)|(.+))|(.+))")
)
