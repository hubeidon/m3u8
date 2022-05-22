package aria2

import "testing"

func TestGetGlobalStat(t *testing.T) {
	a := NewAriaEngine()
	if a.GetGlobalStat().DownloadSpeed != 0 {
		t.Fail()
	}
	if a.GetGlobalStat().NumActive != 0 {
		t.Fail()
	}
}
