docker run --name mysql -e MYSQL_ROOT_PASSWORD=Password123 -p 3306:3306 -d mysql:latest
docker run --name cloudbeaver  -p 8978:8978 -d dbeaver/cloudbeaver:latest
docker network create database
docker network connect database cloudbeaver
docker network connect database mysql
