docker run --name mysql -e MYSQL_ROOT_PASSWORD=Password123 -p 3306:3306 -d mysql:latest
docker run --name cloudbeaver -p 8978:8978 -d dbeaver/cloudbeaver:latest
docker network create database
docker network connect database cloudbeaver
docker network connect database mysql

## sample env vars to be added in docker for the backend

DB_USER="root"
DB_PASSWORD="Password123"
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="todo"
HTTP_LISTEN_ADDR="localhost:3000"

docker build --progress=plain -t frontend .

docker run --name frontend -p 8090:8090 -e BACKEND_SERVICE="10.42.0.142:3000" frontend
docker run --name backend -d -p 3000:3000 -e DB_USER="root" -e DB_PASSWORD="Password123" -e DB_HOST="10.42.0.142" -e DB_PORT="3306" -e DB_NAME="todo" backend
