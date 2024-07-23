build:
	@go build -o bin/bootstrap main.go

run: build
	@./bin/bootstrap

prod:
	@CGO_ENABLED=1 GOOS=linux go build -o bin/bootstrap main.go; \
	zip -jr bootstrap.zip bin/

lambda:
	@GOOS=linux go build -o bin/bootstrap main.go; \
	zip -jr bootstrap.zip bin/