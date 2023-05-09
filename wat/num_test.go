package wat

import (
	"testing"
)

func N(bits NumFlags, v ...string) Num {
	var num Num
	num.Flags = bits
	if len(v) > 0 {
		num.Integer = v[0]
	}
	if len(v) > 1 {
		num.Fraction = v[1]
	}
	if len(v) > 2 {
		num.Exponent = v[2]
	}
	return num
}

func TestNum_AppendTo(t *testing.T) {
	type testCase struct {
		Input    Num
		ExpectGo string
		Expect   string
	}

	testData := [...]testCase{
		{
			Input:    N(0, "0"),
			ExpectGo: `wat.Num{0, "0"}`,
			Expect:   "0",
		},
		{
			Input:    N(FlagSign, "0"),
			ExpectGo: `wat.Num{Sign, "0"}`,
			Expect:   "+0",
		},
		{
			Input:    N(FlagSign|FlagNeg, "0"),
			ExpectGo: `wat.Num{Sign|Neg, "0"}`,
			Expect:   "-0",
		},
		{
			Input:    N(FlagHex, "ff"),
			ExpectGo: `wat.Num{Hex, "ff"}`,
			Expect:   "0xff",
		},
		{
			Input:    N(FlagHex|FlagSign, "ff"),
			ExpectGo: `wat.Num{Hex|Sign, "ff"}`,
			Expect:   "+0xff",
		},
		{
			Input:    N(FlagHex|FlagSign|FlagNeg, "ff"),
			ExpectGo: `wat.Num{Hex|Sign|Neg, "ff"}`,
			Expect:   "-0xff",
		},
		{
			Input:    N(FlagFloat, "123", "456"),
			ExpectGo: `wat.Num{Float, "123", "456"}`,
			Expect:   "123.456",
		},
		{
			Input:    N(FlagFloat|FlagSign, "123", "456"),
			ExpectGo: `wat.Num{Float|Sign, "123", "456"}`,
			Expect:   "+123.456",
		},
		{
			Input:    N(FlagFloat|FlagSign|FlagNeg, "123", "456"),
			ExpectGo: `wat.Num{Float|Sign|Neg, "123", "456"}`,
			Expect:   "-123.456",
		},
		{
			Input:    N(FlagFloat|FlagHex|FlagSign|FlagNeg, "123", "456"),
			ExpectGo: `wat.Num{Float|Hex|Sign|Neg, "123", "456"}`,
			Expect:   "-0x123.456",
		},
		{
			Input:    N(FlagFloat|FlagSign|FlagNeg, "123", "456", "789"),
			ExpectGo: `wat.Num{Float|Sign|Neg, "123", "456", "789"}`,
			Expect:   "-123.456e789",
		},
		{
			Input:    N(FlagFloat|FlagSign|FlagNeg|FlagExpSign, "123", "456", "789"),
			ExpectGo: `wat.Num{Float|Sign|Neg|ExpSign, "123", "456", "789"}`,
			Expect:   "-123.456e+789",
		},
		{
			Input:    N(FlagFloat|FlagSign|FlagNeg|FlagExpSign|FlagExpNeg, "123", "456", "789"),
			ExpectGo: `wat.Num{Float|Sign|Neg|ExpSign|ExpNeg, "123", "456", "789"}`,
			Expect:   "-123.456e-789",
		},
		{
			Input:    N(FlagFloat|FlagHex|FlagSign|FlagNeg, "123", "456", "789"),
			ExpectGo: `wat.Num{Float|Hex|Sign|Neg, "123", "456", "789"}`,
			Expect:   "-0x123.456p789",
		},
		{
			Input:    N(FlagFloat|FlagHex|FlagSign|FlagNeg|FlagExpSign, "123", "456", "789"),
			ExpectGo: `wat.Num{Float|Hex|Sign|Neg|ExpSign, "123", "456", "789"}`,
			Expect:   "-0x123.456p+789",
		},
		{
			Input:    N(FlagFloat|FlagHex|FlagSign|FlagNeg|FlagExpSign|FlagExpNeg, "123", "456", "789"),
			ExpectGo: `wat.Num{Float|Hex|Sign|Neg|ExpSign|ExpNeg, "123", "456", "789"}`,
			Expect:   "-0x123.456p-789",
		},
		{
			Input:    N(FlagFloat | FlagInf),
			ExpectGo: `wat.Num{Float|Inf}`,
			Expect:   "inf",
		},
		{
			Input:    N(FlagFloat | FlagInf | FlagSign),
			ExpectGo: `wat.Num{Float|Inf|Sign}`,
			Expect:   "+inf",
		},
		{
			Input:    N(FlagFloat | FlagInf | FlagSign | FlagNeg),
			ExpectGo: `wat.Num{Float|Inf|Sign|Neg}`,
			Expect:   "-inf",
		},
		{
			Input:    N(FlagFloat | FlagNaN),
			ExpectGo: `wat.Num{Float|NaN}`,
			Expect:   "nan",
		},
		{
			Input:    N(FlagFloat | FlagNaN | FlagSign),
			ExpectGo: `wat.Num{Float|NaN|Sign}`,
			Expect:   "+nan",
		},
		{
			Input:    N(FlagFloat | FlagNaN | FlagSign | FlagNeg),
			ExpectGo: `wat.Num{Float|NaN|Sign|Neg}`,
			Expect:   "-nan",
		},
		{
			Input:    N(FlagFloat|FlagNaN|FlagAcanonical|FlagHex, "deadbeef"),
			ExpectGo: `wat.Num{Float|NaN|Acanonical|Hex, "deadbeef"}`,
			Expect:   "nan:0xdeadbeef",
		},
		{
			Input:    N(FlagFloat|FlagNaN|FlagAcanonical|FlagHex|FlagSign, "deadbeef"),
			ExpectGo: `wat.Num{Float|NaN|Acanonical|Hex|Sign, "deadbeef"}`,
			Expect:   "+nan:0xdeadbeef",
		},
		{
			Input:    N(FlagFloat|FlagNaN|FlagAcanonical|FlagHex|FlagSign|FlagNeg, "deadbeef"),
			ExpectGo: `wat.Num{Float|NaN|Acanonical|Hex|Sign|Neg, "deadbeef"}`,
			Expect:   "-nan:0xdeadbeef",
		},
	}

	for _, row := range testData {
		t.Run(row.ExpectGo, func(t *testing.T) {
			str := row.Input.GoString()
			if str != row.ExpectGo {
				t.Errorf("GoString: wrong output\n\texpect: %s\n\tactual: %s", row.ExpectGo, str)
			}

			str = row.Input.String()
			if str != row.Expect {
				t.Errorf("String: wrong output\n\texpect: %s\n\tactual: %s", row.Expect, str)
			}
		})
	}
}
