package json

import "testing"

func TestMarshal(t *testing.T) {

	type Test struct {
		A  string  `json:"a,omitempty"`
		AP *string `json:"ap,omitempty"`
		B  int     `json:"b,omitempty"`
		BP *int    `json:"bp,omitempty"`
	}

	s := ""
	_ = s
	n := 0
	if b, err := Marshal(&Test{
		A:  "",
		B:  0,
		BP: &n,
	}); err != nil {
		t.Error(err)
	} else {
		t.Log(string(b))
	}
}
