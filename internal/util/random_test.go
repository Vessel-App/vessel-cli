package util

import "testing"

func TestCanGetRandomString(t *testing.T) {
	length := 16
	got := RandomString(length)

	if len(got) != 16 {
		t.Errorf("random string was not correct length. wanted=%d, received=%d", length, len(got))
	}
}
