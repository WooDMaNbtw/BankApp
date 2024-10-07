package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"log"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwarderFor              = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("md: %+v\n", md)
		// from HTTP
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) != 0 {
			mtdt.UserAgent = userAgents[0]
		}
		// from grpc
		if userAgents := md.Get(userAgentHeader); len(userAgents) != 0 {
			mtdt.UserAgent = userAgents[0]
		}

		// from HTTP
		if clientIPs := md.Get(xForwarderFor); len(clientIPs) != 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
