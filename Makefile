up-postgres-docker:
	docker compose up -d postgres

down-postgres-docker:
	docker compose down postgres


install-all: install-gowrap install-iface install-sqlc install-linter

install-gowrap:
	go install github.com/hexdigest/gowrap/cmd/gowrap@v1.4.1

install-iface:
	go install github.com/vburenin/ifacemaker@v1.2.1

install-sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0

install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
