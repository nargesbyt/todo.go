FROM gcr.io/distroless/static-debian11

#COPY ./todo /
COPY ./dist / 
COPY config.yaml.dist /config.yaml

EXPOSE 8080

CMD ["/todo"]