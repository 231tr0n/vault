package crypto_test

import (
	"testing"

	"github.com/231tr0n/vault/pkg/crypto"
)

func failTestCase(t *testing.T, i, o, w any) {
	t.Helper()
	t.Error("Input:", i, "|", "Output:", o, "|", "Want:", w)
}

func TestHmacHash(t *testing.T) {
	t.Parallel()

	pwd := []byte("test")

	tests := [][2][]byte{
		{
			[]byte("normalstringtohash"),
			[]byte("11069b0519652d9432f17c0e8ed81fffbc8571c9e9732aa98c7faf977d6ccf7c"),
		},
		{
			[]byte("normalstingtohash"),
			[]byte("5c5eb63eff8e6e15f83d5fb941d896fc8b496d00f19e51075c80da124a4e88b5"),
		},
	}

	for _, test := range tests {
		t.Log(test)
		out, err := crypto.HmacHash(test[0], pwd, nil)
		if err != nil {
			t.Fatal(err)
		}

		if !crypto.Verify(out, test[1]) {
			failTestCase(t, string(test[0]), string(out), string(test[1]))
		}
	}
}

func TestHash(t *testing.T) {
	t.Parallel()

	tests := [][2][]byte{
		{
			[]byte("normalstringtohash"),
			[]byte("1e7952a21b1ddb915cd04c178139aea0773568727b0d86a013c7dc9c467d3079"),
		},
		{
			[]byte("normalstingtohash"),
			[]byte("39bb86d3e642c133191a354ab5b687fa81512d358bf5667ce5ef396917ab706b"),
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

func TestHmacVerify(t *testing.T) {
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
		out := crypto.HmacVerify(test.a, test.s)

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
