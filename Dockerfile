FROM golang as test
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go mod download
CMD CGO_ENABLED=0 go test ./...

FROM golang as build
WORKDIR /app
COPY --from=test /app .
RUN CGO_ENABLED=0 go build -o project


FROM alpine
WORKDIR /app/server
COPY --from=build /app/project .
EXPOSE 80
ENTRYPOINT ["./project" ]