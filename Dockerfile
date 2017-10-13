FROM golang
ADD .  /go/src/github.com/dcadenas/pagerank
RUN go install github.com/dcadenas/pagerank
WORKDIR /go/src/github.com/dcadenas/pagerank
