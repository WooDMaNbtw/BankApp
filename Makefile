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

new_migration:
	migrate create -ext sql -dir db/migrations -seq ${name}

sqlc:
	sqlc -f db/sqlc.yaml generate

db_docs:
	dbdocs build docs/db.html

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/WooDMaNbtw/BankApp/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/WooDMaNbtw/BankApp/worker TaskDistributor


proto:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bank_app \
	--experimental_allow_proto3_optional \
	proto/*.proto
	statik -src=./docs/swagger -dest=./docs

evans:
	evans --host localhost --port 9090 -r repl --package pb

redis:
	docker run --name redis -p 6375:6379 -d redis:7-alpine

ping_redis:
	docker exec -it redis redis-cli ping


.PHONY: docker_up docker_down create_db drop_db migrate_up migrate_up1 migrate_down migrate_down1 sqlc db_docs db_schema test server mock proto redis ping_redis new_migration
