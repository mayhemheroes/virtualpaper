FROM golang:1.19 as builder

COPY . /virtualpaper
WORKDIR /virtualpaper/mayhem

RUN go install github.com/dvyukov/go-fuzz/go-fuzz@latest github.com/dvyukov/go-fuzz/go-fuzz-build@latest
RUN go get github.com/dvyukov/go-fuzz/go-fuzz-dep
RUN go get github.com/AdaLogics/go-fuzz-headers
RUN apt update && apt install -y clang

RUN go-fuzz-build -libfuzzer -o fuzzvirtualpaper.a && \
    clang -fsanitize=fuzzer fuzzvirtualpaper.a -o fuzzvirtualpaper.libfuzzer

FROM debian:bookworm-slim
COPY --from=builder /virtualpaper/mayhem/fuzzvirtualpaper.libfuzzer /