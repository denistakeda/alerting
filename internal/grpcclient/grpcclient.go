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
	client  proto.AlertingClient
	conn    *grpc.ClientConn
}

var _ ports.Client = (*GRPCClient)(nil)

func NewGRPCClient(address string) (*GRPCClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a connection")
	}

	client := proto.NewAlertingClient(conn)

	return &GRPCClient{
		address: address,
		client:  client,
		conn:    conn,
	}, nil
}

func (c *GRPCClient) SendMetrics(metrics []*metric.Metric) error {
	ms := make([]*proto.Metric, 0, len(metrics))
	for _, m := range metrics {
		ms = append(ms, m.ToProto())
	}

	var req proto.UpdateMetricsRequest
	req.Metrics = ms

	_, err := c.client.UpdateMetrics(context.Background(), &req)
	if err != nil {
		return errors.Wrap(err, "failed to send metrics to the server")
	}

	return nil
}

func (c *GRPCClient) Stop() error {
	return c.conn.Close()
}
