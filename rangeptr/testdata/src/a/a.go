package testdata

func assignPtr() {
	src := []int{1, 2, 3}
	dst := make([]*int, len(src))
	for i, v := range src {
		dst[i] = &v // want `using pointer to the loop iteration variable "v"`
	}
}

func assignSrcPtr() {
	src := []int{1, 2, 3}
	dst := make([]*int, len(src))
	for i := range src {
		dst[i] = &src[i]
	}
}

func assignTmp() {
	src := []int{1, 2, 3}
	dst := make([]*int, len(src))
	for i, v := range src {
		tmp := v
		dst[i] = &tmp
	}
}

func appendPtr() {
	src := []int{1, 2, 3}
	dst := make([]*int, 0)
	for _, v := range src {
		dst = append(dst, &v) // want `using pointer to the loop iteration variable "v"`
	}
}

func appendSrcPtr() {
	src := []int{1, 2, 3}
	dst := make([]*int, 0)
	for i := range src {
		dst = append(dst, &src[i])
	}
}

func appendTmp() {
	src := []int{1, 2, 3}
	dst := make([]*int, 0)
	for _, v := range src {
		tmp := v
		dst = append(dst, &tmp)
	}
}
