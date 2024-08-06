run-local:
	go run ./cmd/web/*


install-mysql:
	docker run --name mysql -v \
		     ./dbdata:/var/lib/mysql -e \
		     MYSQL_ROOT_PASSWORD=password -d mysql:8.4.2

run-mysql:
	docker start mysql
