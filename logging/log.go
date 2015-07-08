package logging

import (
	"time"

	"golang.org/x/net/context"

	"github.com/DanielHeckrath/short/pb"
	"github.com/DanielHeckrath/short/shorten"
	"github.com/go-kit/kit/log"
)

// Log wraps a shorten function with logging
func Log(logger log.Logger, s shorten.Shortener) shorten.Shortener {
	return &loggingShortener{
		logger:    logger,
		shortener: s,
	}
}

type loggingShortener struct {
	logger    log.Logger
	shortener shorten.Shortener
}

func (s *loggingShortener) Shorten(ctx context.Context, url string) (res *pb.ShortURL, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Shorten", "url", url, "result", res, "took", time.Since(begin))
	}(time.Now())
	res, err = s.shortener.Shorten(ctx, url)
	return
}

func (s *loggingShortener) Resolve(ctx context.Context, key string) (res *pb.ShortURL, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Resolve", "key", key, "result", res, "took", time.Since(begin))
	}(time.Now())
	res, err = s.shortener.Resolve(ctx, key)
	return
}

func (s *loggingShortener) Info(ctx context.Context, key string) (res *pb.ShortURL, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Info", "key", key, "result", res, "took", time.Since(begin))
	}(time.Now())
	res, err = s.shortener.Info(ctx, key)
	return
}

func (s *loggingShortener) Latest(ctx context.Context, count int64) (res []*pb.ShortURL, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Latest", "count", count, "result", res, "took", time.Since(begin))
	}(time.Now())
	res, err = s.shortener.Latest(ctx, count)
	return
}
