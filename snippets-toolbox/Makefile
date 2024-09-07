run-local:
	go run ./cmd/web/*


install-mysql:
	docker run --name mysql -p 3306:3306 -v \
		     ./dbdata:/var/lib/mysql -e \
		     MYSQL_ROOT_PASSWORD=password  -d docker.io/mysql:8.4.2

run-mysql:
	docker start mysql
