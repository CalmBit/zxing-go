package internal

import "errors"

func CompareRuneSlices(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// SliceCopyU32 mimics the behaviour of System.arrayCopy - rather than
// settling for the lesser of values, we want assurance we're using two arrays
// of the proper size, and that we can specify arbitrary points to copy within them
// safely.
func SliceCopyU32(src, dest []uint32, srcPos, destPos, length int) error {
	if src == nil || dest == nil {
		return errors.New("src and dest must both be []uint32")
	}

	if srcPos < 0 || destPos < 0 {
		return errors.New("srcPos and destPos must both be positive")
	}

	if length < 0 {
		return errors.New("length must be positive")
	}

	if len(src) < (srcPos + length) {
		return errors.New("srcPos + length is larger than the length of src")
	}

	if len(dest) < (destPos + length) {
		return errors.New("destPos + length is larger than the length of dest")
	}

	for i := 0; i < length; i++ {
		dest[destPos+i] = src[srcPos+i]
	}

	return nil
}

func SliceCloneU32(slice []uint32) []uint32 {
	newSlice := make([]uint32, len(slice))
	copy(newSlice, slice)
	return newSlice
}
