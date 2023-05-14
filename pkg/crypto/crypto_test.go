package crypto_test

import (
	"testing"

	"github.com/231tr0n/vault/pkg/crypto"
)

func failTestCase(t *testing.T, i, o, w any) {
	t.Helper()
	t.Error("Input:", i, "|", "Output:", o, "|", "Want:", w)
}

func TestHash(t *testing.T) {
	t.Parallel()

	tests := [][2][]byte{
		{
			[]byte("normalstringtohash"),
			[]byte("f6be20978f067a1cf3ca91652c3d8d6855539bd726b227c9ce9b4feafd8225d1"),
		},
		{
			[]byte("normalstingtohash"),
			[]byte("be0c31346e35a3b9626c9c1385fda048e2ba530e2959d05127f40519d9c73bf1"),
		},
	}

	for _, test := range tests {
		t.Log(test)
		out, err := crypto.Hash(test[0], nil)
		if err != nil {
			t.Fatal(err)
		}

		if !crypto.Verify(out, test[1]) {
			failTestCase(t, string(test[0]), string(out), string(test[1]))
		}
	}
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	tests := []int{
		2,
		5,
		112,
	}

	for _, test := range tests {
		t.Log(test)
		out, err := crypto.Generate(test)
		if err != nil {
			t.Fatal(err)
		}

		if len(out) != test {
			failTestCase(t, test, len(out), test)
		}
	}
}

func TestVerify(t *testing.T) {
	t.Parallel()

	type test struct {
		s     []byte
		a     []byte
		check bool
	}

	tests := []test{
		{
			s:     []byte("f6be20978f067a1cf3ca91652c3d8d6855539bd726b227c9ce9b4feafd8225d1"),
			a:     []byte("f6be20978f067a1cf3ca91652c3d8d6855539bd726b227c9ce9b4feafd8225d1"),
			check: true,
		},
		{
			s:     []byte("f6be20978f067a1cf3ca91652c3d8d6855539bd726b227c9ce9b4feafd8225d1"),
			a:     []byte("be0c31346e35a3b9626c9c1385fda048e2ba530e2959d05127f40519d9c73bf1"),
			check: false,
		},
	}

	for _, test := range tests {
		t.Log(test)
		out := crypto.Verify(test.a, test.s)

		if out != test.check {
			failTestCase(
				t,
				[]string{
					string(test.s),
					string(test.a),
				},
				out,
				test.check,
			)
		}
	}
}

func TestEncryptAndDecrypt(t *testing.T) {
	t.Parallel()

	tests := [][2][]byte{
		{
			[]byte("hi"),
			[]byte("k"),
		},
		{
			[]byte("h"),
			[]byte("key"),
		},
		{
			[]byte("hihowareyou"),
			[]byte("keykeykekeykeykekeykeykekeykeyke"),
		},
		{
			[]byte("hihihihihihihihihihihihihihihihi"),
			[]byte("hi"),
		},
	}

	for _, test := range tests {
		t.Log(test)
		eout, err := crypto.Encrypt(test[0], test[1])
		if err != nil {
			t.Fatal(err)
		}

		var dout []byte

		dout, err = crypto.Decrypt(eout, test[1])
		if err != nil {
			t.Fatal(err)
		}

		if string(dout) != string(test[0]) {
			failTestCase(
				t,
				[]string{
					string(test[0]),
					string(test[1]),
				},
				string(dout),
				string(test[0]),
			)
		}
	}
}
