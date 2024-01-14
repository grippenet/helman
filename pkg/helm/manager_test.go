package helm

import (
	"testing"
)

func TestBindVars(t *testing.T) {
	s, e := bindVars("${a}/${b}", map[string]string{"a": "1", "b": "2"})
	if e != nil {
		t.Error(e)
	}
	if s != "1/2" {
		t.Errorf("Expected '1/2' got '%s'", s)
	}
	_, e = bindVars("${toto}", map[string]string{})
	if e == nil {
		t.Error("Expected error")
	}
}
