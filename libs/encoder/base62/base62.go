package base62

func Decode(s string) uint64 {
	dict := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	base := uint64(len(dict))
	d := uint64(0)
	for _, ch := range s {
		for i, a := range dict {
			if a == ch {
				d = d*base + uint64(i)
			}
		}
	}
	return d
}

func Encode(i uint64) string {
	dict := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	base := uint64(len(dict))
	digits := []uint64{}
	for i > 0 {
		r := i % base
		digits = append([]uint64{r}, digits...)
		i = i / base
	}

	r := []rune{}
	for _, d := range digits {
		r = append(r, dict[d])
	}
	return string(r)
}
