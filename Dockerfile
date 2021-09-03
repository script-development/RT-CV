# build backend
FROM golang:1.17-buster AS backend

RUN mkdir /project
WORKDIR /project
COPY ./ ./

RUN go build -o rtcv

# build dashboard
FROM node:16-alpine AS dashboard

WORKDIR /app
COPY ./dashboard/ .

ENV NEXT_TELEMETRY_DISABLED=1
RUN npm ci && npm run build

# Setup the runtime
FROM ubuntu AS runtime

RUN ln -fs /usr/share/zoneinfo/Europe/Amsterdam /etc/localtime \
    && mkdir -p /project/dashboard \
    && apt update \
    && apt install -y --no-install-recommends poppler-utils \
    fonts-dejavu fonts-freefont-ttf fonts-ubuntu ttf-bitstream-vera \
    && apt autoremove \
    && apt clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=backend /project/rtcv /project/rtcv
COPY --from=dashboard /app/out /project/dashboard/out
COPY assets /project/assets

WORKDIR /project

EXPOSE 4000

CMD [ "sh", "-c", "./rtcv" ]
