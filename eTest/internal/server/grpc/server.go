package grpc

import (
	"context"
	"log"
	"net"
	"path/filepath"

	"github.com/fuwensun/goms/eTest/api"
	"github.com/fuwensun/goms/eTest/internal/model"
	"github.com/fuwensun/goms/eTest/internal/service"
	"github.com/fuwensun/goms/pkg/conf"

	"google.golang.org/grpc"
	xrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	cfgfile = "grpc.yml"
	addr    = ":50051"
)

type ServerCfg struct {
	Addr string `yaml:"addr"`
}

//
type Server struct {
	gs  *grpc.Server
	svc service.Svc
}

//
func New(cfgpath string, s service.Svc) (*Server, error) {
	var sc ServerCfg
	path := filepath.Join(cfgpath, cfgfile)
	if err := conf.GetConf(path, &sc); err != nil {
		log.Printf("get config file: %v", err)
		// err = fmt.Errorf("get config file: %w", err)
		// return nil, err
	}
	if sc.Addr != "" {
		addr = sc.Addr
	}
	log.Printf("grpc server addr: %v", addr)

	gs := xrpc.NewServer()
	server := &Server{svc: s, gs: gs}
	api.RegisterUserServer(gs, server)
	reflection.Register(gs)
	server.start()
	return server, nil
}

func (s *Server) start() {
	gs := s.gs
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panicf("tcp listen: %v", err)
		// fmt.Errorf("tcp listen: %w", err)
		// return server, err
	}
	go func() {
		if err := gs.Serve(lis); err != nil {
			log.Panicf("failed to serve: %v", err)
		}
	}()
}

var pingcount model.PingCount

// example for grpc request handler.
func (s *Server) Ping(ctx context.Context, req *api.Request) (res *api.Reply, err error) {
	svc := s.svc
	message := "pong" + " " + req.Message
	res = &api.Reply{Message: message}
	log.Printf("grpc" + " " + message)
	pingcount++
	svc.UpdateGrpcPingCount(ctx, pingcount)
	pc := svc.ReadGrpcPingCount(ctx)
	log.Printf("grpc ping count: %v\n", pc)
	return res, nil
}
