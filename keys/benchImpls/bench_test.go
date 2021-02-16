package bench

import "github.com/ijc/Gotty"
import "github.com/xo/terminfo"
import "testing"
import "path/filepath"
import "math/rand"
import "time"
import "os"

var keys []string = []string{ "kf1", "kclr", "kbs", "cl" }

func randomTerms() []string {
	var terms []string

	filepath.Walk("/usr/lib/terminfo/", func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			terms = append(terms, info.Name())
		}
		return nil
	})

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(terms), func(i, j int) { terms[i], terms[j] = terms[j], terms[i] })

	return terms
}

func BenchmarkIjc(b *testing.B) {
	terms := randomTerms()
	var err error
	var info *gotty.TermInfo
	var key string
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		info, err = gotty.OpenTermInfo(terms[i % len(terms)])
		if err != nil {
			b.Logf("#%d: Error reading %s: %s, skipping", i, terms[i % len(terms)], err.Error())
			continue
		}

		for _, key = range keys {
			info.GetAttribute(key)
		}
	}
}

func BenchmarkXo(b *testing.B) {
	terms := randomTerms()
	var info *terminfo.Terminfo
	var err error
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		info, err = terminfo.Load(terms[i % len(terms)])
		if err != nil {
			b.Logf("#%d: Error reading %s: %s, skipping", i, terms[i % len(terms)], err.Error())
			continue
		}

		_ = info.StringCaps()["key_f1"]
		_ = info.StringCaps()["key_clear"]
		_ = info.StringCaps()["key_bs"]
		_ = info.StringCaps()["clear"]
	}
}

/*func TestIjc(t *testing.T) {
	info, err := gotty.OpenTermInfo("termite")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	val, err := info.GetAttribute("kf1")

	if err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		t.Log(val)
	}
}*/

func TestXo(t *testing.T) {
	info, err := terminfo.Load("termite")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	val := info.StringCaps()["key_f1"]
	t.Log(val)
}
