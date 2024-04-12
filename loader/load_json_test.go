package loader

import "testing"

func TestLoadSeed(t *testing.T) {
	seed, err := LoadSeed("./seed.json")
	if err != nil {
		t.Error(err)
	}
	t.Log(seed)

}
