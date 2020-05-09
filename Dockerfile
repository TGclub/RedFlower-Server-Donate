FROM golang
WORKDIR /app
ADD . .
COPY ./config.json /app/config.json
EXPOSE 8000

CMD ["/app/main"]
