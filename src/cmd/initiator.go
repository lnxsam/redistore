package main

import (
	"fmt"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redistore/internal/api/rest"
	"redistore/internal/domain/creating"
	"redistore/internal/domain/listing"
	"redistore/internal/domain/searching"
	"redistore/internal/domain/updating"
	"strconv"

	"redistore/pkg/configs"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	db           *gorm.DB
	redisClient  *redis.Client
	searchEngine *redisearch.Client
)

func provideDB() *gorm.DB {
	if db == nil {
		var err error
		switch configs.Env("DB_DRIVER") {
		case "postgres":
			db, err = gorm.Open(
				postgres.Open(
					"host="+configs.Env("POSTGRES_HOST")+
						" user="+configs.Env("POSTGRES_USER")+
						" password="+configs.Env("POSTGRES_PASSWORD")+
						" dbname="+configs.Env("POSTGRES_DB_NAME")+
						" port="+configs.Env("POSTGRES_DB_PORT")+
						" sslmode=disable",
				), &gorm.Config{
					NamingStrategy: schema.NamingStrategy{
						SingularTable: true,
					},
					SkipDefaultTransaction: true,
					PrepareStmt:            true,
				},
			)
			if err != nil {
				panic("failed to connect database")
			}
			break
		default:
			panic("please choose valid db name")
		}
	}
	return db
}

func provideSearchEngine() *redisearch.Client {
	if searchEngine == nil {

		// Create a client. By default a client is schemaless
		// unless a schema is provided when creating the index
		searchEngine = redisearch.NewClient("localhost:6379", "redistore_index")
		// Create a schema
		sc := redisearch.NewSchema(redisearch.DefaultOptions).
			AddField(redisearch.NewTextFieldOptions("Title", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true}))

		// Drop an existing index. If the index does not exist an error is returned
		searchEngine.Drop()

		// Create the index with the given schema
		if err := searchEngine.CreateIndex(sc); err != nil {
			log.Fatal(err)
		}
	}
	return searchEngine
}

func provideCache() *redis.Client {
	if redisClient == nil {
		redisDB, err := strconv.Atoi(configs.Env("REDIS_DB"))
		if err != nil {
			panic("invalid redis DB")
		}
		redisClient = redis.NewClient(&redis.Options{
			Addr:     configs.Env("REDIS_HOST") + ":" + configs.Env("REDIS_PORT"),
			Password: configs.Env("REDIS_PASSWORD"),
			DB:       redisDB,
		})
	}
	return redisClient
}

func loadConfigFile() {
	err := godotenv.Load("./../.env")
	if err != nil {
		panic(err)
		// panic("can not find env file")
	}
}

func startRestServer(creatingSvc creating.Service, updatingSvc updating.Service, searchingSvc searching.Service, listingSvc listing.Service) {
	handler := rest.New(creatingSvc, updatingSvc, searchingSvc, listingSvc)
	router := gin.New()
	router.POST("/create_product", handler.CreateProduct)
	router.POST("/products", handler.GetProductList)
	router.POST("/search_products_by_title", handler.SearchProductsByTitle)
	router.POST("/create_card", handler.CreateCard)
	router.POST("/add_products_to_card", handler.AddProductToCard)
	router.POST("/remove_card_item", handler.RemoveCardItem)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", "8081"),
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()
}
