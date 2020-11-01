# hashbang
Instagram tag manager
## Run the tests
```
$ go test github.com/arctair/hashbang/v1
$ go test -tags acceptance
```
or
```
$ nodemon
```
### Run the tests against a deployment
```
$ BASE_URL=https://hashbang.arctair.com go test -tags acceptance
```
## Run the server
```
$ go run .
$ curl localhost:5000
```
## Build, deploy, and verify
```
$ scripts/deploy
```
