FROM alpine:latest
LABEL authors="Connor Molz"

ENTRYPOINT ["top", "-b"]