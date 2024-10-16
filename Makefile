up:
	docker run -d --rm \
		--name postgres \
		-p 5432:5432 \
		-v ${PWD}/data/psql:/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD=local \
		-e POSTGRES_DB=metrics \
		postgres:16.4-alpine3.20

down:
	docker stop postgres

exec:
	docker exec -it postgres psql -U postgres