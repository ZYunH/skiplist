package skiplist

func strcmp(a, b string) int {
	mLen := len(a)
	if len(a) > len(b) {
		mLen = len(b)
	}
	for i := 0; i < mLen; i++ {
		if a[i] != b[i] {
			return int(a[i] - b[i])
		}
	}
	return 0
}
