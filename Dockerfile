# build backend
FROM golang:1.18-buster AS backend

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

# Build cv pdf generator
FROM dart as pdf_generator

COPY ./pdf_generator/ .
RUN dart pub get && dart compile exe bin/pdf_generator.dart

# Setup the runtime
FROM ubuntu AS runtime

RUN ln -fs /usr/share/zoneinfo/Europe/Amsterdam /etc/localtime \
    && mkdir -p /project/dashboard /project/pdf_generator/bin \
    && rm -rf /var/lib/apt/lists/*

COPY --from=backend /project/rtcv /project/rtcv
COPY --from=dashboard /app/out /project/dashboard/out
COPY --from=pdf_generator /root/bin/pdf_generator.exe /project/pdf_generator/bin/pdf_generator.exe
COPY pdf_generator/fonts /project/pdf_generator/fonts
COPY assets /project/assets

WORKDIR /project

EXPOSE 4000

CMD [ "sh", "-c", "./rtcv" ]
