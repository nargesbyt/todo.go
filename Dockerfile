FROM gcr.io/distroless/static-debian11

COPY ./todo /
COPY config.yaml /

EXPOSE 8080

CMD ["/todo"]