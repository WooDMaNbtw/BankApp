DB_URL = postgresql://bank_admin:bank_password_qwertyuiop@bankapp.cpomksa8kmhk.eu-north-1.rds.amazonaws.com:5432/bank_app

docker_up:
	docker compose --env-file .env.docker up -d

docker_down:
	docker compose --env-file .env.docker down

create_db:
	docker exec -it pgdb-bank createdb --username=bank_admin --owner=bank_admin bank_app

drop_db:
	docker exec -it pgdb-bank dropdb --username=bank_admin bank_app

migrate_up:
	migrate -path db/migrations -database "${DB_URL}" -verbose up

migrate_up1:
	migrate -path db/migrations -database "${DB_URL}" -verbose up 1

migrate_down:
	migrate -path db/migrations -database "${DB_URL}" -verbose down

migrate_down1:
	migrate -path db/migrations -database "${DB_URL}" -verbose down 1

sqlc:
	sqlc -f db/sqlc.yaml generate

db_docs:
	dbdocs build docs/db.html

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbm

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/WooDMaNbtw/BankApp/db/sqlc Store

.PHONY: docker_up docker_down create_db drop_db migrate_up migrate_up1 migrate_down migrate_down1 sqlc db_docs db_schema test server mock
