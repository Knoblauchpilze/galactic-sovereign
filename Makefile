
# Provide tag as:
# GIT_COMMIT_HASH=$(git rev-parse --short HEAD)

user-service-build:
	docker build --build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} --tag user-service:${GIT_COMMIT_HASH} -f build/users/Dockerfile .

user-service-run: user-service-build
	docker run -p 5432 -p 60001:60001 user-service:${GIT_COMMIT_HASH}
