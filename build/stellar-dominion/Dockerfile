FROM node:22 AS builder
ARG API_BASE_URL="http://stellar-dominion-service:80/v1"
ARG USER_API_BASE_URL="http://user-service:80/v1/users"
WORKDIR /build
COPY frontend/stellar-dominion ./
RUN echo "PUBLIC_API_BASE_URL='${API_BASE_URL}'" > .env
RUN echo "PUBLIC_USER_API_BASE_URL='${USER_API_BASE_URL}'" >> .env
RUN npm install
RUN npm run build

FROM node:22-alpine3.18
ARG GIT_COMMIT_HASH="undefined"
ARG SERVER_ORIGIN="http://localhost:3001"
ARG NODE_PORT=3001
ENV GIT_COMMIT_HASH=$GIT_COMMIT_HASH
ENV PORT=$NODE_PORT
ENV ORIGIN=$SERVER_ORIGIN
WORKDIR /app
COPY --from=builder build/build bin/
COPY --from=builder build/package.json bin/
COPY --from=builder build/.env bin/
WORKDIR /app/bin
CMD ["node", "."]
