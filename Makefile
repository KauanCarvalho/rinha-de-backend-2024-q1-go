database-check:
	until nc -z -v -w30 database 5432; do \
	  sleep 1; \
	done

database-drop:
	psql ${DATABASE_DSN} -c "\c postgres; DROP DATABASE backendfight;" || exit 0

database-create:
	psql $(DATABASE_DSN) -c "CREATE DATABASE backendfight;"

database-migration-up:
	migrate -path ./migrations -database ${DATABASE_DSN} -verbose up
