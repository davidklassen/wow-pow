package challenge

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	for i := 0; i < 200; i++ {
		if len(Generate(i)) != i {
			t.Errorf("incorrect data for len %d", i)
		}
	}
}

func TestChallenge(t *testing.T) {
	data := "foobar"
	bits := 5
	sln := Solve(data, bits)
	if !Verify(data, sln, bits) {
		t.Errorf("expected valid solution")
	}
	if Verify("barfoo", sln, bits) {
		t.Errorf("expected invalid solution")
	}
	if Verify(data, sln, bits+5) {
		t.Errorf("expected invalid solution")
	}
}

func FuzzChallenge(f *testing.F) {
	bits := 3
	f.Fuzz(func(t *testing.T, data string) {
		sln := Solve(data, bits)
		if !Verify(data, sln, bits) {
			t.Errorf("expected valid solution")
		}
	})
}
