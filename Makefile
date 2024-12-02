test-unit:
	go test -cover -v ./internal/applications/core/... -coverprofile=coverage.out

up_db:
	mkdir -p data
	docker compose --project-directory . --file common/config/docker-compose.yml up db -d

clean:
	docker compose --project-directory . --file common/config/docker-compose.yml down db
	docker container prune
	sudo rm -r data

migrate_gen:
	migrate create -ext sql -dir common/schema -seq $(filter-out $@,$(MAKECMDGOALS))

migrate_up:
	migrate -database "postgresql://postgres:SuiseiKawaii@localhost:5432" -path "common/schema" up

migrate_drop:
	migrate -database "postgresql://postgres:SuiseiKawaii@localhost:5432" -path "common/schema" drop -f

test:
	go test -cover -coverprofile=coverage.out $(shell go list ./internal/... | grep -vE "logger|applications|utils") -parallel 7


test-core:
	go test -cover ./internal/applications/core/... -coverprofile=coverage.out

test-secondary:
	go test -cover ./internal/adapters/secondary/...

test-primary:
	go test -cover ./internal/test/httpprotocol/...

test-verbose:
	go test -cover -v  ./internal/... -coverprofile=coverage.out

watch:
	air -c ./common/config/.air.toml
