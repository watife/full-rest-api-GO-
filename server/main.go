package main

import (
	"database/sql"
	"fakorede-bolu/full-rest-api/server/pkg/models/postgres"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gopkg.in/go-playground/validator.v9"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	todo     *postgres.TodoModel
	user     *postgres.UserModel
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
	databasePassword := os.Getenv("DATABASE_PASSWORD")

	dbConn := "postgres://" + databaseUser + ":" + databasePassword + "@localhost/" + databaseName + "?sslmode=disable"

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

	// db connection
	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		todo:     &postgres.TodoModel{DB: db},
		user:     &postgres.UserModel{DB: db},
	}

	// custom server struct to make use of custom errorLog
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Listen to server
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()

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
