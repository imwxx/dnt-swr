package ginit

import (
    "flag"
    "fmt"
    "dnt-swr/lib"
)


type Flags struct {
    ip   string
    port string
    file string
}

const (
    ip   = "0.0.0.0"
    port = "9999"
    file = "./conf/appconf.yml"
)

func parseFlags() *Flags {
    argsflag := new(Flags)
    flag.StringVar(&argsflag.ip, "i", ip, "bond listen ip, default 0.0.0.0")
    flag.StringVar(&argsflag.port, "p", port, "bond listen port, default 9999")
    flag.StringVar(&argsflag.file, "f", file, "config file, default ./conf/appconf.yml")
    flag.Parse()
    return argsflag
}

func InitServer() string {
    argsflag := parseFlags()
    info := argsflag.ip + ":" + argsflag.port
    file := argsflag.file
    conf, err := lib.LoadConfig(file)
    if err != nil {
        fmt.Println(err)
    } else {
        lib.AppConf = conf
    }

    // redis
    redisInfo := conf.REDISHOST + ":" + conf.REDISPORT
    lib.RedisClient, err = lib.InitRedis(redisInfo)

    // mysql
    db, err := lib.InitDb(conf.DBINFO)
    if err != nil {
        fmt.Printf("Init db falied: %s", err)
    } else {
        lib.MDBC = db
    }

    return info
}

