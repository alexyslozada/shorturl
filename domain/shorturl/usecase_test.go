package shorturl

import (
	"testing"
)

func Test_randomPATH(t *testing.T) {
	resp := randomPATH()
	if len(resp) != MaxLetters {
		t.Fatalf("se esperaban %d letras, se obtuvieron %d", MaxLetters, len(resp))
	}

	t.Logf("randomPath() %s", resp)
}
