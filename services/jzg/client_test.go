package jzg

import "testing"

func TestSignPayloadMatchesDocumentExample(t *testing.T) {
	payload := map[string]interface{}{
		"body": map[string]interface{}{
			"barCode": "ADAHJTC5QA",
			"idNo":    "",
		},
		"header": map[string]interface{}{
			"distributorCode":      "JZ-1-1",
			"ticketMachineAccount": "",
			"timestamp":            float64(156843421678454),
			"terminalNo":           "",
			"locatePort":           "",
			"secretKey":            "old-value",
		},
	}
	got, err := signPayload(payload, "123456")
	if err != nil {
		t.Fatalf("signPayload returned error: %v", err)
	}
	want := "e3dfedcb973986d97162a4e02ce3ea6c"
	if got != want {
		t.Fatalf("signature mismatch: got %s want %s", got, want)
	}
}
