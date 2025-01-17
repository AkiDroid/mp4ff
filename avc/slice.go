package avc

import (
	"bytes"
	"errors"

	"github.com/edgeware/mp4ff/bits"
)

// Errors for parsing and handling AVC slices
var (
	ErrNoSliceHeader      = errors.New("No slice header")
	ErrInvalidSliceType   = errors.New("Invalid slice type")
	ErrTooFewBytesToParse = errors.New("Too few bytes to parse symbol")
)

// SliceType - AVC slice type
type SliceType uint

func (s SliceType) String() string {
	switch s {
	case SLICE_I:
		return "I"
	case SLICE_P:
		return "P"
	case SLICE_B:
		return "B"
	case SLICE_SI:
		return "SI"
	case SLICE_SP:
		return "SP"
	default:
		return ""
	}
}

// AVC slice types
const (
	SLICE_P  = SliceType(0)
	SLICE_B  = SliceType(1)
	SLICE_I  = SliceType(2)
	SLICE_SP = SliceType(3)
	SLICE_SI = SliceType(4)
)

// GetSliceTypeFromNALU - parse slice header to get slice type in interval 0 to 4
func GetSliceTypeFromNALU(data []byte) (sliceType SliceType, err error) {

	if len(data) <= 1 {
		err = ErrTooFewBytesToParse
		return
	}

	naluType := GetNaluType(data[0])
	switch naluType {
	case 1, 2, 5, 19:
		// slice_layer_without_partitioning_rbsp
		// slice_data_partition_a_layer_rbsp

	default:
		err = ErrNoSliceHeader
		return
	}
	r := bits.NewEBSPReader(bytes.NewReader((data[1:])))

	// first_mb_in_slice
	if _, err = r.ReadExpGolomb(); err != nil {
		return
	}

	// slice_type
	var st uint
	if st, err = r.ReadExpGolomb(); err != nil {
		return
	}
	sliceType = SliceType(st)
	if sliceType > 9 {
		err = ErrInvalidSliceType
		return
	}

	if sliceType >= 5 {
		sliceType -= 5 // The same type is repeated twice to tell if all slices in picture are the same
	}
	return
}
