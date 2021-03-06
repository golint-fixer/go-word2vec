package word2vec

import (
	"math"
	"os"
	"testing"
)

func TestModelLoadError(t *testing.T) {
	tests := []struct {
		path string
		kind int
	}{
		{
			"test/broken.model.bin",
			1,
		},
		{
			"/",
			1,
		},
		{
			"invalid_path",
			0,
		},
	}

	for _, test := range tests {
		inf, err := os.Open(test.path)
		defer inf.Close()
		if test.kind == 0 && err == nil {
			t.Errorf("Invalid path should make errors")
		}
		_, err = NewModel(inf)
		if test.kind == 1 && err == nil {
			t.Errorf("Broken model should make errors")
		}
	}

}

func TestModelLoadOpen(t *testing.T) {
	ifname := "test/model.bin"
	inf, err := os.Open(ifname)
	defer inf.Close()
	if err != nil {
		t.Errorf("%q\n", err)
	}

	model, err := NewModel(inf)
	if model == nil {
		t.Errorf("Model was not loaded\t%q\n", err)
		return
	}
	t.Logf("Model\tvocab=%d\tvector=%d\n", model.GetVocabSize(), model.GetVectorSize())

	if model.GetVocabSize() != len(model.GetNorms()) {
		t.Errorf("Vocab norms error")
	}
	if model.GetVocabSize() != len(model.GetVocab()) {
		t.Errorf("Vocab size error")
	}
	if len(model.data) != len(model.GetConnectedVector()) {
		t.Errorf("GetConnectedVector error")
	}
}

func TestModelLoadSim(t *testing.T) {
	ifname := "test/model.bin"
	inf, err := os.Open(ifname)
	defer inf.Close()
	if err != nil {
		t.Errorf("%q\n", err)
	}

	model, err := NewModel(inf)
	if model == nil {
		t.Errorf("Model was not loaded\t%q\n", err)
		return
	}
	t.Logf("Model\tvocab=%d\tvector=%d\n", model.GetVocabSize(), model.GetVectorSize())

	wv0 := model.GetNormalizedVector("NOT_FOUND_WORD")
	if wv0 != nil {
		t.Errorf("%v is not nil", wv0)
	}

	wv1 := model.GetNormalizedVector("he")
	wv2 := model.GetNormalizedVector("she")
	wv3 := model.GetNormalizedVector("with")

	val1 := wv1.Dot(wv2)
	val2 := wv1.Dot(wv3)
	if val1 < val2 {
		t.Errorf("wv1・w2 < w1・w3\n")
	}

	sim1, err := model.Similarity("he", "she")
	if err != nil || sim1 != val1 {
		t.Errorf("Error: Similarity()")
	}
	sim2, err := model.Similarity("he", "NOT_FOUND_WORD")
	if sim2 != 0 || err == nil {
		t.Errorf("Error: Similarity() for not found word")
	}
	sim3, err := model.Similarity("NOT_FOUND_WORD", "he")
	if sim3 != 0 || err == nil {
		t.Errorf("Error: Similarity() for not found word")
	}
}

func TestModelLoadVec(t *testing.T) {
	ifname := "test/model.bin"
	inf, err := os.Open(ifname)
	defer inf.Close()
	if err != nil {
		t.Errorf("%q\n", err)
	}
	w2vtest(t, inf)
}

func TestModelWithoutNewline(t *testing.T) {
	ifname := "test/model_nonewline.bin"
	inf, err := os.Open(ifname)
	defer inf.Close()
	if err != nil {
		t.Errorf("%q\n", err)
	}
	w2vtest(t, inf)
}

func w2vtest(t *testing.T, inf *os.File) {
	model, err := NewModel(inf)
	if model == nil {
		t.Errorf("Model was not loaded\t%q\n", err)
		return
	}
	t.Logf("Model\tvocab=%d\tvector=%d\n", model.GetVocabSize(), model.GetVectorSize())

	if wv := model.GetNormalizedVector("</s>"); wv == nil {
		t.Errorf("Error:\tNot found the first word\n")
	}
	if wv := model.GetNormalizedVector(`'Just`); wv == nil {
		t.Errorf("Error:\tNot found the last word\n")
	}

	wv1 := model.GetNormalizedVector("he")
	wv2 := model.GetNormalizedVector("she")
	wv3 := model.GetNormalizedVector("with")

	if n1 := wv1.GetNorm(); round(n1, 5) != 1 {
		t.Errorf("Error:\tNorm1 is not 1\t%f\n", n1)
	}
	if n2 := wv2.GetNorm(); round(n2, 5) != 1 {
		t.Errorf("Error:\tNorm2 is not 1\t%f\n", n2)
	}
	if n3 := wv3.GetNorm(); round(n3, 5) != 1 {
		t.Errorf("Error:\tNorm3 is not 1\t%f\n", n3)
	}

	if n := model.GetNorm("she"); n < 0 {
		t.Errorf("Error:\tNorm1 is negative\t%f\n", n)
	}
	if n := model.GetNorm("NOT_FOUND_WORD"); n != 0 {
		t.Errorf("Error:\tNorm1 is not 0\t%f\n", n)
	}

	wv4, norm4 := model.GetVector("it")
	if wv4 == nil {
		t.Errorf("Error:\tnil Vector of GetVector(\"it\")")
	}
	if round(norm4, 6) != 0.569891 {
		t.Errorf("Error:\tInorrect norm of GetVector(\"it\")")
	}
	wv, norm := model.GetVector("NOT_FOUND_WORD")
	if wv != nil {
		t.Errorf("Error:\tVector should be nil for non-existing word")
	}
	if norm != 0 {
		t.Errorf("Error:\tThe norm should be 0 for non-existing word")
	}

}

func round(f float32, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(float64(f)*shift+.5) / shift
}
