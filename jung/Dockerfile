FROM golang:1.20-buster
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get -yq update && apt-get -yq install gnupg2 ca-certificates cl-base64
COPY ./jung/ /jung/
WORKDIR /jung
RUN go mod download
RUN go mod tidy
RUN go install .

CMD ["/go/bin/de.janmeckelholt.jung"]
#CMD ["sleep", "5000"]