package sender

import (
	"fmt"
	"net/http"

	"github.com/bjlag/go-metrics/internal/agent/collector"
)

const (
	urlTemplate = "http://127.0.0.1:8080/update/%s/%s/%v"
)

type MetricSender struct {
	client *http.Client
}

func NewHttpSender(client *http.Client) *MetricSender {
	return &MetricSender{
		client: client,
	}
}

func (s MetricSender) Send(metric *collector.Metric) (*http.Response, error) {
	url := fmt.Sprintf(urlTemplate, metric.Kind(), metric.Name(), metric.Value())
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to '%s', error %v", url, err)
	}

	request.Header.Set("Content-Type", "text/plain")

	response, err := s.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request to '%s', error %v", url, err)
	}

	return response, nil
}
