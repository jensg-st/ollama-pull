FROM docker.io/library/golang:1.22 as builder

COPY . . 

RUN CGO_ENABLED=false go build -tags osusergo,netgo -o /ollama-pull .

FROM alpine:3.19

COPY --from=builder /ollama-pull /ollama-pull

ENTRYPOINT ["/ollama-pull"]