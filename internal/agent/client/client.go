package client

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/bjlag/go-metrics/internal/agent/collector"
	"github.com/bjlag/go-metrics/internal/model"
)

const (
	baseURLTemplate = "http://%s:%d"
	urlTemplate     = "%s/update/"

	timeout       = 100 * time.Millisecond
	maxRetries    = 2
	retryWaitTime = 500 * time.Millisecond
)

type MetricSender struct {
	client  *resty.Client
	baseURL string
}

func NewHTTPSender(host string, port int) *MetricSender {
	client := resty.New()
	client.SetTimeout(timeout)
	client.SetRetryCount(maxRetries)
	client.SetRetryWaitTime(retryWaitTime)

	return &MetricSender{
		client:  client,
		baseURL: fmt.Sprintf(baseURLTemplate, host, port),
	}
}

func (s MetricSender) Send(metric *collector.Metric) (*resty.Response, error) {
	in := &model.UpdateIn{
		ID:    metric.Name(),
		MType: metric.Kind(),
	}

	switch metric.Kind() {
	case collector.Gauge:
		value, err := metric.GaugeValue()
		if err != nil {
			return nil, err
		}
		in.Value = &value
	case collector.Counter:
		value, err := metric.CounterValue()
		if err != nil {
			return nil, err
		}
		in.Delta = &value
	default:
		return nil, fmt.Errorf("unknown metric kind: %s", metric.Kind())
	}

	url := fmt.Sprintf(urlTemplate, s.baseURL)
	request := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(in)

	response, err := request.Post(url)
	if err != nil {
		return nil, fmt.Errorf("error sending request to '%s', error %v", url, err)
	}

	return response, nil
}
