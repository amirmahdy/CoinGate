package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgent := md.Get("user-agent"); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}
	}

	if peer, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = peer.Addr.String()
	}
	return &mtdt
}
