package main

import (
	"os"
	"os/signal"
	"redistore/internal/data/datasource/redisearch"
	"redistore/internal/domain/creating"
	"redistore/internal/domain/listing"
	"redistore/internal/domain/searching"
	"redistore/internal/domain/updating"

	"redistore/internal/data"
	"redistore/internal/data/datasource/postgres"
	"redistore/internal/data/datasource/redis"
)

func main() {

	// third parties
	loadConfigFile()
	pgDB := provideDB()
	cache := provideCache()
	searchEngine := provideSearchEngine()

	// data_sources
	pgDS := postgres.NewDBDataSource(pgDB)
	cacheDs := redis.NewCacheDataSource(cache)
	searchEngineDs := redisearch.NewSearchDataSource(searchEngine)

	err := pgDS.AutoMigrate()
	if err != nil {
		panic(err)
	}

	// data
	accRepo := data.NewRepository(pgDS, cacheDs, searchEngineDs)

	// domain
	creatingSvc := creating.New(accRepo)
	updatingSvc := updating.New(accRepo)
	searchingSvc := searching.New(accRepo)
	listingSvc := listing.New(accRepo)

	// api
	startRestServer(creatingSvc, updatingSvc, searchingSvc, listingSvc)

	//
	//grpcServer := grpc.GetInstance(server)
	//grpcServer.Start()
	//defer grpcServer.Stop()

	// ctrl c
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	//block until signal is received
	<-ch
}
