package main

func Crop(a string, n int) (ret string) {

	r := []rune(a)
	if len(r) > n {
		r = r[:n]
	}

	ret = string(r)

	return
}

func CropBS(a []byte, n int) (ret string) {

	return Crop(string(a), n)
}
