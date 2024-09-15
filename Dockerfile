FROM ubuntu:latest
LABEL authors="woodman"

ENTRYPOINT ["top", "-b"]