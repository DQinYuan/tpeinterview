package shellcmd

import "testing"

func TestGenRecord(t *testing.T) {
	record := genRecord(9)

	result := []string{
		"9",
		"99",
		"999",
		"9999",
		"99999",
		"999999",
		"9999999",
		"99999999",
		"999999999",
		"9999999999",
		"99999999999",
	}

	for i, v := range result{
		if record[i] != v{
			t.Errorf("gen result error")
		}
	}


}
