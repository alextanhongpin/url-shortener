-include .env
export

start:
	@go run main.go

install:
	@go get -u github.com/pressly/goose/cmd/goose

migrate:
	@goose -dir migrations mysql "${DB_USER}:${DB_PASS}@tcp(127.0.0.1:3306)/${DB_NAME}?parseTime=true" up

mysql:
	mysql -h 127.0.0.1 -u ${DB_USER} -p ${DB_NAME}
