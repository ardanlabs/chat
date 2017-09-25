package msg_test

import (
	"testing"

	"github.com/ardanlabs/chat/internal/msg"
)

const succeed = "\u2713"
const failed = "\u2717"

// TestEncode test that the encoding of a message works.
func TestEncode(t *testing.T) {
	m := msg.MSG{
		Name: "0123456789",
		Data: "hello",
	}

	t.Log("Given the need to test encoding.")
	{
		t.Log("\tTest 0:\tWhen checking for basic message.")
		{
			data := msg.Encode(m)
			if len(data) != 17 {
				t.Fatalf("\t%s\tShould have the correct number of bytes : exp[17] got[%d]\n", failed, len(data))
			}
			t.Logf("\t%s\tShould have the correct number of bytes.\n", succeed)
		}
	}
}
