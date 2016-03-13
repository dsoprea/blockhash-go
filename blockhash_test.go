package blockhash

import (
	"fmt"
	"log"
	"os"
	"testing"
)

const defaultFile = "fixtures/tulips-1083572_1920.jpg"

var fileHashMap = map[string]string{
	"fixtures/tulips-1083572_1920.jpg": "801e3818fd80ffecfffc0ffc0e38000899ff193c1d38083c5c6c9e781e310fb0",
	"fixtures/almond-blossom-1229138_1920.jpg": "27822586176e47fc005d02fa351df7d87bfa7bf003908624cc764facc5887e80",
}

func TestBlockhash(t *testing.T) {
	for file, fileHash := range fileHashMap {
		reader, err := os.Open(file)
		if err != nil {
			t.Error(err)
		}

		hash := Blockhash(reader, 16)

		if hash.ToHex() != fileHash {
			t.Error("Hash doesn't match for", file)
		}
		reader.Close()
	}
}

func ExampleBlockhash() {
	reader, err := os.Open(defaultFile) // Get the file reader for the image
	if err != nil {
		log.Fatal(err)
	}
	hash := Blockhash(reader, 16) // Divide the image into 16x16 blocks, and calculate the hash
	fmt.Println("Bits", hash.Bits[:32]) // You should use the entire content in `Bits` though
	fmt.Println("Hex", hash.ToHex())
	// Output:
	// Bits [1 0 0 0 0 0 0 0 0 0 0 1 1 1 1 0 0 0 1 1 1 0 0 0 0 0 0 1 1 0 0 0]
	// Hex 801e3818fd80ffecfffc0ffc0e38000899ff193c1d38083c5c6c9e781e310fb0
}

func BenchmarkBlockhash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		reader, err := os.Open(defaultFile)
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()

		Blockhash(reader, 16)

		b.StopTimer()
		reader.Close()
		b.StartTimer()
	}
}
