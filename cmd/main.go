package main

import (
	"crypto/tls"
	"database/sql"
	"fakorede-bolu/full-rest-api/pkg/models/postgres"

	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gopkg.in/go-playground/validator.v9"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	todo     *postgres.TodoModel
	user     *postgres.UserModel
	inbox    *postgres.InboxModel
}

type ReminderEmails struct {
	// Filtered
}

func (e ReminderEmails) Run() {
	// Queries the DB
	// Sends some email
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

}

var validate *validator.Validate

func main() {

	validate = validator.New()

	databaseName := os.Getenv("DATABASE_NAME")
	databaseUser := os.Getenv("DATABASE_USER")
	databaseHost := os.Getenv("DATABASE_HOST")
	// databasePort := os.Getenv("DATABASE_PORT")
	databasePassword := os.Getenv("DATABASE_PASSWORD")

	// dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, databasePort, databaseUser, databaseName, databasePassword)
	dbConn := "postgres://" + databaseUser + ":" + databasePassword + "@" + databaseHost + ":5432/" + databaseName + "?sslmode=disable"

	fmt.Println(dbConn)

	// dynamic http address from command-line flag with default :4000
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", dbConn, "MySQL data source name")

	// parse the command-line flag
	flag.Parse()

	// Info log
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Error log
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Lshortfile) //can also tage log.Llongfile for full path

	// pg db connection
	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// pool := &redis.Pool{
	// 	MaxIdle:   80,
	// 	MaxActive: 12000,
	// 	Dial: func() (redis.Conn, error) {
	// 		conn, err := redis.Dial("tcp", "localhost:6379")
	// 		if err != nil {
	// 			log.Printf("ERROR: fail init redis pool: %s", err.Error())
	// 			os.Exit(1)
	// 		}
	// 		return conn, err
	// 	},
	// }

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		todo:     &postgres.TodoModel{DB: db},
		user:     &postgres.UserModel{DB: db},
		inbox:    &postgres.InboxModel{DB: db},
	}

	// outbox()

	go app.cronn()

	// custom server struct to make use of custom errorLog
	srv := &http.Server{
		Addr:           *addr,
		ErrorLog:       errorLog,
		Handler:        app.routes(),
		TLSConfig:      tlsConfig,
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 524288,
	}

	// Listen to server
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	// err = srv.ListenAndServe()

	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func task() {
	fmt.Println("I am runnning task.")
}

func outbox() {
	gocron.Every(1).Second().Do(task)
	_, time := gocron.NextRun()
	fmt.Println(time)
	fmt.Println("I am runnning task.")

	<-gocron.Start()
}

// func (app *application) cronn() {
// 	fmt.Println("Job starts")

// }
