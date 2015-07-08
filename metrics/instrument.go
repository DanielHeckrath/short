package metrics

import (
	"fmt"
	"io/ioutil"
	"time"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"

	"github.com/DanielHeckrath/short/pb"
	"github.com/DanielHeckrath/short/shorten"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/expvar"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/metrics/statsd"
)

func makeInstrumentation(namespace, name, helpCounter, helpDuration string) (metrics.Counter, metrics.TimeHistogram) {
	counter := metrics.NewMultiCounter(
		expvar.NewCounter(fmt.Sprintf("requests_%s", name)),
		statsd.NewCounter(ioutil.Discard, fmt.Sprintf("requests_%s_total", name), time.Second),
		prometheus.NewCounter(stdprometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: name,
			Name:      "requests_total",
			Help:      helpCounter,
		}, []string{}),
	)
	duration := metrics.NewTimeHistogram(time.Nanosecond, metrics.NewMultiHistogram(
		expvar.NewHistogram(fmt.Sprintf("duration_%s_nanoseconds_total", name), 0, 1e9, 3, 50, 95, 99),
		statsd.NewHistogram(ioutil.Discard, fmt.Sprintf("duration_%s_nanoseconds_total", name), time.Second),
		prometheus.NewSummary(stdprometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: name,
			Name:      "duration_nanoseconds_total",
			Help:      helpDuration,
		}, []string{}),
	))

	return counter, duration
}

// Instrument wraps a shorten function with instrumentation
func Instrument(s shorten.Shortener) shorten.Shortener {
	// `package metrics` domain
	requestsShorten, durationShorten := makeInstrumentation(
		"shortener",
		"shorten",
		"Total number of received shorten requests.",
		"Total nanoseconds spend serving shorten requests.",
	)
	requestsResolve, durationResolve := makeInstrumentation(
		"shortener",
		"resolve",
		"Total number of received resolve requests.",
		"Total nanoseconds spend serving resolve requests.",
	)
	requestsLatest, durationLatest := makeInstrumentation(
		"shortener",
		"latest",
		"Total number of received latest requests.",
		"Total nanoseconds spend serving latest requests.",
	)
	requestsInfo, durationInfo := makeInstrumentation(
		"shortener",
		"info",
		"Total number of received info requests.",
		"Total nanoseconds spend serving info requests.",
	)

	return &instrumentedShortener{
		requestsShorten: requestsShorten,
		requestsResolve: requestsResolve,
		requestsInfo:    requestsInfo,
		requestsLatest:  requestsLatest,

		durationShorten: durationShorten,
		durationResolve: durationResolve,
		durationInfo:    durationInfo,
		durationLatest:  durationLatest,

		shortener: s,
	}
}

type instrumentedShortener struct {
	requestsShorten metrics.Counter
	requestsResolve metrics.Counter
	requestsInfo    metrics.Counter
	requestsLatest  metrics.Counter

	durationShorten metrics.TimeHistogram
	durationResolve metrics.TimeHistogram
	durationInfo    metrics.TimeHistogram
	durationLatest  metrics.TimeHistogram

	shortener shorten.Shortener
}

func (s *instrumentedShortener) Shorten(ctx context.Context, url string) (*pb.ShortURL, error) {
	defer func(begin time.Time) {
		s.requestsShorten.Add(1)
		s.durationShorten.Observe(time.Since(begin))
	}(time.Now())
	return s.shortener.Shorten(ctx, url)
}

func (s *instrumentedShortener) Resolve(ctx context.Context, key string) (*pb.ShortURL, error) {
	defer func(begin time.Time) {
		s.requestsResolve.Add(1)
		s.durationResolve.Observe(time.Since(begin))
	}(time.Now())
	return s.shortener.Resolve(ctx, key)
}

func (s *instrumentedShortener) Info(ctx context.Context, key string) (*pb.ShortURL, error) {
	defer func(begin time.Time) {
		s.requestsInfo.Add(1)
		s.durationInfo.Observe(time.Since(begin))
	}(time.Now())
	return s.shortener.Info(ctx, key)
}

func (s *instrumentedShortener) Latest(ctx context.Context, count int64) ([]*pb.ShortURL, error) {
	defer func(begin time.Time) {
		s.requestsLatest.Add(1)
		s.durationLatest.Observe(time.Since(begin))
	}(time.Now())
	return s.shortener.Latest(ctx, count)
}
