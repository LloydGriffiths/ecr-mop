FROM golang:1.9

ADD cmd src/github.com/LloydGriffiths/ecr-mop/cmd/
ADD mop src/github.com/LloydGriffiths/ecr-mop/mop/
ADD vendor src/github.com/LloydGriffiths/ecr-mop/vendor/

RUN go install github.com/LloydGriffiths/ecr-mop/cmd/ecr-mop/
