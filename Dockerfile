# Command to build: docker build . -t docker-tutorial-image
# Command to run: docker run -p 3030:3001 -it --rm --name docker-tutorial-run docker-tutorial-image

FROM golang:alpine
ADD . /go/src/app
WORKDIR /go/src/app
CMD ["go", "run", "main.go"]
