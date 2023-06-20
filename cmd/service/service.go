package service

import (
	"flag"
	"fmt"
	ep "kitprc/endpoint"
	"kitprc/repository"
	"kitprc/service"
	"kitprc/tool"
	"kitprc/transport/grpc"
	"kitprc/transport/grpc/pb"
	"kitprc/transport/http"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	netHttp "net/http"

	mw "kitprc/middleware"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/juju/ratelimit"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/log/level"
	googleGrpc "google.golang.org/grpc"

	"github.com/oklog/oklog/pkg/group"
)

var logger log.Logger

var (
	fs       = flag.NewFlagSet("world", flag.ExitOnError)
	httpAddr = fs.String("http-addr", ":7001", "Http listen address")
	grpcAddr = fs.String("grpc-addr", ":7002", "Http listen address")
)

func Run() {

	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	store := repository.New()
	svc := service.New(store)

	ems := mw.TokenBucketLimitter(mw.Bucket{ratelimit.NewBucket(time.Hour, 2)})

	eps := ep.NewEndpoint(svc, ems)

	g := &group.Group{}
	initHttpHandler(eps, g)
	initGRPCHandler(eps, g)
	initCancelInterrupt(g)

	_ = level.Error(logger).Log("exit", g.Run())

}

func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}

func initHttpHandler(endpoints ep.Endpoints, g *group.Group) {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(level.Error(logger))),
	}

	httpHandler := http.MakeHTTPHandler(endpoints, opts...)

	// 注册
	go func() {
		err := tool.RegService("127.0.0.1:8500", "1", "测试1", "127.0.0.1", 8000, "5s", "http://172.16.8.232:8000/health", "test")
		if err != nil {
			_ = level.Error(logger).Log("transport", "HTTP", "during", "Reg", "err", err)
		}
		_ = netHttp.ListenAndServe("0.0.0.0:8000", httpHandler)
	}()
}

func initGRPCHandler(endpoints ep.Endpoints, g *group.Group) {
	grpcOpts := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		_ = logger.Log("transport", "gRPC", "during", "Listen", "err", err)
	}

	g.Add(func() error {
		baseServer := googleGrpc.NewServer()
		pb.RegisterServiceServer(baseServer, grpc.MakeGRPCHandler(endpoints, grpcOpts...))
		return baseServer.Serve(grpcListener)
	}, func(error) {
		_ = level.Error(logger).Log("grpcListener.Close", grpcListener.Close())
	})
}
