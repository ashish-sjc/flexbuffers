package flexbuffers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFromJson(t *testing.T) {
	cases := []struct {
		input   string
		buildFn func(b *Builder)
	}{
		{
			// empty object
			input: `{}`,
			buildFn: func(b *Builder) {
				b.Map(func(b *Builder) {
				})
			},
		},
		{
			// empty array
			input: `[]`,
			buildFn: func(b *Builder) {
				b.Vector(false, false, func(bld *Builder) {
				})
			},
		},
		{
			// simple array
			input: `[1, 2, -3]`,
			buildFn: func(b *Builder) {
				b.Vector(false, false, func(b *Builder) {
					b.Int(1)
					b.Int(2)
					b.Int(-3)
				})
			},
		},
		{
			// nested array
			input: `["foo", [1, -3], 1.2]`,
			buildFn: func(b *Builder) {
				b.Vector(false, false, func(b *Builder) {
					b.StringValue("foo")
					b.Vector(false, false, func(b *Builder) {
						b.Int(1)
						b.Int(-3)
					})
					b.Float64(1.2)
				})
			},
		},
		{
			// nested array and object
			input: `{"aaa": "foo", "ccc": -456, "ddd": [7, 8, 9], "eee": {"fff": "xy"}, "bbb": 123}`,
			buildFn: func(b *Builder) {
				b.Map(func(b *Builder) {
					b.StringValueField([]byte("aaa"), "foo")
					b.IntField([]byte("ccc"), -456)
					b.VectorField([]byte("ddd"), false, false, func(b *Builder) {
						b.Int(7)
						b.Int(8)
						b.Int(9)
					})
					b.MapField([]byte("eee"), func(b *Builder) {
						b.StringValueField([]byte("fff"), "xy")
					})
					b.IntField([]byte("bbb"), 123)
				})
			},
		},
		{
			// string escape
			input: `["\"\n\\\u3042"]`,
			buildFn: func(b *Builder) {
				b.Vector(false, false, func(b *Builder) {
					b.StringValue("\"\n\\\u3042") // \u3042 == 'あ'
				})
			},
		},
	}
	for _, cas := range cases {
		r, err := FromJson(s2b(cas.input))
		if err != nil {
			t.Errorf("'%s': %v", cas.input, err)
			continue
		}
		b := NewBuilder()
		cas.buildFn(b)
		if err := b.Finish(); err != nil {
			t.Errorf("'%s': %v", cas.input, err)
			continue
		}
		if diff := cmp.Diff(b.Buffer(), r); diff != "" {
			t.Errorf("'%s': %s", cas.input, diff)
		}
	}
}
