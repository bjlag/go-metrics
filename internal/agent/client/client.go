package client

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/bjlag/go-metrics/internal/agent/collector"
)

const (
	urlTemplate = "%s/update/%s/%s/%v"

	timeout       = 100 * time.Millisecond
	maxRetries    = 2
	retryWaitTime = 500 * time.Millisecond
)

type MetricSender struct {
	client  *resty.Client
	baseUrl string
}

func NewHTTPSender(baseUrl string) *MetricSender {
	client := resty.New()
	client.SetTimeout(timeout)
	client.SetRetryCount(maxRetries)
	client.SetRetryWaitTime(retryWaitTime)

	return &MetricSender{
		client:  client,
		baseUrl: baseUrl,
	}
}

func (s MetricSender) Send(metric *collector.Metric) (*resty.Response, error) {
	url := fmt.Sprintf(urlTemplate, s.baseUrl, metric.Kind(), metric.Name(), metric.Value())
	request := s.client.R().SetHeader("Content-Type", "text/plain")

	response, err := request.Post(url)
	if err != nil {
		return nil, fmt.Errorf("error sending request to '%s', error %v", url, err)
	}

	return response, nil
}
