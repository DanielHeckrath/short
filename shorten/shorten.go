package shorten

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/redis.v3"

	"github.com/DanielHeckrath/short/pb"
	"github.com/juju/errors"
)

const counter = "__counter__"

// Shortener is the fundamental interface for all shorterning operations
type Shortener interface {
	Shorten(context.Context, string) (*pb.ShortURL, error)
	Resolve(context.Context, string) (*pb.ShortURL, error)
	Info(context.Context, string) (*pb.ShortURL, error)
	Latest(context.Context, int64) ([]*pb.ShortURL, error)
}

// New return a concrete implementation of a shorten function
func New(proto, host string, client *redis.Client) Shortener {
	return &shortener{
		redis: client,
		proto: proto,
		host:  host,
	}
}

type shortener struct {
	redis *redis.Client
	proto string
	host  string
}

func (s *shortener) store(key, shortURL, longURL string) (*pb.ShortURL, error) {
	url := newShortURL(key, shortURL, longURL)
	status := s.redis.HMSet(url.Key,
		"LongUrl", url.LongUrl,
		"ShortUrl", url.ShortUrl,
		"CreationDate", fmt.Sprint(url.CreationDate),
		"Clicks", fmt.Sprint(url.Clicks),
	)

	if status.Err() != nil {
		return nil, status.Err()
	}

	return url, nil
}

func (s *shortener) load(key string) (*pb.ShortURL, error) {
	if ok, _ := s.redis.HExists(key, "ShortUrl").Result(); !ok {
		return nil, errors.New("Unknown key: " + key)
	}

	url := &pb.ShortURL{}
	url.Key = key

	res, err := s.redis.HMGet(key, "LongUrl", "ShortUrl", "CreationDate", "Clicks").Result()

	if err != nil {
		return nil, errors.Annotatef(err, "Unable to load values for key: %s", key)
	}

	url.LongUrl = res[0].(string)
	url.ShortUrl = res[1].(string)

	date := res[2].(string)
	ts, err := strconv.ParseInt(date, 10, 64)

	if err != nil {
		return nil, errors.Annotatef(err, "Unable to convert timestamp for key: %s", key)
	}

	url.CreationDate = ts

	var clicks = res[3].(string)
	count, err := strconv.ParseInt(clicks, 10, 64)

	if err != nil {
		return nil, errors.Annotatef(err, "Unable to convert clicks for key: %s", key)
	}

	url.Clicks = count
	return url, nil
}

func (s *shortener) Shorten(ctx context.Context, longURL string) (*pb.ShortURL, error) {
	url, err := url.Parse(longURL)

	if err != nil {
		return nil, errors.Annotate(err, "Unable to parse input url")
	}

	count, err := s.redis.Incr(counter).Result()

	if err != nil {
		return nil, errors.Annotate(err, "Unable to increase link counter")
	}

	encoded := encode(count)
	location := fmt.Sprintf("%s://%s/%s", s.proto, s.host, encoded)

	short, err := s.store(encoded, location, url.String())

	if err != nil {
		return nil, errors.Annotate(err, "Unable to save short url")
	}

	return short, nil
}

func (s *shortener) Resolve(ctx context.Context, key string) (*pb.ShortURL, error) {
	url, err := s.load(key)

	if err != nil {
		return nil, errors.Annotatef(err, "Unable to load short url for key: %s", key)
	}

	_, err = s.redis.HIncrBy(key, "Clicks", 1).Result()

	if err != nil {
		return nil, errors.Annotatef(err, "Unable to increase click count for key: %s", key)
	}

	return url, nil
}

func (s *shortener) Info(ctx context.Context, key string) (*pb.ShortURL, error) {
	url, err := s.load(key)

	if err != nil {
		return nil, errors.Annotatef(err, "Unable to load short url for key: %s", key)
	}

	return url, nil
}

func (s *shortener) Latest(ctx context.Context, count int64) ([]*pb.ShortURL, error) {
	start, err := s.redis.Get(counter).Int64()

	if err != nil {
		return nil, errors.Annotate(err, "Unable to load current link count")
	}

	end := (start - count)

	var urls = make([]*pb.ShortURL, count)

	var i int64
	for ; start-i > end && i < count && start-i > 0; i++ {
		url, err := s.load(encode(start - i))

		if err != nil {
			return nil, errors.Annotate(err, "Unable to load short url")
		}

		urls[i] = url
	}

	return urls[:i], nil
}

// Creates a new KurzUrl instance. The Given key, shorturl and longurl will
// be used. Clicks will be set to 0 and CreationDate to time.Nanoseconds()
func newShortURL(key, shorturl, longurl string) *pb.ShortURL {
	url := &pb.ShortURL{}
	url.CreationDate = time.Now().UnixNano()
	url.Key = key
	url.LongUrl = longurl
	url.ShortUrl = shorturl
	url.Clicks = 0
	return url
}
