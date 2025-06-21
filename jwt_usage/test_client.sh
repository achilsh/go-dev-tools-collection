#! /usr/bin/env  bash 
## 生成token
#curl http://127.0.0.1:8080/api/

## 使用 token 来访问接口
 curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6IkNocmlzdG9waGVyIiwiYWFhIjo5OTkwMSwiZXhwIjoxNzUwNTE4NDY3LCJpc3MiOiIxMTExMSIsInN1YiI6IjIyMjIyMjIifQ.HA3Nl2OuquFwWr3-U7f7T1zxzhKCR67Xxfls9FqZ8Mc" http://127.0.0.1:8080/api/private/
