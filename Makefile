docker_up:
	docker compose --env-file .env.docker up -d

docker_down:
	docker compose --env-file .env.docker down

create_db:
	docker exec -it pgdb-bank createdb --username=bank_admin --owner=bank_admin bank_app

drop_db:
	docker exec -it pgdb-bank dropdb --username=bank_admin bank_app

migrate_up:
	migrate -path db/migrations -database "postgresql://bank_admin:bank_password_qwertyuiop@localhost:5433/bank_app?sslmode=disable" -verbose up

migrate_up1:
	migrate -path db/migrations -database "postgresql://bank_admin:bank_password_qwertyuiop@localhost:5433/bank_app?sslmode=disable" -verbose up 1

migrate_down:
	migrate -path db/migrations -database "postgresql://bank_admin:bank_password_qwertyuiop@localhost:5433/bank_app?sslmode=disable" -verbose down

migrate_down1:
	migrate -path db/migrations -database "postgresql://bank_admin:bank_password_qwertyuiop@localhost:5433/bank_app?sslmode=disable" -verbose down 1

sqlc:
	sqlc -f db/sqlc.yaml generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/WooDMaNbtw/BankApp/db/sqlc Store

.PHONY: docker_up docker_down create_db drop_db migrate_up migrate_up1 migrate_down migrate_down1 sqlc test server mock
