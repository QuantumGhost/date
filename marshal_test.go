// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"
	"time"
)

func TestGobEncoding(t *testing.T) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	decoder := gob.NewDecoder(&b)
	cases := []Date{
		New(-11111, time.February, 3),
		New(-1, time.December, 31),
		New(0, time.January, 1),
		New(1, time.January, 1),
		New(1970, time.January, 1),
		New(2012, time.June, 25),
		New(12345, time.June, 7),
	}
	for _, c := range cases {
		var d Date
		err := encoder.Encode(&c)
		if err != nil {
			t.Errorf("Gob(%v) encode error %v", c, err)
		} else {
			err = decoder.Decode(&d)
			if err != nil {
				t.Errorf("Gob(%v) decode error %v", c, err)
			} else if d != c {
				t.Errorf("Gob(%v) decode got %v", c, d)
			}
		}
	}
}

func TestDateJSONMarshalling(t *testing.T) {
	cases := []struct {
		value Date
		want  string
	}{
		{New(-11111, time.February, 3), `"-11111-02-03"`},
		{New(-1, time.December, 31), `"-0001-12-31"`},
		{New(0, time.January, 1), `"0000-01-01"`},
		{New(1, time.January, 1), `"0001-01-01"`},
		{New(1970, time.January, 1), `"1970-01-01"`},
		{New(2012, time.June, 25), `"2012-06-25"`},
		{New(12345, time.June, 7), `"+12345-06-07"`},
	}
	for _, c := range cases {
		var d Date
		bb, err := json.Marshal(c.value)
		if err != nil {
			t.Errorf("JSON(%v) marshal error %v", c, err)
		} else if string(bb) != c.want {
			t.Errorf("JSON(%v) == %v, want %v", c.value, string(bb), c.want)
		} else {
			err = json.Unmarshal(bb, &d)
			if err != nil {
				t.Errorf("JSON(%v) unmarshal error %v", c.value, err)
			} else if d != c.value {
				t.Errorf("JSON(%v) unmarshal got %v", c.value, d)
			}
		}
	}
}

func TestDateTextMarshalling(t *testing.T) {
	cases := []struct {
		value Date
		want  string
	}{
		{New(-11111, time.February, 3), "-11111-02-03"},
		{New(-1, time.December, 31), "-0001-12-31"},
		{New(0, time.January, 1), "0000-01-01"},
		{New(1, time.January, 1), "0001-01-01"},
		{New(1970, time.January, 1), "1970-01-01"},
		{New(2012, time.June, 25), "2012-06-25"},
		{New(12345, time.June, 7), "+12345-06-07"},
	}
	for _, c := range cases {
		var d Date
		bb, err := c.value.MarshalText()
		if err != nil {
			t.Errorf("Text(%v) marshal error %v", c, err)
		} else if string(bb) != c.want {
			t.Errorf("Text(%v) == %v, want %v", c.value, string(bb), c.want)
		} else {
			err = d.UnmarshalText(bb)
			if err != nil {
				t.Errorf("Text(%v) unmarshal error %v", c.value, err)
			} else if d != c.value {
				t.Errorf("Text(%v) unmarshal got %v", c.value, d)
			}
		}
	}
}

func TestDateBinaryMarshalling(t *testing.T) {
	cases := []struct {
		value Date
	}{
		{New(-11111, time.February, 3)},
		{New(-1, time.December, 31)},
		{New(0, time.January, 1)},
		{New(1, time.January, 1)},
		{New(1970, time.January, 1)},
		{New(2012, time.June, 25)},
		{New(12345, time.June, 7)},
	}
	for _, c := range cases {
		bb, err := c.value.MarshalBinary()
		if err != nil {
			t.Errorf("Binary(%v) marshal error %v", c, err)
		} else {
			var d Date
			err = d.UnmarshalBinary(bb)
			if err != nil {
				t.Errorf("Binary(%v) unmarshal error %v", c.value, err)
			} else if d != c.value {
				t.Errorf("Binary(%v) unmarshal got %v", c.value, d)
			}
		}
	}
}

func TestDateBinaryUnmarshallingErrors(t *testing.T) {
	var d Date
	err1 := d.UnmarshalBinary([]byte{})
	if err1 == nil {
		t.Errorf("unmarshal no empty data error")
	}

	err2 := d.UnmarshalBinary([]byte("12345"))
	if err2 == nil {
		t.Errorf("unmarshal no wrong length error")
	}
}

func TestInvalidDateText(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{`not-a-date`, `Date.ParseISO: cannot parse not-a-date: incorrect syntax`},
		{`215-08-15`, `Date.ParseISO: cannot parse 215-08-15: invalid year`},
	}
	for _, c := range cases {
		var d Date
		err := d.UnmarshalText([]byte(c.value))
		if err == nil || err.Error() != c.want {
			t.Errorf("InvalidText(%v) == %v, want %v", c.value, err, c.want)
		}
	}
}
