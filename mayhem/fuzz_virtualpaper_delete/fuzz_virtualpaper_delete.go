package fuzz_virtualpaper_delete
import (
    fuzz "github.com/AdaLogics/go-fuzz-headers"

    "tryffel.net/go/virtualpaper/process"
)

func mayhemit(bytes []byte) int {

    fuzzConsumer := fuzz.NewConsumer(bytes)

    fuzzId, err := fuzzConsumer.GetString()
    if err != nil {
        return 0
    }

    process.DeleteDocument(fuzzId)
    return 0
}

func Fuzz(data []byte) int {
    _ = mayhemit(data)
    return 0
}