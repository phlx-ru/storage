package metrics

import (
	"errors"
	"time"

	"gopkg.in/alexcesaro/statsd.v2"

	"github.com/go-kratos/kratos/v2/log"
)

type Metrics interface {
	Count(bucket string, n interface{})
	Increment(bucket string)
	Gauge(bucket string, value interface{})
	Timing(bucket string, value interface{})
	Histogram(bucket string, value interface{})
	NewTiming() statsd.Timing
	Close()
}

func New(address, name string, mute bool) (*statsd.Client, error) {
	metrics, err := statsd.New(
		statsd.Address(address),
		statsd.ErrorHandler(
			func(err error) {
				log.Warnf(`failed to send metrics: %v`, err)
			},
		),
		statsd.Mute(mute),
		statsd.Prefix(name),
		statsd.FlushPeriod(1*time.Second),
	)
	if metrics == nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("metrics client is undefined")
	}
	if err != nil && !mute {
		return nil, err
	}
	return metrics, nil
}
