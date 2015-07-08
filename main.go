package main

import (
	// import expvar for it's side effect of registering a handler at /debug/vars
	_ "expvar"
	"flag"
	"fmt"
	stdlog "log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/redis.v3"

	"github.com/DanielHeckrath/short/logging"
	"github.com/DanielHeckrath/short/metrics"
	"github.com/DanielHeckrath/short/pb"
	"github.com/DanielHeckrath/short/shorten"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

const httpAddr = ":80"
const debugAddr = ":8000"
const grpcAddr = ":8001"

var defaultTimestampUTCNano kitlog.Valuer = func() interface{} {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func main() {
	// Flag domain. Note that gRPC transitively registers flags via its import
	// of glog. So, we define a new flag set, to keep those domains distinct.
	fs := flag.NewFlagSet("", flag.ExitOnError)
	var (
		redisPass = fs.String("redis.pass", "", "Redis server password")
		redisDB   = fs.Int64("redis.db", 0, "Redis server database")

		shortenProto    = fs.String("short.proto", "http", "Protocol to use for short urls")
		shortenHost     = fs.String("short.host", "localhost", "Host to use for short urls")
		shortenNotFound = fs.String("short.na", "/", "Redirect address in case no url was found")
	)
	flag.Usage = fs.Usage // only show our flags
	fs.Parse(os.Args[1:])

	// `package log` domain
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.With(logger, "ts", defaultTimestampUTCNano)
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger)) // redirect stdlib logging to us
	stdlog.SetFlags(0)                                // flags are handled in our logger

	// Server domain

	// Shortener needs a redis connection
	redisHost := os.Getenv("REDIS_PORT_6379_TCP_ADDR")
	redisPort := os.Getenv("REDIS_PORT_6379_TCP_PORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	logger.Log("srv", "redis", "addr", redisAddr, "pw", *redisPass, "db", *redisDB)

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: *redisPass, // no password set
		DB:       *redisDB,   // use default DB
	})
	defer client.Close()

	res, err := client.Ping().Result()

	if err != nil {
		logger.Log("fatal", err)
		os.Exit(1)
	}

	logger.Log("srv", "redis", "ping", res)

	logger.Log("srv", "shorten", "proto", *shortenProto, "host", *shortenHost)
	var s = shorten.New(*shortenProto, *shortenHost, client)
	s = metrics.Instrument(s)
	s = logging.Log(logger, s)

	shortenEndpoint := shorten.NewShortenEndpoint(s)
	resolveEndpoint := shorten.NewResolveEndpoint(s)
	infoEndpoint := shorten.NewInfoEndpoint(s)
	latestEndpoint := shorten.NewLatestEndpoint(s)

	// Mechanical stuff
	rand.Seed(time.Now().UnixNano())
	root := context.Background()
	errc := make(chan error)

	go func() {
		errc <- interrupt()
	}()

	// Transport: HTTP (debug/instrumentation)
	go func() {
		logger.Log("addr", debugAddr, "transport", "debug")
		http.Handle("/metrics", prometheus.Handler())
		errc <- http.ListenAndServe(debugAddr, nil)
	}()

	// Transport: HTTP (JSON)
	go func() {
		ctx, cancel := context.WithCancel(root)
		defer cancel()
		before := []httptransport.BeforeFunc{}
		after := []httptransport.AfterFunc{}
		shortenHandler := prometheus.InstrumentHandler(
			"shorten",
			shorten.NewShortenHandler(ctx, shortenEndpoint, before, after),
		)
		resolveHandler := prometheus.InstrumentHandler(
			"resolve",
			shorten.NewResolveHandler(ctx, resolveEndpoint, *shortenNotFound),
		)
		infoHandler := prometheus.InstrumentHandler(
			"info",
			shorten.NewInfoHandler(ctx, infoEndpoint, before, after),
		)
		latestHandler := prometheus.InstrumentHandler(
			"latest",
			shorten.NewLatestHandler(ctx, latestEndpoint, before, after),
		)

		router := mux.NewRouter()
		router.Path("/shorten").Methods("POST").HandlerFunc(shortenHandler)
		router.Path("/{key:([a-zA-Z0-9]+$)}").Methods("GET").HandlerFunc(resolveHandler)
		router.Path("/info/{key:[a-zA-Z0-9]+}").Methods("GET").HandlerFunc(infoHandler)
		router.Path("/latest/{count:[0-9]+}").Methods("GET").HandlerFunc(latestHandler)

		logger.Log("addr", httpAddr, "transport", "HTTP/JSON")
		errc <- http.ListenAndServe(httpAddr, logging.Handler(logger, router))
	}()

	// Transport: gRPC
	go func() {
		ln, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			errc <- err
			return
		}
		s := grpc.NewServer() // uses its own context?
		pb.RegisterShortServer(s, grpcBinding{
			shorten: shortenEndpoint,
			resolve: resolveEndpoint,
			latest:  latestEndpoint,
		})
		logger.Log("addr", grpcAddr, "transport", "gRPC")
		errc <- s.Serve(ln)
	}()

	logger.Log("fatal", <-errc)
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
