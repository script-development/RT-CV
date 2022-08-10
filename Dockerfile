# build backend
FROM golang:1.18-alpine3.16 AS backend

RUN mkdir /project
WORKDIR /project
COPY ./ ./

RUN go build -ldflags "-X main.AppVersion=$(git log --format='%H' -n 1)" -o rtcv

# build dashboard
FROM node:16-buster AS dashboard

WORKDIR /app
COPY ./dashboard/ .

ENV NEXT_TELEMETRY_DISABLED=1
RUN npm ci \
    && npm run build

# Setup the runtime
FROM alpine:3.16 AS runtime

RUN mkdir -p /project/dashboard \
    && apk add --no-cache tzdata

ENV TZ="Europe/Amsterdam"
WORKDIR /project
EXPOSE 4000

COPY --from=backend /project/rtcv /project/rtcv
COPY --from=dashboard /app/out /project/dashboard/out
COPY assets /project/assets


CMD [ "sh", "-c", "./rtcv" ]
