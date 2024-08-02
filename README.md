# Crypto Currency Quant Trading Application


```
├── README.md
├── cmd
│   └── main.go
├── internal
│   ├── data        //get price data via database or api
│   ├── strategy    //run strategy and tools for it
│   └── util        //sqlite, utils etc...
├── pkg
│   └── upbit
│       └── model
└── quant.db
```

### .env
```
UPBIT_OPEN_API_ACCESS_KEY=access-key
UPBIT_OPEN_API_SECRET_KEY=secret-key
```


### Run
```bash
go run cmd/main.go
```



