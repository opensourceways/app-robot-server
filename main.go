package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/opensourceways/app-robot-server/config"
	"github.com/opensourceways/app-robot-server/dbmodels"
	"github.com/opensourceways/app-robot-server/logs"
	"github.com/opensourceways/app-robot-server/mongodb"
	"github.com/opensourceways/app-robot-server/router"
)

// @title Swagger app-robot-server API
// @version 0.0.1
// @description plugin maintenance server api doc
// contact.name WeiZhi Xie
// contact.email 986740642@qq.com
// @securityDefinitions.apikey ApiKeyAuth
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @in header
// @name access-token
// @BasePath /v1
func main() {
	//TODO: config file path parse from the args by flag package
	if err := config.InitConfig(); err != nil {
		logs.Logger.Fatal(err)
	}

	if err := logs.Init(); err != nil {
		logs.Logger.Fatal(err)
	}

	db, err := mongodb.Initialize(&config.Application.Mongo)
	if err != nil {
		logs.Logger.Fatal(err)
	}
	dbmodels.RegisterDB(db)
	err = dbmodels.GetDB().InitCUsers()
	if err != nil {
		logs.Logger.Fatal(err)
	}
	runServer()
}

func runServer() {
	route := router.Init()
	address := fmt.Sprintf(":%s", config.Application.Port)
	s := &http.Server{
		Addr:           address,
		Handler:        route,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logs.Logger.Debug("server run success on ", address)
	fmt.Printf("welcome %s api. the default doc address: http://locahost%s/swagger/index.html\n", logs.Module, s.Addr)
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Logger.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logs.Logger.Error("shutdown server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logs.Logger.Fatalf("sever shutdown: %s\n", err)
	}
	logs.Logger.Debug("sever exiting")
}
