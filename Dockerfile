FROM golang:bullseye AS build

WORKDIR /app
COPY ./ ./
RUN go mod download
RUN go build -o /app/api main.go

FROM debian:bullseye-slim
LABEL maintainer ="Pazari Team <hello@pazari.io>"
ARG DEBIAN_FRONTEND=noninteractive
WORKDIR /app

COPY --from=build /app/api /app/api

#ENV SECRET_KEY=NOT_SAME_AS_PRODUCTION_KEY
#ENV PORT=1337

RUN apt-get update && apt-get install -y \
    curl\
    git

# install Magick, Aubio, ffmpeg, ffprobe, pdfcpu, 7z
RUN apt-get install -y \
    ffmpeg\
    imagemagick\
    aubio-tools\
    libaubio-dev\
    libaubio-doc\
    xz-utils\
    p7zip-full

# install PDFCPU 
RUN curl "https://pazari-storage.sgp1.digitaloceanspaces.com/pdfcpu_0.3.13_Linux_x86_64.tar.xz" -o pdfcpu.tar.xz\
    && tar -xvf pdfcpu.tar.xz\
    && rm pdfcpu.tar.xz\
    && mv pdfcpu_0.3.13_Linux_x86_64 pdfcpu\
    && cd pdfcpu\
    && cp pdfcpu /usr/bin/pdfcpu\
    && cd ..\
    && rm -rf pdfcpu

COPY ./data ./data
COPY ./uploads ./uploads
#COPY ./certs   ./cert-cache

#DEV
#COPY ./.env ./.env
COPY ./.env.prod ./.env
#DEV   
#EXPOSE 1337    

#Production
VOLUME ["/app/cert-cache"]

EXPOSE 443
CMD ["./api"]