up-postgres-docker:
	docker compose up -d postgres

down-postgres-docker:
	docker compose down postgres

install-gowrap:
	go install github.com/hexdigest/gowrap/cmd/gowrap@latest

install-iface:
	go install github.com/vburenin/ifacemaker@latest

intall-sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest