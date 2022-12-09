
EXECUTABLES = buffalo docker docker-compose go
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

migrate:
	$(info performing migrations)
	buffalo pop migrate up

dev: migrate
	buffalo dev

up:
	docker-compose up -d

db-seed:
	$(info seeding db)
	buffalo task db:seed

db-reset:
	$(info resetting db)
	buffalo pop reset

reset: db-reset db-seed

down:
	docker-compose down