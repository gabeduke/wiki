
preflight:
	mkdir -p $(CURDIR)/.data

dev:
	buffalo dev

db: preflight
	docker run -d --name postgres \
	-e POSTGRES_PASSWORD=postgres \
	-e PGDATA=/var/lib/postgresql/data/pgdata \
    -v $(CURDIR)/.data:/var/lib/postgresql/data \
	-p 5432:5432 \
	postgres

clean:
	docker stop postgres && docker rm postgres