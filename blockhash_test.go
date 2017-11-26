package blockhash

import (
	"fmt"
	"log"
	"os"
	"testing"
)

const defaultFile = "fixtures/tulips-1083572_1920.jpg"

var fileHashMap = map[string]string{
	"fixtures/tulips-1083572_1920.jpg":         "801f3818fd80ffecfffc4ffc0e38000899ff193e1d38083c5c6c9e781e310ff0",
	"fixtures/almond-blossom-1229138_1920.jpg": "278225c6176e47fc005d02fa351dffd87ffa7bf003908624cc764face5887e80",
}

func TestBlockhash(t *testing.T) {
	for file, fileHash := range fileHashMap {
		reader, err := os.Open(file)
		if err != nil {
			t.Error(err)
		}

		hash, _ := Blockhash(reader, 16)
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
	hash, _ := Blockhash(reader, 16)    // Divide the image into 16x16 blocks, and calculate the hash
	fmt.Println("Bits", hash.Bits[:32]) // You should use the entire content in `Bits` though
	fmt.Println("Hex", hash.ToHex())
	// Output:
	// Bits [1 0 0 0 0 0 0 0 0 0 0 1 1 1 1 1 0 0 1 1 1 0 0 0 0 0 0 1 1 0 0 0]
	// Hex 801f3818fd80ffecfffc4ffc0e38000899ff193e1d38083c5c6c9e781e310ff0
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

func TestDustinBlockhash1(t *testing.T) {
	filepaths := []string {
		"/home/dustin/Downloads/20170618_155330.jpg",
		"/home/dustin/Downloads/20170618_155330-small.jpg",
		"/home/dustin/Downloads/amazing-mountain-valley-wallpaper-29910-30628-hd-wallpapers.jpg",
		"/home/dustin/Downloads/amazing-mountain-valley-wallpaper-29910-30628-hd-wallpapers-small.jpg",
	}

	for _, filepath := range filepaths {
		fmt.Println(filepath)
		reader, err := os.Open(filepath)
		if err != nil {
			t.Error(err)
		}

		hash, _ := Blockhash(reader, 16)
		fmt.Println(hash.ToHex())

		fmt.Println("")
	}
}
