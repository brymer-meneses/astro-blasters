
tidy:
	cd src && go mod tidy

# run the native client
run-native:
	cd src && go run cmd/native/main.go

# run the server
run-server:
	cd src && go run cmd/server/main.go
