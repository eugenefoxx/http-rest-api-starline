package apiserver

import (

	//	"fmt"

	"context"
	"flag"
	"github.com/eugenefoxx/http-rest-api-starline/pkg/logging"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	//	"os"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store/sqlstore"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store_redis/redisstore"
	"github.com/go-redis/redis/v8"

	///"github.com/go-redis/redis"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

// Start ...
func Start(config *Config) error {
	logger := logging.GetLogger()
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	//	defer db.Close()
	defer db.Close()
	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	//html := config.HTML

	//	router := mux.NewRouter()
	//	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	//srv := newServer(store, sessionStore, html)

	redis, err := newRedis(config.Redis)
	if err != nil {
		//	log.Fatalf("Could not initialize Redis client %s", err)
		log.Printf("Could not initialize Redis client %s", err)

	}
	//defer redis.Close()

	redis_store := redisstore.New(redis)

	//fmt.Printf("redis is run  %v\n", redis)

	srv := newServer(store, sessionStore, redis_store)

	//return http.ListenAndServe(config.BindAddr, srv)

	//servv := http.ListenAndServe(config.BindAddr, srv)
	server := &http.Server{
		Addr: config.BindAddr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      srv, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Errorf(err.Error())
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	logger.Error("shutting down")
	os.Exit(0)

	return nil
}

func newDB(databaseURL string) (*sqlx.DB, error) {
	//func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	//	db, err := sql.Open("pgx", databaseURL)
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
	//		os.Exit(1)
	//	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

/*
type Client struct {
	client *redis.Client
}
*/
func newRedis(redisAddr string) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:        redisAddr,
		DB:          0,
		DialTimeout: 100 * time.Millisecond,
		ReadTimeout: 100 * time.Millisecond,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	/*
		if _, err := client.Ping().Result(); err != nil {
			return nil, err
		}
	*/
	/*return &Client{
		client: client,
	}, nil
	*/
	return client, nil
}
