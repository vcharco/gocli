package gocliutils

func RepeatString(str string, times int) string {
	res := ""
	for i := 0; i < times; i++ {
		res += str
	}

	return res
}
