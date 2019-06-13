FROM golang:1.11.2

RUN rm -rf /go/src/*

WORKDIR /go/src/app

RUN go get -u github.com/gorilla/mux

RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

RUN go get -u github.com/jinzhu/gorm

RUN go get -u golang.org/x/crypto/argon2

RUN go get -u github.com/dgrijalva/jwt-go

RUN go get -u github.com/go-ozzo/ozzo-validation

RUN go get -u github.com/go-ozzo/ozzo-validation/is