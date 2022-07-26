package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/navisot/go-url-shortener/api"
	"github.com/navisot/go-url-shortener/repository/mongo"
	"github.com/navisot/go-url-shortener/repository/redis"
	"github.com/navisot/go-url-shortener/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	// Choose repo
	repo := chooseRepo()
	// Init service
	service := shortener.NewRedirectService(repo)
	// Init handler
	handler := api.NewHandler(service)

	// Init router and utilize some middlewares
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Endpoints
	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)

	go func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(httpPort(), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

// httpPort returns the port that app will listen to
func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

// chooseRepo returns a repository based on env var
func chooseRepo() shortener.RedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisUrl := os.Getenv("REDIS_URL")
		repo, err := redis.NewRedisRepository(redisUrl)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoUrl := os.Getenv("MONGO_URL")
		mongoDb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mongo.NewMongoRepository(mongoUrl, mongoDb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}
