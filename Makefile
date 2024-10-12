
tidy:
	cd src && go mod tidy

run-client-web:
	cd src && go run cmd/client/web/main.go

run-client-native:
	cd src && go run cmd/client/native/main.go 

run-server:
	cd src && go run cmd/server.go
