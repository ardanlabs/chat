package cache_test

import (
	"net"
	"testing"

	"github.com/ardanlabs/chat/internal/platform/cache"
)

const succeed = "\u2713"
const failed = "\u2717"

// TestCache test that the caching system work.
func TestCache(t *testing.T) {
	cc := cache.New()

	id := "bill"
	tcpAddr := net.TCPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 6000,
		Zone: "Test",
	}
	address := tcpAddr.String()

	t.Log("Given the need to test caching.")
	{
		t.Logf("\tTest 0:\tBasic mechanics ID[ %s ] Address[ %s ]", id, address)
		{
			if err := cc.Add(id, &tcpAddr); err != nil {
				t.Fatalf("\t%s\tShould be able to add this client : %v\n", failed, err)
			}
			t.Logf("\t%s\tShould be able to add this client.\n", succeed)

			client, err := cc.GetID(id)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get this client by ID : %v\n", failed, err)
			}
			t.Logf("\t%s\tShould be able to get this client by ID.\n", succeed)

			if client.TCPAddr.Zone != tcpAddr.Zone {
				t.Errorf("\t%s\tShould be able to get the right zone.\n", failed)
				t.Errorf("\tWant[ %s ]\n", tcpAddr.Zone)
				t.Errorf("\tGot [ %s ]\n", client.TCPAddr.Zone)
				return
			}
			t.Logf("\t%s\tShould be able to get the right zone.\n", succeed)

			client, err = cc.GetAddress(address)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get this client by Address : %v\n", failed, err)
			}
			t.Logf("\t%s\tShould be able to get this client by Address.\n", succeed)

			if client.TCPAddr.Zone != tcpAddr.Zone {
				t.Errorf("\t%s\tShould be able to get the right zone.\n", failed)
				t.Errorf("\tWant[ %s ]\n", tcpAddr.Zone)
				t.Errorf("\tGot [ %s ]\n", client.TCPAddr.Zone)
				return
			}
			t.Logf("\t%s\tShould be able to get the right zone.\n", succeed)

			if err := cc.Remove(address); err != nil {
				t.Fatalf("\t%s\tShould be able to remove this client by Address : %v\n", failed, err)
			}
			t.Logf("\t%s\tShould be able to remove this client by Address.\n", succeed)

			if _, err = cc.GetAddress(address); err == nil {
				t.Errorf("\t%s\tShould NOT be able to get this client by Address : %v\n", failed, err)
			}
			t.Logf("\t%s\tShould NOT be able to get this client by Address.\n", succeed)

			if _, err = cc.GetID(id); err == nil {
				t.Errorf("\t%s\tShould NOT be able to get this client by ID : %v\n", failed, err)
			}
			t.Logf("\t%s\tShould NOT be able to get this client by ID.\n", succeed)
		}
	}
}
