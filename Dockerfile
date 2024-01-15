#FROM gcr.io/distroless/static-debian11
FROM debian:latest

COPY todo/todo /
#COPY ./dist / 
COPY config.yaml.dist /config.yaml

EXPOSE 8080
RUN chmod +x /todo
ENTRYPOINT ["/todo"]