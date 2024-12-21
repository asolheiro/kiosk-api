## Stage: Build
FROM golang:1.23.4-alpine AS build

WORKDIR /app

COPY ./ .

RUN go mod download &&\
    go mod verify 
    
RUN go build -o ./bin/kiosk ./cmd/kiosk/kiosk.go

## Stage: Run
FROM scratch

WORKDIR /app

COPY --from=build /app/bin/kiosk/. .

EXPOSE 8080

ENTRYPOINT [ "./kiosk" ]