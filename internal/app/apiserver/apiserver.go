package apiserver

import (

	//	"fmt"

	"context"
	"fmt"
	"log"
	"net/http"
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

	fmt.Printf("redis is ok %v\n", redis)

	srv := newServer(store, sessionStore, redis_store)
	/*
		srvv := &http.Server{
			Addr: config.BindAddr,
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
		}
	*/
	//	return http.ListenAndServe(config.BindAddr, srv)
	// http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./web/images"))))
	//return http.ListenAndServe(config.BindAddr, srv)
	serv := http.ListenAndServe(config.BindAddr, srv)

	return serv
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
