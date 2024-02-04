package challenge

import (
	"encoding/base64"
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Run("should generate challenges with correct length", func(t *testing.T) {
		for i := 0; i < 200; i++ {
			if len(Generate(i, 1)) != i+2 { // add 2 chars for prefix.
				t.Errorf("incorrect data for len %d", i)
			}
		}
	})

	t.Run("should generate correct prefix", func(t *testing.T) {
		if got, want := Generate(5, 123)[:4], "123:"; got != want {
			t.Errorf("unexpected prefix, got %q, want %q", got, want)
		}
	})
}

func TestChallenge(t *testing.T) {
	data := "4:foobar"
	sln, err := Solve(data)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err = Verify(data, sln); err != nil {
		t.Errorf("unexpected verification error: %v", err)
	}
	if err = Verify("barfoo", sln); err == nil {
		t.Errorf("expected verification error, got nil")
	}
	if err = Verify("9:foobar", sln); err == nil {
		t.Errorf("expected verification error, got nil")
	}
}

func FuzzChallenge(f *testing.F) {
	f.Fuzz(func(t *testing.T, data string) {
		data = "3:" + base64.RawURLEncoding.EncodeToString([]byte(data))
		sln, err := Solve(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err = Verify(data, sln); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
