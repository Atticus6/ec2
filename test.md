#mac build mac or window build window
#go build -o main

#mac build linux
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build main.go


# run
./main KEY_ID ACCESS_KEY REGIN PAYH TEXT