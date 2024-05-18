build:
	@go build -o bin/template-multimodal-db

run: build
	@./bin/template-multimodal-db


docker:
	echo "building docker file..."
	@docker build -t template-multimodal-db .
	echo "running API inside Docker container..."
	@docker run -p 5000:5000 template-multimodal-db