package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// HTTPClient is a rate-limited client
type HTTPClient struct {
	bus chan *task
}

type task struct {
	request  *http.Request
	respChan chan<- *http.Response
	errChan  chan<- error
}

// New instantiates a new HTTPClient
func New(rateLimit int, cert string) (*HTTPClient, error) {
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
		bus: bus,
	}, nil
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
