package base62

func Decode(s string) int64 {
	dict := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	base := int64(len(dict))
	d := int64(0)
	for _, ch := range s {
		for i, a := range dict {
			if a == ch {
				d = d*base + int64(i)
			}
		}
	}
	return d
}

func Encode(i int64) string {
	dict := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	base := int64(len(dict))
	digits := []int64{}
	for i > 0 {
		r := i % base
		digits = append([]int64{r}, digits...)
		i = i / base
	}

	r := []rune{}
	for _, d := range digits {
		r = append(r, dict[d])
	}
	return string(r)
}
