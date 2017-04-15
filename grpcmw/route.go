package grpcmw

import "regexp"

var (
	routeRegexp = regexp.MustCompile(`\/(?:(.+)\.)?(.+)\/(.+)`)
)
