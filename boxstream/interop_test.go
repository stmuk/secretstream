package boxstream

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/cryptix/go/logging/logtest"
	"github.com/cryptix/go/proc"
)

const cnt = 1000

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestInterop_WriteToJS(t *testing.T) {
	var key [32]byte
	var nonce [24]byte

	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		t.Fatal(err)
	}

	node, err := proc.StartStdioProcess("node", logtest.Logger("unbox.js", t), "unbox_test.js",
		base64.StdEncoding.EncodeToString(key[:]),
		base64.StdEncoding.EncodeToString(nonce[:]),
	)
	if err != nil {
		t.Fatal(err)
	}

	w := NewBoxer(node, &nonce, &key)
	want := strings.Repeat("Hello, Tests!", cnt)
	if _, err := fmt.Fprintln(w, want); err != nil {
		t.Fatal(err)
	}

	s := bufio.NewScanner(node)
	for s.Scan() {
		got := s.Text()
		t.Log(got)
		if got == want {
			break
		}
		t.Error("test data missmatch! got:%q", got)
	}
}

func TestInterop_ReadFromJS(t *testing.T) {
	var key [32]byte
	var nonce [24]byte

	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		t.Fatal(err)
	}

	node, err := proc.StartStdioProcess("node", logtest.Logger("box.js", t), "box_test.js",
		base64.StdEncoding.EncodeToString(key[:]),
		base64.StdEncoding.EncodeToString(nonce[:]),
	)
	if err != nil {
		t.Fatal(err)
	}

	r := NewUnboxer(node, &nonce, &key)

	want := strings.Repeat("Hello, Tests!", cnt)
	if _, err := fmt.Fprintln(node, want); err != nil {
		t.Fatal(err)
	}

	s := bufio.NewScanner(r)
	for s.Scan() {
		got := s.Text()
		t.Log(got)
		if got == want {
			break
		}
		t.Error("test data missmatch! got:%q", got)
	}

}
