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
