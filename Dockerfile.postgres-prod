# syntax=docker/dockerfile:1

FROM postgres:12.9 as stockfy-postgres-prod

WORKDIR /stockfy-postgresql

# COPY sql_files/create_database.sql .
COPY sql_files/create_database.sql /docker-entrypoint-initdb.d/1-create_database.sql

ENTRYPOINT ["docker-entrypoint.sh"]
# RUN psql -U postgres -h localhost -f sql_files/create_database.sql

EXPOSE 5432

CMD ["postgres"]