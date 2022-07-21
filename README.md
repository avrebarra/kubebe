# 🤺 kubebe

Basic HTTP server.

```sh
➜  ~ curl --location --request GET 'http://localhost:2323' \
--header 'Content-Type: application/json' \
{"code":"ok","data":{"id":"kubebe","state":"healthy","start_time":"2022-07-21T16:11:24.231377+07:00","uptime":"41s"}}%

➜  ~ curl --location --request POST 'http://localhost:2323/diceroll' \
--header 'Content-Type: application/json' \
--data-raw '{
    "dice_number": 5
}'
{"code":"ok","data":{"total":9,"sequence":[4,1,1,2,1]}}%
```
