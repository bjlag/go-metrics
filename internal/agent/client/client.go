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
	urlTemplate     = "%s/updates/"

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

func (s MetricSender) Send(metrics []*collector.Metric) (*resty.Response, error) {
	req := make([]model.UpdateIn, 0, len(metrics))
	for _, m := range metrics {
		in := model.UpdateIn{
			ID:    m.Name(),
			MType: m.Kind(),
		}

		switch m.Kind() {
		case collector.Gauge:
			value, err := m.GaugeValue()
			if err != nil {
				return nil, err
			}
			in.Value = &value
		case collector.Counter:
			value, err := m.CounterValue()
			if err != nil {
				return nil, err
			}
			in.Delta = &value
		default:
			continue
		}

		req = append(req, in)
	}

	jsonb, err := json.Marshal(req)
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
