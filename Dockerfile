FROM golang:1.16.5 as development
WORKDIR /app
COPY . /app 
RUN ["go", "install"]
ENTRYPOINT ["pretgo"]
