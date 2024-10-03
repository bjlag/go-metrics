package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
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

	jsonb, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metric: %s", err)
	}

	compressed, err := compress(jsonb)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(urlTemplate, s.baseURL)
	request := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(compressed)

	response, err := request.Post(url)
	if err != nil {
		return nil, fmt.Errorf("error sending request to '%s', error %v", url, err)
	}

	return response, nil
}

func compress(src []byte) ([]byte, error) {
	var buf bytes.Buffer

	wr, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip writer: %s", err)
	}
	_, err = wr.Write(src)
	if err != nil {
		return nil, fmt.Errorf("failed to compress metric: %s", err)
	}

	err = wr.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %s", err)
	}

	return buf.Bytes(), nil
}
