// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dwarf_test

import (
	. "debug/dwarf";
	"debug/elf";
	"debug/macho";
	"testing";
)

var typedefTests = map[string]string{
	"t_ptr_volatile_int": "*volatile int",
	"t_ptr_const_char": "*const char",
	"t_long": "long int",
	"t_ushort": "short unsigned int",
	"t_func_int_of_float_double": "func(float, double) int",
	"t_ptr_func_int_of_float_double": "*func(float, double) int",
	"t_func_ptr_int_of_char_schar_uchar": "func(char, signed char, unsigned char) *int",
	"t_func_void_of_char": "func(char) void",
	"t_func_void_of_void": "func() void",
	"t_func_void_of_ptr_char_dots": "func(*char, ...) void",
	"t_my_struct": "struct my_struct {vi volatile int@0; x char@4 : 1@7; y int@4 : 4@27; array [40]long long int@8}",
	"t_my_union": "union my_union {vi volatile int@0; x char@0 : 1@7; y int@0 : 4@28; array [40]long long int@0}",
	"t_my_enum": "enum my_enum {e1=1; e2=2; e3=-5; e4=1000000000000000}",
	"t_my_list": "struct list {val short int@0; next *t_my_list@8}",
	"t_my_tree": "struct tree {left *struct tree@0; right *struct tree@8; val long long unsigned int@16}",
}

func elfData(t *testing.T, name string) *Data {
	f, err := elf.Open(name);
	if err != nil {
		t.Fatal(err);
	}

	d, err := f.DWARF();
	if err != nil {
		t.Fatal(err);
	}
	return d;
}

func machoData(t *testing.T, name string) *Data {
	f, err := macho.Open(name);
	if err != nil {
		t.Fatal(err);
	}

	d, err := f.DWARF();
	if err != nil {
		t.Fatal(err);
	}
	return d;
}


func TestTypedefsELF(t *testing.T) {
	testTypedefs(t, elfData(t, "testdata/typedef.elf"));
}

func TestTypedefsMachO(t *testing.T) {
	testTypedefs(t, machoData(t, "testdata/typedef.macho"));
}

func testTypedefs(t *testing.T, d *Data) {
	r := d.Reader();
	seen := make(map[string]bool);
	for {
		e, err := r.Next();
		if err != nil {
			t.Fatal("r.Next:", err);
		}
		if e == nil {
			break;
		}
		if e.Tag == TagTypedef {
			typ, err := d.Type(e.Offset);
			if err != nil {
				t.Fatal("d.Type:", err);
			}
			t1 := typ.(*TypedefType);
			var typstr string;
			if ts, ok := t1.Type.(*StructType); ok {
				typstr = ts.Defn();
			} else {
				typstr = t1.Type.String();
			}

			if want, ok := typedefTests[t1.Name]; ok {
				if _, ok := seen[t1.Name]; ok {
					t.Errorf("multiple definitions for %s", t1.Name);
				}
				seen[t1.Name] = true;
				if typstr != want {
					t.Errorf("%s:\n\thave %s\n\twant %s", t1.Name, typstr, want);
				}
			}
		}
		if e.Tag != TagCompileUnit {
			r.SkipChildren();
		}
	}

	for k := range typedefTests {
		if _, ok := seen[k]; !ok {
			t.Errorf("missing %s", k);
		}
	}
}
