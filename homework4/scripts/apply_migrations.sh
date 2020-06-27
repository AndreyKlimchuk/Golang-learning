 export POSTGRESQL_URL='postgres://gorello:12345@localhost:5432/gorello?sslmode=disable'
 migrate -database ${POSTGRESQL_URL} -path ../db/migrations $1