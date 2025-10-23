ARG alpine_version=3.19
ARG go_version=1.23.2
ARG DATABASE_MAIN_HOST
ARG DATABASE_MAIN_PORT
ARG DATABASE_MAIN_USERNAME
ARG DATABASE_MAIN_PASSWORD


## Builder
FROM golang:$go_version-alpine$alpine_version AS builder

RUN apk update && \
    apk add make && \
    apk add ca-certificates gcc g++ libc-dev

WORKDIR /app

COPY . .

RUN go mod download

RUN make build


## Distribution
FROM alpine:latest
RUN apk add --no-cache bash

WORKDIR /app

RUN mkdir /app/config \
	&& mkdir /app/logs \
	&& mkdir /app/config/resources \
    && mkdir /app/storage \
    && mkdir /app/storage/admin \
    && mkdir /app/storage/admin/profile \
    && mkdir /app/storage/article \
    && mkdir /app/storage/estatement \
    && mkdir /app/static \
    && mkdir /app/static/images

COPY --from=builder /app/main /app/
COPY --from=builder /app/config /app/config
COPY --from=builder /app/db /app/db
COPY --from=builder /app/static /app/static

RUN chmod +x main \
    && chmod -R 770 /app/logs

RUN ls -l
RUN ls -la /app/config
RUN ls -la /app/config/resources
RUN pwd

# CMD /bin/bash -c ./start.sh
CMD ["./main", "start"]