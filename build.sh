set -xe

go get github.com/gin-gonic/gin
go get github.com/dgrijalva/jwt-go v3.2.0+incompatible
go get github.com/gofrs/uuid v4.0.0+incompatible
go get github.com/jackc/pgx/v4 v4.13.0
go get golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
go get golang.org/x/net v0.0.0-20211105192438-b53810dc28af

go buil -o bin/application main.go
