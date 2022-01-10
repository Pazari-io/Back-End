# Welcome to Pazari Engine
This repository holds the current version of our backend engine, which helps us protect and distribute digital publications. 

## Purpose

An API server written in Go . Which handles

- Keep track of tasks in a database 

- Downloads and uploads of original copies 

- Extract extra information about different file types (e.g. BMP from audio files and measurement from media files)

- Watermark Audio, Video, Image, PDF files

- Encrypt and decrypt archive file using AES-256 algorithm 

## Version

  **1.0.0-alpha.4**

## Dev

- [Download and install - The Go Programming Language](https://go.dev/doc/install) , version 1.17+ on your target machine 

- Clone this repository 

- Set your environment variables on the .env file accordingly

- Build the project
  
  - ```bash
    go build -o api main.go
    ```

- Execute the api server
  
  - ```bash
    ./api
    ```

- To run tests
  
  - ```bash
    go test ./...
    ```

## Docker container 

Because the Pazari engine requires careful configuration, Dockerfile containerizes the whole dependencies on Debian bullseye.


```bash
docker build -t pazari-engine .
# development 
docker run -d -p 1337:1337 pazari-engine:latest 
# production (need to cache certificate on host and restart after crashes )
docker run -d -p 443:443 --restart=always -v certs:/certs pazari-engine:latest 
```

## API

In the current version, API requires relays on a secret key for authentication, which is set on the .env file.To connect with the HTTP client to this API server,
you need to add an Authorization header to your request.

```http
Authorization: Bearer YOUR_SECRET_KEY
```

`GET /api/health (unencrypted)`

Check if the server is available.

`POST /api/v1/auth/upload (authenticated)`

On success returns task, which can be used later to retrieve different versions (original, encrypted, watermarked), should only be shared with the file owner.

```json
{
    "taskID": "Li91cGxvYWRzL29yaWdpbmFsLzI1NDA0MjI0MDg3ODgxMjM0Y2IxZjE3MTY1ODZhYTE3NDlhNWFhOTQxMmNlNGNiNjQ4ODE4NmZlZDUzNDkxNWIucG5n"
}
```

 `GET /api/v1/auth/watermark?fileID=taskID (authenticated)`

To download the watermarked edition (works for audio, video, and images) for showing on the marketplace. 

`GET /api/v1/auth/download?fileID=taskID (authenticated)`

On success downloads the original or encrypted (for example, PDF). This should only be shared with the buyer after the smart contract confirms and verifies the transaction.



### Depencensies

Our engine uses well-known open-source software to handle encryption, authorization, database management, and media processing.You can read more about each of them.

[ImageMagic](https://github.com/ImageMagick/ImageMagick)

[FFmpeg](https://github.com/FFmpeg/FFmpeg)

[Fiber](https://github.com/gofiber/fiber)

[Gorm](https://github.com/go-gorm/gorm)

[Aubio](https://github.com/aubio/aubio)

[Pdfcpu](https://github.com/pdfcpu/pdfcpu)



## Security

We handle security and security issues with great care. Please contract `security [at] pazari.io` as soon as you find a valid vulnerability. 



## Important

Currently alpha and under development .




