## Run docker-compose

### run docker containers
```dockerfile
docker compose --env-file .env.docker up -d
```

### stop docker containers
```dockerfile
docker compose --env-file .env.docker down
```

### connect to the postgres via docker | password - bank_password_qwertyuiop
```cmd
docker exec -it pgdb-bank psql -U bank_admin -d bank_app
```

### create a new db migration:
```cmd
migrate create -ext sql -dir db/migration -seq <migration_name>
```
