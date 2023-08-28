include config.env
export

run:
	docker-compose up --build -d

stop:
	docker-compose down -v
	
conn:
	PGPASSWORD=$$PG_PASSWORD psql -h 0.0.0.0 -p 54344 -U $$PG_USER -d $$PG_DB -c '\q';