package httpclient

import "net/http"

type HTTPClient struct {
	bus chan *task
}

type task struct {
	request  *http.Request
	respChan chan<- *http.Response
	errChan  chan<- error
}

func New(rateLimit int) *HTTPClient {
	bus := make(chan *task)
	client := &http.Client{}

	for i := 0; i < rateLimit; i++ {
		go handleRequests(bus, client)
	}

	return &HTTPClient{
		bus: bus,
	}
}

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
