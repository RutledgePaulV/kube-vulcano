FROM golang:1.5.1-onbuild


ENTRYPOINT ["go-wrapper", "run"]