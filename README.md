# SantaWebsite
Website of gifts for orphans for the new year. Everyone can become Santa for an orphan child

## Developers
1. Adilzhan Shirbayev
2. Nargiz Skakova
3. Baurzhan Saliyev

## Instructions to run (with docker)
1. ```cd SantaWebsite```
2. ```docker build -t santaweb .```
3. ```docker run -dp 8080:8080 santaweb```

## Instructions to run (without docker)
1. To run server
```console
go run app.go
```
2. if the programme runs successfully, you will get the following result
```console
$Connected to MongoDB!
$Starting server on port :8080...
```
3. Our server runs on port 8080. Just past localhost:8080
4. Done, you are on the website


* Sign up as a child if you would like to receive a gift and write down what you want

* Register as a volunteer if you want to fulfil the wishes of children

## Dependencies
Mongo driver: go get go.mongodb.org/mongo-driver
Gorilla mux: go get -u github.com/gorilla/mux
Bcrypt: go get golang.org/x/crypto/bcrypt
