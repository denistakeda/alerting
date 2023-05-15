package grpcclient

import (
	"context"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/ports"
	"github.com/denistakeda/alerting/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	address string
}

var _ ports.Client = (*GRPCClient)(nil)

func NewGRPCClient(address string) (*GRPCClient, error) {
	return &GRPCClient{address: address}, nil
}

func (c *GRPCClient) SendMetrics(metrics []*metric.Metric) error {
	ms := make([]*proto.Metric, 0, len(metrics))
	for _, m := range metrics {
		ms = append(ms, m.ToProto())
	}

	conn, err := grpc.Dial(c.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.Wrap(err, "failed to create a connection")
	}
	defer conn.Close()

	client := proto.NewAlertingClient(conn)

	var req proto.UpdateMetricsRequest
	req.Metrics = ms

	_, err = client.UpdateMetrics(context.Background(), &req)
	if err != nil {
		return errors.Wrap(err, "failed to send metrics to the server")
	}
	return nil
}
