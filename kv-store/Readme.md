# Usage

- v1

```
curl -X POST -H "Content-Type: application/json" -d
'{"database": "database1", "key": "name", "value": {"John Doe","Other name","Third name"}}' http://localhost:3000/api/v1/set

```

- v2

```
curl -X POST -H "Content-Type: application/json" -d
'{
"database": "db6",
"key": "sp_secret_uat1",
"value": {
    "keyvault": "https://kv02uat1.azurekv.net",
    "expirydate": "expiry02243434uat1",
    "metadata": "test434343"
  }
}'
```

## Get all keys

- v1

```
curl -XGET http://localhost:3000/api/v1/all/db5
```

- v2

```
curl -XGET http://localhost:3000/api/v2/all/db6
```

## Get a value from "database"

- v1

```
curl -XGET http://localhost:3000/api/v1/get/db5/color3
```

- v2

```
curl -XGET http://localhost:3000/api/v2/get/db6/sp_secret_uat
```

## Delete a key from "database1"

- v1

curl -X DELETE "http://localhost:3000/api/v1/del/db5/color3"

- v2

http://localhost:3000/api/v2/del/db6/sp_secret_dev
