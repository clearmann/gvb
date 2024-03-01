package main

import (
	"context"
	"github.com/gin-gonic/gin"
	ginblog "gvb/internal"
	g "gvb/internal/global"
	"gvb/internal/middleware"
	sonwflake "gvb/internal/utils/snowflake"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(r *gin.Engine, serverName string, addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	//保证下面的优雅启停
	go func() {
		log.Printf("%s running in %s \n", serverName, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	//SIGINT 用户发送INTR字符(Ctrl+C)触发
	//SIGTERM 结束程序(可以被捕获、阻塞或忽略)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting Down project %s... \n", serverName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s Shutdown, cause by : %v \n", serverName, err)
	}
	select {
	case <-ctx.Done():
		log.Println("关闭超时")
	}
	log.Printf("%s stop success...", serverName)
}
func main() {
	r := gin.Default()
	conf := g.ReadConfig()
	ginblog.InitLogger(conf)
	db := ginblog.InitDatabase(conf)
	rdb := ginblog.InitRedis(conf)
	r.Use(middleware.CORS())
	r.Use(middleware.WithGormDB(db), middleware.WithRedisDB(rdb))
	r.Use(middleware.WithCookieStore(conf.Session.Name, conf.Session.Salt))
	sonwflake.Init(conf.Server.StartTime, conf.Server.MachineID)
	ginblog.RegisterHandlers(r)
	Run(r, "gvb", conf.Server.Port)
}
