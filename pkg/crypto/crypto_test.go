package crypto_test

import (
	"testing"

	"github.com/231tr0n/vault/pkg/crypto"
)

func notok(t *testing.T, i any, o any, w any) {
	t.Error("Input:", i, "|", "Output:", o, "|", "Want:", w)
}

func TestHash(t *testing.T) {
	var p = []byte("passwdkey")

	var tests = [][2][]byte{
		{
			[]byte("normalstringtohash"),
			[]byte("f6be20978f067a1cf3ca91652c3d8d6855539bd726b227c9ce9b4feafd8225d1"),
		},
		{
			[]byte("normalstingtohash"),
			[]byte("be0c31346e35a3b9626c9c1385fda048e2ba530e2959d05127f40519d9c73bf1"),
		},
	}

	for _, val := range tests {
		var out, err = crypto.Hash(val[0], p, nil)
		if err != nil {
			t.Fatal(err)
		}
		if !crypto.Verify(out, val[1]) {
			notok(t, string(val[0]), string(out), string(val[1]))
		}
	}
}

func TestGenerate(t *testing.T) {
	var tests = []int{
		2,
		5,
		112,
	}

	for _, val := range tests {
		var out, err = crypto.Generate(val)
		if err != nil {
			t.Fatal(err)
		}
		if len(out) != val {
			notok(t, val, len(out), val)
		}
	}
}

func TestVerify(t *testing.T) {
	type test struct {
		s     []byte
		a     []byte
		check bool
	}

	var tests = []test{
		{
			s:     []byte("f6be20978f067a1cf3ca91652c3d8d6855539bd726b227c9ce9b4feafd8225d1"),
			a:     []byte("f6be20978f067a1cf3ca91652c3d8d6855539bd726b227c9ce9b4feafd8225d1"),
			check: true,
		},
		{
			s:     []byte("f6be20978f067a1cf3ca91652c3d8d6855539bd726b227c9ce9b4feafd8225d1"),
			a:     []byte("be0c31346e35a3b9626c9c1385fda048e2ba530e2959d05127f40519d9c73bf1"),
			check: true,
		},
	}

	for _, val := range tests {
		var out = crypto.Verify(val.a, val.s)
		if out != val.check {
			notok(
				t,
				[]string{
					string(val.s),
					string(val.a),
				},
				out,
				val.check,
			)
		}
	}
}
