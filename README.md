<!-- go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/swaggo/swag/cmd/swag@latest -->
make setup

API SWAGGER: [http://localhost:8080/swagger/]
Token swagger login: 'Bearer TOKEN"
Create migrate: migrate create -ext sql -dir db/migrations -seq create_name_table