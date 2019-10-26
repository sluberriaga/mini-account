# Account

This project is a webapp that manages an account balance and transactions.
Transactions can be credits or debits anda are stored on a ledger for traceability.

## Instructions

In order to run the app simply type
```sh
$ ./build/account
```

In order to build the app just run 
```sh
$ make
```

## API
#### /api/account/balance
```sh
    curl --request GET \
        --url http://localhost:8080/api/account/balance \
        --header 'content-type: application/json'
```
#### /api/account/transaction
```sh
    curl --request POST \
      --url http://localhost:8080/api/account/transaction \
      --header 'content-type: application/json' \
      --data '{
    	"type": "credit",
    	"amount": 10
      }'
```
#### /api/account/transaction/:id
```sh
    curl --request GET \
      --url http://localhost:8080/api/account/transaction/7e4eeca5-2614-476b-8097-eddf09d819b \
      --header 'content-type: application/json' \
```
#### /api/account/transaction
```sh
    curl --request GET \
      --url 'http://localhost:8080/api/account/transaction?offset=0&limit=10' \
      --header 'content-type: application/json'
```


