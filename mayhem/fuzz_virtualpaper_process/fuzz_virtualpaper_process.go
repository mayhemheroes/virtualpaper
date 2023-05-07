package fuzz_virtualpaper_process

import (
    fuzz "github.com/AdaLogics/go-fuzz-headers"

    "tryffel.net/go/virtualpaper/process"
)

func mayhemit(bytes []byte) int {

    if len(bytes) > 2 {
        num := int(bytes[0])
        bytes = bytes[1:]
        fuzzConsumer := fuzz.NewConsumer(bytes)
        
        switch num {
            case 1:
                fuzzName, err := fuzzConsumer.GetString()
                if err != nil {
                    return 0
                }

                process.GetHash(fuzzName)
                return 0

            case 0:
                fuzzName, err := fuzzConsumer.GetString()
                if err != nil {
                    return 0
                }

                process.MimeTypeFromName(fuzzName)
                return 0
        
            default:
                fuzzType, err := fuzzConsumer.GetString()
                fuzzName, err := fuzzConsumer.GetString()
                if err != nil {
                    return 0
                }

                process.MimeTypeIsSupported(fuzzType, fuzzName)
                return 0
        }
    }
    return 0
}

func Fuzz(data []byte) int {
    _ = mayhemit(data)
    return 0
}