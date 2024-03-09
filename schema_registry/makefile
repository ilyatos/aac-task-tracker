generate:
	rm -rf ./pkg/** && protoc -I ./schemas --go_out paths=source_relative:./pkg schemas/**/*.proto schemas/**/**/*.proto
