# CMPE-273-Assignment-1# 


This is a virtual stock trading system for whoever wants to learn how to invest in stocks.


## Usage
Clone the repository CMPE-273-Assignment-1

###Start the  server:

```
cd CMPE-273-Assignment-1
go run server.go
```
###Start the Client :

######Make RPC call from client for purchasing stocks
Eg: 
go run client.go GOOG:50%,YHOO:50% 3000
for buying 50% google and 50% yahoo stocks of the available 

######Make RPC call from client for Checking your portfolio (loss/gain)
go run client.go "Trade Id"


