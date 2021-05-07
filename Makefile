buid:
	go build ./cmd/main.go
run:
	go run ./cmd/main.go
postgres:
	docker run --name psql -p 5432\:5432 -e POSTGRES_USER\=root -e POSTGRES_PASSWORD\=pass -d postgres\:13-alpine
psqlstart:
	docker start psql
psqlstop:
	docker stop psql
createdb:
	docker exec -it psql createdb --username=root --owner=root usr
dropdb:
	docker exec -it psql dropdb usr
migrateup:
	migrate -path schema -database "postgresql://root:pass@localhost:5432/usr?sslmode=disable" -verbose up
migratedown:
	migrate -path schema -database "postgresql://root:pass@localhost:5432/usr?sslmode=disable" -verbose down
redis:	
	docker run --name usr-redis -p 6379:6379 -d redis
.PHONY: run, buid

.DEFAULT_GOAL := buid