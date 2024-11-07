package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

func Start(ctx context.Context, addr, srvName string, registerServer func(s *grpc.Server)) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatal(err)
	}

	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)

	registerServer(srv)

	logrus.Infof("%s grpc server listening...", srvName)

	go listenForStopped(ctx, srv, srvName)

	if err = srv.Serve(lis); err != nil {
		logrus.Fatal(err)
	}
}

func CreateClientConnection(addr string) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cc, err := grpc.NewClient(addr, opts...)
	if err != nil {
		logrus.Fatal(err)
	}

	return cc
}

func listenForStopped(ctx context.Context, grpcServer *grpc.Server, srvName string) {
	defer func() {
		logrus.Warnf("%s grpc server stopped", srvName)
	}()

	for {
		select {
		case <-ctx.Done():
			grpcServer.Stop()
			return
		}
	}
}
