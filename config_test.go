package components

import (
	"encoding/json"
	"testing"
)

type Sample struct {
	unexported   int      `env:"UNEXPORTED, 110"` // 非导出字段
	Int          int      `env:"INT, 1"`
	Uint         uint     `env:"UINT, 1"`
	Bool         bool     `env:"BOOL, true"`
	String       string   `env:"STRING, Mariah Carey"`
	SliceInt     []int    `env:"SLICE_INT,1,2,3"`
	SliceString  []string `env:"SLICE_STRING,Fly, Like, A, Bird"`
	SubStruct    sub
	SubStructPtr *sub1
	Embed
}

func toString(s *Sample) string {
	b, _ := json.Marshal(s)
	return string(b)
}

type Embed struct {
	Array [2]bool `env:"ARRAY,0,1"`
}

type sub struct {
	SliceString []string `env:"SLICE_STRING_1,We, Belong, Together"`
}

type sub1 struct {
	String string
}

func TestMustMapConfig(t *testing.T) {
	cases := []struct {
		input  *Sample
		expect *Sample
		err    bool
		fn     []func(env, value string)
	}{
		{
			input:  &Sample{},
			expect: &Sample{},
			err:    true,
			fn:     nil,
		},
	}
	for _, c := range cases {
		err := MustMapConfig(c.input)
		if c.err {
			if err == nil {
				t.Fatal("error should not be nil")
			}
		} else if !c.err {
			if err != nil {
				t.Fatalf("error should be nil, got: %s", err)
			}
		} else {
			s1 := toString(c.input)
			s2 := toString(c.expect)
			if s1 != s2 {
				t.Fatalf("expect result: %s, got: %s", s2, s1)
			}
		}
	}
}
