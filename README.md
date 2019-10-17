A basic Implementation of full rest api with golang

Create a .env file at the root of the project with key: 

DATABASE_NAME 
DATABASE_USER
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_PASSWORD 
JWT_KEY 

run:   go run ./server -addr=":9999"  (where :9999 can be any address)

Running with Docker

Create a .env file having 

DATABASE_NAME
DATABASE_USER
DATABASE_HOST=fullstack-postgres
DATABASE_PORT=5432
DATABASE_PASSWORD
JWT_KEY
GMAIL_USERNAME (If you want to use the mail function)
GMAIL_PASSWORD  (If you want to use the mail function)

Docker compose up --build