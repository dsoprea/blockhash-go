// Go implementation of blockhash (http://blockhash.io), a perceptual hash of images
//
// Reference implementation: https://github.com/commonsmachinery/blockhash-python/tree/0d76144cf5b6ac149ff7612ab433a48b0fb9139b
package blockhash

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"sort"
)

type Hash struct {
	// bit representation of the hash
	Bits []int
}

// Generates the block hash for a given image
// `bits` should be a power of 2 number (e.g. 2, 4, 8, 16, etc.); the number of output bits is equal to bits^2
func Blockhash(reader io.Reader, bits int) (*Hash, error) {
	img, _, err := image.Decode(reader)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	size := bounds.Size()
	blockSizeX, blockSizeY := size.X/bits, size.Y/bits

	result := make([]int, bits*bits)

	// divide the image into `bits * bits` grids
	// of size blockSizeX * blockSizeY
	// then add up the r+g+b values for each pixel
	for y := 0; y < bits; y++ {
		yBlockSizeY := y * blockSizeY
		yBits := y * bits
		for x := 0; x < bits; x++ {
			var value int
			xBlockSizeX := x * blockSizeX

			for iy := 0; iy < blockSizeY; iy++ {
				for ix := 0; ix < blockSizeX; ix++ {
					cx := xBlockSizeX + ix
					cy := yBlockSizeY + iy
					r, g, b, _ := img.At(cx, cy).RGBA()
					value += int(r + g + b)
				}
			}

			result[yBits+x] = value
		}
	}

	translateBlocksToBits(result, blockSizeX*blockSizeY)

	hash := &Hash{
		Bits: result,
	}
	return hash, nil
}

// A string representation of the hash in hex format
func (hash *Hash) ToHex() string {
	size := 16
	chunks := len(hash.Bits) / size
	var buffer bytes.Buffer
	for i := 0; i < chunks; i++ {
		bitChunk := reverse(hash.Bits[i*size : (i+1)*size])
		var num uint16
		for i, bit := range bitChunk {
			if bit == 1 {
				num |= 1 << uint(i)
			}
		}
		buffer.WriteString(fmt.Sprintf("%04x", num))
	}

	return buffer.String()
}

func translateBlocksToBits(blocks []int, pixelsPerBlock int) []int {
	const numHorizontalBands = 4
	halfBlockValue := pixelsPerBlock * 256 * 3 / 2
	bandSize := len(blocks) / numHorizontalBands

	for i := 0; i < numHorizontalBands; i++ {
		begin, end := i*bandSize, (i+1)*bandSize
		m := median(blocks[begin:end])
		mGreaterThanHalfBlockValue := m > halfBlockValue
		for j := begin; j < end; j++ {
			v := blocks[j]
			if v > m || (mGreaterThanHalfBlockValue && abs(v-m) < 1) {
				blocks[j] = 1
			} else {
				blocks[j] = 0
			}
		}
	}

	return blocks
}

func reverse(ints []int) []int {
	for i, j := 0, len(ints)-1; i < j; i, j = i+1, j-1 {
		ints[i], ints[j] = ints[j], ints[i]
	}
	return ints
}

func median(data []int) int {
	len := len(data)
	cpy := make([]int, len)
	copy(cpy, data)
	sort.Ints(cpy)

	mid := len / 2
	if len%2 == 0 {
		return (cpy[mid] + cpy[mid+1]) / 2
	}

	return cpy[mid]
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}
