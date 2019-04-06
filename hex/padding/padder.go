package padding

func Pad(input string, outputLength int) string {
	diff := outputLength - len(input)
	for i := 0; i < diff; i++ {
		input = "0" + input
	}

	return input
}
