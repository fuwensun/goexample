package http

import (
	"log"
	"path/filepath"

	"github.com/fuwensun/goms/eTest/internal/service"
	"github.com/fuwensun/goms/pkg/conf"

	"github.com/gin-gonic/gin"
)

type httpcfg struct {
	Addr string `yaml:"addr"`
}

type Server struct {
	cfg *httpcfg
	eng *gin.Engine
	svc service.Svc
}

func getHttpConfig(cfgpath string) (httpcfg, error) {
	var cfg httpcfg
	path := filepath.Join(cfgpath, "http.yml")
	if err := conf.GetConf(path, &cfg); err != nil {
		log.Printf("get config file: %v", err)
	}
	if cfg.Addr != "" {
		log.Printf("get config addr: %v", cfg.Addr)
	}
	//todo get env
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
		log.Printf("use default addr: %v", cfg.Addr)
	}
	log.Printf("http server addr: %v", cfg.Addr)
	return cfg, nil
}

//
func New(cfgpath string, s service.Svc) (*Server, error) {
	cfg, err := getHttpConfig(cfgpath)
	if err != nil {
		return nil, err
	}
	engine := gin.Default()
	server := &Server{cfg: &cfg, eng: engine, svc: s}
	server.initRouter()
	return server, nil
}

//
func (srv *Server) initRouter() {
	e := srv.eng
	e.GET("/ping", srv.ping)
	user := e.Group("/user")
	{
		user.POST("", srv.createUser)
		user.PUT("/:uid", srv.updateUser)
		user.GET("/:uid", srv.readUser)
		user.DELETE("/:uid", srv.deleteUser)
		user.GET("", srv.readUser)
	}
}

//
func (srv *Server) Start() {
	go func() {
		if err := srv.eng.Run(srv.cfg.Addr); err != nil {
			log.Panicf("failed to server: %v", err)
		}
	}()
}
