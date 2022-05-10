postgresql:
	docker run \
    -e POSTGRES_USER=user \
    -e POSTGRES_PASSWORD=pass \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -p 5432:5432 \
    postgres