package keys

import "testing"

import "unicode/utf8"

func TestDelete(t *testing.T) {
	var in []byte = []byte{ 0x1B, 0x5B, 0x33, 0x7E }
	var should []Key = []Key{
		Key{KeySpecial, SpecialDelete, utf8.RuneError},
	}

	var parser *Parser
	var err error
	parser, err = Init()
	if err != nil {
		t.Error("Parser initialization failed:")
		t.Error(err.Error())
		t.FailNow()
	}

	var is []Key = parser.ParseUTF8(in)

	compareKeys(t, is, should)

	// SUCCESS
}

func TestSpecial(t *testing.T) {
	var in []byte = []byte{ 0x7F, 0x1B, 0x5B, 0x33, 0x7E, 0x61, 0x71 }
	//                        bs,  x1B,    [,    3,    ~,    a,    q
	var should []Key = []Key{
		Key{KeySpecial, SpecialBackspace, utf8.RuneError},
		Key{KeySpecial, SpecialDelete, utf8.RuneError},
		Key{KeyLetter, 0, 'a'},
		Key{KeyLetter, 0, 'q'},
	}

	var parser *Parser
	var err error
	parser, err = Init()
	if err != nil {
		t.Error("Parser initialization failed:")
		t.Error(err.Error())
		t.FailNow()
	}

	t.Log("Hello world")
	var is []Key = parser.ParseUTF8(in)
	t.Log("Hello world")

	compareKeys(t, is, should)

	// SUCCESS
}

func TestArrowKeys(t *testing.T) {
	var in []byte = []byte{ 0x1B, 0x5B, 0x44, 0x1B, 0x5B, 0x43 }
	//                       x1b,    [,    D,  x1b,    [,    C
	var should []Key = []Key{
		Key{KeySpecial, SpecialArrowLeft, utf8.RuneError},
		Key{KeySpecial, SpecialArrowRight, utf8.RuneError},
		/*Key{KeySpecial, SpecialArrowUp, utf8.RuneError},
		Key{KeySpecial, SpecialArrowDown, utf8.RuneError},*/
	}

	var parser *Parser
	var err error
	parser, err = Init()
	if err != nil {
		t.Error("Parser intialization failed:")
		t.Error(err.Error())
		t.FailNow()
	}

	var is []Key = parser.ParseUTF8(in)

	compareKeys(t, is, should)
}

func compareKeys(t *testing.T, is []Key, should []Key) {
	t.Helper()

	if len(should) != len(is) {
		t.Errorf("Wrong length: is %d, should be %d", len(is), len(should))
	}

	var end int = len(should)
	if len(is) < end {
		end = len(is)
	}

	var sK, iK Key
	var i int
	for i = 0; i < end; i++ {
		sK = should[i]
		iK = is[i]
		if sK.Type != iK.Type {
			t.Errorf("%d: Wrong type: is %d, should be %d", i, iK.Type, sK.Type)
		}
		if sK.Mod != iK.Mod {
			t.Errorf("%d: Wrong modifier: is %d, should be %d", i, iK.Mod, sK.Mod)
		}
		if sK.Value != iK.Value {
			t.Errorf("%d: Wrong value: is %c, should be %c", i, iK.Value, sK.Value)
		}
	}

	for ; i < len(is); i++ {
		t.Errorf("Extra key: %d, %d, %c", is[i].Type, is[i].Mod, is[i].Value)
	}

	for ; i < len(should); i++ {
		t.Errorf("Missing: %d, %d, %c", should[i].Type, should[i].Mod, should[i].Value)
	}
}
