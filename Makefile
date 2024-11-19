LOCAL_BIN = $(CURDIR)/bin

generate-structs:
	mkdir -p src/pkg
	protoc --go_out=./src/pkg --go_opt=paths=source_relative \
			model/auth.proto
	protoc --go_out=./src/pkg --go_opt=paths=source_relative \
		model/profile.proto
	protoc --go_out=./src/pkg --go_opt=paths=source_relative \
			model/post.proto