FROM golang:1.19 as builder

COPY . /virtualpaper
WORKDIR /virtualpaper/mayhem

RUN go install github.com/dvyukov/go-fuzz/go-fuzz@latest github.com/dvyukov/go-fuzz/go-fuzz-build@latest
RUN go get github.com/dvyukov/go-fuzz/go-fuzz-dep
RUN go get github.com/AdaLogics/go-fuzz-headers
RUN apt update && apt install -y clang

RUN cd fuzz_virtualpaper_process && go-fuzz-build -libfuzzer -o fuzz_virtualpaper_process.a && \
    clang -fsanitize=fuzzer fuzz_virtualpaper_process.a -o fuzz_virtualpaper_process.libfuzzer

RUN cd fuzz_virtualpaper_delete && go-fuzz-build -libfuzzer -o fuzz_virtualpaper_delete.a && \
    clang -fsanitize=fuzzer fuzz_virtualpaper_delete.a -o fuzz_virtualpaper_delete.libfuzzer

FROM debian:bookworm-slim
COPY --from=builder /virtualpaper/mayhem/fuzz_virtualpaper_process/fuzz_virtualpaper_process.libfuzzer /
COPY --from=builder /virtualpaper/mayhem/fuzz_virtualpaper_delete/fuzz_virtualpaper_delete.libfuzzer /