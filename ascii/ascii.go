package ascii

func IsLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
