include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the application
.PHONY: run
run:
	@go run ./services/contact/cmd -db-dsn=${DB_DSN}

## db/psql: connect to the database using psql
.PHONY: psql
psql:
	psql ${DB_DSN}

.PHONY: protobuf/generate
protobuf/generate:
	  protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative services/contact/protobuf/*.proto
.PHONY: protobuf/clean
protobuf/clean:
	rm services/contact/protobuf/*.go