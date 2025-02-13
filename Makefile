docker-api:
	echo "building docker api file"
	@docker build -f Dockerfile.api -t hotel-io/api .
	echo "running API inside docker container"
	@docker run --rm --env-file .env -p 5000:5000 hotel-io/api:latest
