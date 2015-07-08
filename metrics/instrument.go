package metrics

import (
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

// Instrument wraps a shorten function with instrumentation
func Instrument(s shorten.Shortener) shorten.Shortener {
	// `package metrics` domain
	requestsShorten := metrics.NewMultiCounter(
		expvar.NewCounter("requests_shorten"),
		statsd.NewCounter(ioutil.Discard, "requests_shorten_total", time.Second),
		prometheus.NewCounter(stdprometheus.CounterOpts{
			Namespace: "shortener",
			Subsystem: "shorten",
			Name:      "requests_total",
			Help:      "Total number of received shorten requests.",
		}, []string{}),
	)
	durationShorten := metrics.NewTimeHistogram(time.Nanosecond, metrics.NewMultiHistogram(
		expvar.NewHistogram("duration_shorten_nanoseconds_total", 0, 1e9, 3, 50, 95, 99),
		statsd.NewHistogram(ioutil.Discard, "duration_shorten_nanoseconds_total", time.Second),
		prometheus.NewSummary(stdprometheus.SummaryOpts{
			Namespace: "shortener",
			Subsystem: "shorten",
			Name:      "duration_nanoseconds_total",
			Help:      "Total nanoseconds spend serving shorten requests.",
		}, []string{}),
	))

	requestsResolve := metrics.NewMultiCounter(
		expvar.NewCounter("requests_resolve"),
		statsd.NewCounter(ioutil.Discard, "requests_resolve_total", time.Second),
		prometheus.NewCounter(stdprometheus.CounterOpts{
			Namespace: "shortener",
			Subsystem: "resolve",
			Name:      "requests_total",
			Help:      "Total number of received resolve requests.",
		}, []string{}),
	)
	durationResolve := metrics.NewTimeHistogram(time.Nanosecond, metrics.NewMultiHistogram(
		expvar.NewHistogram("duration_resolve_nanoseconds_total", 0, 1e9, 3, 50, 95, 99),
		statsd.NewHistogram(ioutil.Discard, "duration_resolve_nanoseconds_total", time.Second),
		prometheus.NewSummary(stdprometheus.SummaryOpts{
			Namespace: "shortener",
			Subsystem: "resolve",
			Name:      "duration_nanoseconds_total",
			Help:      "Total nanoseconds spend serving resolve requests.",
		}, []string{}),
	))

	requestsLatest := metrics.NewMultiCounter(
		expvar.NewCounter("requests_latest"),
		statsd.NewCounter(ioutil.Discard, "requests_latest_total", time.Second),
		prometheus.NewCounter(stdprometheus.CounterOpts{
			Namespace: "shortener",
			Subsystem: "latest",
			Name:      "requests_total",
			Help:      "Total number of received latest requests.",
		}, []string{}),
	)
	durationLatest := metrics.NewTimeHistogram(time.Nanosecond, metrics.NewMultiHistogram(
		expvar.NewHistogram("duration_latest_nanoseconds_total", 0, 1e9, 3, 50, 95, 99),
		statsd.NewHistogram(ioutil.Discard, "duration_latest_nanoseconds_total", time.Second),
		prometheus.NewSummary(stdprometheus.SummaryOpts{
			Namespace: "shortener",
			Subsystem: "latest",
			Name:      "duration_nanoseconds_total",
			Help:      "Total nanoseconds spend serving latest requests.",
		}, []string{}),
	))

	return &instrumentedShortener{
		requestsShorten: requestsShorten,
		requestsResolve: requestsResolve,
		requestsLatest:  requestsLatest,

		durationShorten: durationShorten,
		durationResolve: durationResolve,
		durationLatest:  durationLatest,

		shortener: s,
	}
}

type instrumentedShortener struct {
	requestsShorten metrics.Counter
	requestsResolve metrics.Counter
	requestsLatest  metrics.Counter

	durationShorten metrics.TimeHistogram
	durationResolve metrics.TimeHistogram
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

func (s *instrumentedShortener) Latest(ctx context.Context, count int64) ([]*pb.ShortURL, error) {
	defer func(begin time.Time) {
		s.requestsLatest.Add(1)
		s.durationLatest.Observe(time.Since(begin))
	}(time.Now())
	return s.shortener.Latest(ctx, count)
}
