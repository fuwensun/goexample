package dao

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fuwensun/goms/eTest/internal/model"
	"github.com/fuwensun/goms/pkg/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
)

// dbcfg mysql config.
type dbcfg struct {
	DSN string `yaml:"dsn"`
}

//
type cccfg struct {
	Name string
	Addr string `yaml:"addr"`
}

// Dao dao interface
type Dao interface {
	Close()

	Ping(ctx context.Context) (err error)
	//count
	UpdatePingCount(c context.Context, t model.PingType, v model.PingCount) error
	ReadPingCount(c context.Context, t model.PingType) (model.PingCount, error)
	//user-cc
	ExistUserCC(c context.Context, uid int64) (bool, error)
	SetUserCC(c context.Context, user *model.User) error
	GetUserCC(c context.Context, uid int64) (model.User, error)
	DelUserCC(c context.Context, uid int64) error
	//user-db
	CreateUserDB(c context.Context, user *model.User) error
	UpdateUserDB(c context.Context, user *model.User) error
	ReadUserDB(c context.Context, uid int64) (user model.User, err error)
	DeleteUserDB(c context.Context, uid int64) error
	//user
	CreateUser(c context.Context, user *model.User) error
	UpdateUser(c context.Context, user *model.User) error
	ReadUser(c context.Context, uid int64) (model.User, error)
	DeleteUser(c context.Context, uid int64) error
}

// dao dao.
type dao struct {
	db    *sql.DB
	redis redis.Conn
}

func getDBConfig(cfgpath string) (dbcfg, error) {
	var df dbcfg
	path := filepath.Join(cfgpath, "mysql.yml")
	if err := conf.GetConf(path, &df); err != nil {
		log.Printf("get db config file: %v", err)
	}
	if df.DSN != "" {
		log.Printf("get config db DSN: %v", df.DSN)
	}
	if dsn := os.Getenv("MYSQL_SVC_DSN"); dsn != "" {
		df.DSN = dsn
		log.Printf("get env db DSN: %v", dsn)
	}
	if df.DSN == "" {
		err := fmt.Errorf("get db DSN: %w", model.ErrNotFoundData)
		return df, err
	}
	return df, nil
}
func getCCConfig(cfgpath string) (cccfg, error) {
	var cf cccfg
	path := filepath.Join(cfgpath, "redis.yml")
	if err := conf.GetConf(path, &cf); err != nil {
		log.Printf("get cc config file: %v", err)
	}
	if cf.Addr != "" {
		log.Printf("get config cc Addr: %v", cf.Addr)
	}
	if addr := os.Getenv("REDIS_SVC_ADDR"); addr != "" {
		cf.Addr = addr
		log.Printf("get env cc Addr: %v", addr)
	}
	if cf.Addr == "" {
		err := fmt.Errorf("get cc Addr: %w", model.ErrNotFoundData)
		return cf, err
	}
	return cf, nil
}

// New new a dao.
func New(cfgpath string) (Dao, func(), error) {
	//db
	df, err := getDBConfig(cfgpath)
	if err != nil {
		return nil, nil, err //?
	}
	mdb, err := sql.Open("mysql", df.DSN)
	if err != nil {
		log.Panicf("open db: %v", err)
	}
	if err := mdb.Ping(); err != nil {
		log.Panicf("ping db: %v", err)
	}
	//cc
	cf, err := getCCConfig(cfgpath)
	if err != nil {
		return nil, nil, err //?
	}
	mcc, err := redis.Dial("tcp", cf.Addr)
	if err != nil {
		log.Panicf("dial cc: %v", err)
	}
	if _, err := mcc.Do("PING"); err != nil {
		log.Panicf("ping cc: %v", err)
	}
	//
	d := &dao{
		db:    mdb,
		redis: mcc,
	}
	return d, d.Close, nil
}

// Close close the resource.
func (d *dao) Close() {
	d.redis.Close()
	d.db.Close()
}

// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	if _, err = d.redis.Do("PING"); err != nil {
		return
	}
	return d.db.PingContext(ctx)
}
