package main

func appendAtIndex[T any](slice []T, index int, toAppend T) []T {
	return append(append(slice[:index], toAppend), slice[index:]...)
}
