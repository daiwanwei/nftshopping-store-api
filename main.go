package main

import (
	"context"
	"flag"
	adapterRouters "nftshopping-store-api/adapter/routers"
	_ "nftshopping-store-api/docs"
	eventRouters "nftshopping-store-api/event/routers"
	"nftshopping-store-api/pkg/config"
	"nftshopping-store-api/pkg/flags"
	"nftshopping-store-api/pkg/log"
	"strconv"
	"sync"
)

// @title Swagger API
// @version 1.0
// @description Swagger API for Golang Project Blueprint.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email alienwade007@gmail.com

// @license.name MIT
// @license.url https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE
// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @BasePath /
func main() {
	flag.StringVar(&flags.Env, "env", "local", "environment")
	flag.Parse()
	c, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	logger, err := log.GetLog()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	//sub,err:=subscribers.NewItemSubscriber()
	//if err != nil {
	//	return
	//}
	//go func(){
	//	defer wg.Done()
	//	err:=sub.Run()
	//	if err != nil {
	//		log.Log.Error(err)
	//	}
	//}()
	eventRouter, err := eventRouters.GetRouter()
	if err != nil {
		logger.Error(err)
		return
	}
	go func() {
		defer wg.Done()
		err := eventRouter.Run(context.Background())
		if err != nil {
			logger.Error(err)
		}
	}()
	adapterRouter, err := adapterRouters.GetRouter()
	if err != nil {
		logger.Error(err)
		return
	}
	go func() {
		defer wg.Done()
		err := adapterRouter.Run(":" + strconv.Itoa(c.Server.Port))
		if err != nil {
			logger.Error(err)
		}
	}()
	wg.Wait()
}
