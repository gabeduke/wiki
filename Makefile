
EXECUTABLES = buffalo docker docker-compose go
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

migrate: db
	$(info performing migrations)
	buffalo pop migrate up

dev: migrate
	buffalo dev

db:
	docker-compose up -d

db-seed: db
	$(info creating db tables)
	buffalo pop create -a

clean:
	docker-compose down