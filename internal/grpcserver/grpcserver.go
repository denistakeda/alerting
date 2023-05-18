package grpcserver

import (
	"context"
	"net"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/services/loggerservice"
	"github.com/denistakeda/alerting/internal/storage"
	"github.com/denistakeda/alerting/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	proto.UnimplementedAlertingServer

	address string
	store   storage.Storage
	logger  zerolog.Logger

	server *grpc.Server
}

func NewGRPCServer(log *loggerservice.LoggerService, store storage.Storage, address string) *GRPCServer {
	return &GRPCServer{
		store:   store,
		address: address,
		logger:  log.ComponentLogger("GRPCServer"),
	}
}

func (s *GRPCServer) Start() <-chan error {
	res := make(chan error, 1)

	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		res <- errors.Wrapf(err, "failed to listen address %s", s.address)
		return res
	}

	s.server = grpc.NewServer()
	proto.RegisterAlertingServer(s.server, s)

	s.logger.Info().Msgf("GRPC server is listening on address %s", s.address)
	go func() {
		if err := s.server.Serve(listen); err != nil {
			res <- errors.Wrap(err, "GRPC server failed")
		}
	}()

	return res
}

func (s *GRPCServer) Stop() {
	s.server.Stop()
	s.logger.Info().Msgf("GRPC server was stopped")
}

func (s *GRPCServer) UpdateMetrics(ctx context.Context, req *proto.UpdateMetricsRequest) (*empty.Empty, error) {
	s.logger.Debug().Msgf("got %d metrics", len(req.Metrics))

	ms := make([]*metric.Metric, 0, len(req.Metrics))
	for _, m := range req.Metrics {
		ms = append(ms, metric.FromProto(m))
	}

	if err := s.store.UpdateAll(ctx, ms); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store metrics")
	}

	return &emptypb.Empty{}, nil
}
