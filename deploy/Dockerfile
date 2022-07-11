FROM golang:1.15-buster AS build

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod tidy && go mod vendor

COPY ./ ./

RUN go build -o ./marketplace-application
RUN ls

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/.env /app/.env
COPY --from=build /app/marketplace-application /app/marketplace-application

EXPOSE 4040

ENTRYPOINT ["/app/marketplace-application"]