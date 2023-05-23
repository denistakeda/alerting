package httpclient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/ports"
	"github.com/pkg/errors"
)

// HTTPClient is a rate-limited client
type HTTPClient struct {
	bus     chan *task
	address string
}

var _ ports.Client = (*HTTPClient)(nil)

type task struct {
	request  *http.Request
	respChan chan<- *http.Response
	errChan  chan<- error
}

// New instantiates a new HTTPClient
func New(rateLimit int, cert string, address string) (*HTTPClient, error) {
	client := &http.Client{}

	if cert != "" {
		caCert, err := os.ReadFile(cert)
		if err != nil {
			return nil, errors.Wrap(err, "unable to find certificate file")
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		}
	}

	bus := make(chan *task)
	for i := 0; i < rateLimit; i++ {
		go handleRequests(bus, client)
	}

	return &HTTPClient{
		bus:     bus,
		address: address,
	}, nil
}

func (c *HTTPClient) SendMetrics(metrics []*metric.Metric) error {
	url := fmt.Sprintf("%s/updates/", c.address)
	m, err := json.Marshal(metrics)
	if err != nil {
		return errors.Wrap(err, "failed to marshal metrics")
	}
	body := bytes.NewBuffer(m)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return errors.Wrap(err, "failed to create a request")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return errors.Wrapf(err, "unable to file a request to URL: %s", url)
	}
	if err := resp.Body.Close(); err != nil {
		return errors.Wrap(err, "unable to close a body")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("not successfull status %d", resp.StatusCode)
	}

	return nil
}

func (*HTTPClient) Stop() error {
	// Do nothing
	return nil
}

// Do executes the query with rate-limiting mechanism
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	respChan := make(chan *http.Response)
	errChan := make(chan error)
	t := &task{
		request:  req,
		respChan: respChan,
		errChan:  errChan,
	}

	c.bus <- t
	select {
	case res := <-respChan:
		return res, nil
	case err := <-errChan:
		return nil, err
	}
}

func handleRequests(bus <-chan *task, client *http.Client) {
	for t := range bus {
		resp, err := client.Do(t.request)
		if err != nil {
			t.errChan <- err
			continue
		}

		if err := resp.Body.Close(); err != nil {
			t.errChan <- err
			continue
		}

		t.respChan <- resp
	}
}
