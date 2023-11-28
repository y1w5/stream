#!/bin/bash

QUERY=""
if [[ "$LIMIT" != "" ]]; then
	QUERY="?limit=$LIMIT"
fi

curl -s "localhost:8080/v1/pages.list.std$QUERY" > list.json
curl -s "localhost:8080/v1/pages.list.std$QUERY" > list.json
curl -s "localhost:8080/v1/pages.list.std$QUERY" > list.json

curl -s "localhost:8080/v1/pages.list.exp$QUERY" > list.json
curl -s "localhost:8080/v1/pages.list.exp$QUERY" > list.json
curl -s "localhost:8080/v1/pages.list.exp$QUERY" > list.json

curl -s "localhost:8080/v2/pages.list$QUERY" > list.json
curl -s "localhost:8080/v2/pages.list$QUERY" > list.json
curl -s "localhost:8080/v2/pages.list$QUERY" > list.json

curl -s "localhost:8080/v2/pages.stream$QUERY" > stream.json
curl -s "localhost:8080/v2/pages.stream$QUERY" > stream.json
curl -s "localhost:8080/v2/pages.stream$QUERY" > stream.json

curl -s "localhost:8080/v3/pages.list$QUERY" > list.json
curl -s "localhost:8080/v3/pages.list$QUERY" > list.json
curl -s "localhost:8080/v3/pages.list$QUERY" > list.json

curl -s "localhost:8080/v3/pages.stream$QUERY" > stream.json
curl -s "localhost:8080/v3/pages.stream$QUERY" > stream.json
curl -s "localhost:8080/v3/pages.stream$QUERY" > stream.json
