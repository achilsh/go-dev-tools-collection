#! /usr/bin/env bash 
curl -X POST http://127.0.0.1:5657/demo-server/v1/x2/generic \
  -H "Content-Type: application/json" \
  -d @- <<EOF | jq
  {
    "aa": 1,
    "age": 2
  }
EOF