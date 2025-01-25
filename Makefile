start-app:
	docker-compose up -d
	timeout /t 2
	migrate -path Server/schema/migrations -database 'postgres://postgres:1111@localhost:8080/ChatDB?sslmode=disable' up

migrate-up:
	migrate -path Server/schema/migrations -database 'postgres://postgres:1111@localhost:8080/ChatDB?sslmode=disable' up

migrate-down:
	migrate -path Server/schema/migrations -database 'postgres://postgres:1111@localhost:8080/ChatDB?sslmode=disable' down

migrate-down-1:
	migrate -path Server/schema/migrations -database 'postgres://postgres:1111@localhost:8080/ChatDB?sslmode=disable' down 1

stop-app:
	docker-compose down -v