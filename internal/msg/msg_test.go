package msg_test

import (
	"testing"

	"github.com/ardanlabs/chat/internal/msg"
)

const succeed = "\u2713"
const failed = "\u2717"

// TestEncode test that the encoding of a message works.
func TestEncode(t *testing.T) {
	tt := []struct {
		name   string
		m      msg.MSG
		length int
	}{
		{
			name: "length",
			m: msg.MSG{
				Sender:    "BillKenned",
				Recipient: "JillKenned",
				Data:      "hello",
			},
			length: 27,
		},
		{
			name: "shortname",
			m: msg.MSG{
				Sender:    "Bill",
				Recipient: "Cory",
				Data:      "helloworld",
			},
			length: 32,
		},
	}

	t.Log("Given the need to test encoding/decoding.")
	{
		for i, tst := range tt {
			t.Logf("\tTest %d:\t%s", i, tst.name)
			{
				data := msg.Encode(tst.m)
				if len(data) != tst.length {
					t.Fatalf("\t%s\tShould have the correct number of bytes : exp[%d] got[%d]\n", failed, tst.length, len(data))
				}
				t.Logf("\t%s\tShould have the correct number of bytes.\n", succeed)

				m := msg.Decode(data)
				if m.Sender != tst.m.Sender {
					t.Fatalf("\t%s\tShould have the correct Sender : exp[%v] got[%v]\n", failed, tst.m.Sender, m.Sender)
				}
				t.Logf("\t%s\tShould have the correct Sender.\n", succeed)

				if m.Recipient != tst.m.Recipient {
					t.Fatalf("\t%s\tShould have the correct Recipient : exp[%v] got[%v]\n", failed, tst.m.Recipient, m.Recipient)
				}
				t.Logf("\t%s\tShould have the correct Recipient.\n", succeed)

				if m.Data != tst.m.Data {
					t.Fatalf("\t%s\tShould have the correct data : exp[%s] got[%s]\n", failed, tst.m.Data, m.Data)
				}
				t.Logf("\t%s\tShould have the correct data.\n", succeed)
			}
		}
	}
}
