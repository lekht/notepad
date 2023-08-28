include config.env
export

app:
	go run ./cmd/service/service.go 

run:
	docker-compose up --build -d

stop:
	docker-compose down -v
	
conn:
	PGPASSWORD=$$PG_PASSWORD psql -h $$PG_HOST -p $$PG_PORT -U $$PG_USER -d $$PG_DB -c '\q';