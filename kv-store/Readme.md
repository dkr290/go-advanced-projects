# Usage

curl -X POST -H "Content-Type: application/json" -d '{"database": "database1", "key": "name", "value": "John Doe"}' http://localhost:3000/api/v1/set

## Get a value from "database1"

curl "http://localhost:3000/api/v1/get?database=database1&key=name"

## Delete a key from "database1"

curl -X DELETE "http://localhost:3000/api/v1/delete?database=database1&key=name"
