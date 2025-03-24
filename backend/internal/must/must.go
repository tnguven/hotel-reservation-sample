package must

func Panic[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}
