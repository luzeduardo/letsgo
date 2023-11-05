package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Fatalf("\ngot  %v;\nwant %v", actual, expected)
	}
}
