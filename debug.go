package skiplist

func printSkipList(l *SkipList) {
	for i := l.level - 1; i >= 0; i-- {
		print(i, " ")
		hdr := l.header.levels[i]
		print("[hdr span:", hdr.span, "] -> ")
		x := hdr.next

		for x != nil {
			print("[val:", x.val, " span:", x.levels[i].span, "] -> ")
			x = x.levels[i].next
		}
		print("nil")
		print("\r\n")
	}
	print("\r\n")
}
