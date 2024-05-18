FROM golang:1.21.1-alpine3.12 as development
WORKDIR /app
COPY . /app
RUN ["go", "install"]
ENTRYPOINT ["pretgo"]
