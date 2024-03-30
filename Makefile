
user-service-build:
	docker build --tag user-service -f build/users/Dockerfile .

user-service-run:
	docker run -p 5432 -p 60001:60001 user-service