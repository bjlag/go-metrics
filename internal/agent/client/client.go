// HTTP клиент для отправки метрик на сервер.
package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/bjlag/go-metrics/internal/agent/collector"
	"github.com/bjlag/go-metrics/internal/agent/limiter"
	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/signature"
)

const (
	baseURLTemplate = "http://%s:%d"
	urlTemplate     = "%s/updates/"

	timeout          = 100 * time.Millisecond
	maxRetries       = 3
	retryWaitTime    = 200 * time.Millisecond
	retryMaxWaitTime = 500 * time.Millisecond
)

type MetricSender struct {
	client  *resty.Client
	sign    *signature.SignManager
	limiter *limiter.RateLimiter
	baseURL string
	log     log
}

func NewHTTPSender(host string, port int, sign *signature.SignManager, limiter *limiter.RateLimiter, log log) *MetricSender {
	client := resty.New()
	client.SetTimeout(timeout)
	client.SetRetryCount(maxRetries)
	client.SetRetryWaitTime(retryWaitTime)
	client.SetRetryMaxWaitTime(retryMaxWaitTime)

	return &MetricSender{
		client:  client,
		sign:    sign,
		limiter: limiter,
		baseURL: fmt.Sprintf(baseURLTemplate, host, port),
		log:     log,
	}
}

func (s MetricSender) Send(metrics []*collector.Metric) error {
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
				return err
			}
			in.Value = &value
		case collector.Counter:
			value, err := m.CounterValue()
			if err != nil {
				return err
			}
			in.Delta = &value
		default:
			continue
		}

		req = append(req, in)
	}

	jsonb, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %s", err)
	}

	compressed, err := compress(jsonb)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(urlTemplate, s.baseURL)
	request := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(compressed)

	if s.sign.Enable() {
		request = request.SetHeader("HashSHA256", s.sign.Sing(jsonb))
	}

	s.limiter.Acquire()
	defer s.limiter.Release()

	response, err := request.Post(url)
	if err != nil {
		return fmt.Errorf("error sending request to '%s', error %v", url, err)
	}

	s.log.WithField("uri", response.Request.URL).
		WithField("response", string(response.Body())).
		WithField("status", response.StatusCode()).
		Info("sent request")

	return nil
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
