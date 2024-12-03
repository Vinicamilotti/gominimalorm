package error_helper

func MustExecute(err error) {
	if err != nil {
		panic(err)
	}
}

func MustReturn[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
