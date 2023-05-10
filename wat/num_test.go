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
		Name     string
		Input    Num
		ExpectGo string
		Expect   string
	}

	testData := [...]testCase{
		{
			Name:     "UnsignedInt",
			Input:    N(0, "0"),
			ExpectGo: `wat.Num{0, "0"}`,
			Expect:   "0",
		},
		{
			Name:     "PosInt",
			Input:    N(FlagSign, "0"),
			ExpectGo: `wat.Num{Sign, "0"}`,
			Expect:   "+0",
		},
		{
			Name:     "NegInt",
			Input:    N(FlagSign|FlagNeg, "0"),
			ExpectGo: `wat.Num{Sign|Neg, "0"}`,
			Expect:   "-0",
		},
		{
			Name:     "UnsignedHexInt",
			Input:    N(FlagHex, "ff"),
			ExpectGo: `wat.Num{Hex, "ff"}`,
			Expect:   "0xff",
		},
		{
			Name:     "PosHexInt",
			Input:    N(FlagHex|FlagSign, "ff"),
			ExpectGo: `wat.Num{Hex|Sign, "ff"}`,
			Expect:   "+0xff",
		},
		{
			Name:     "NegHexInt",
			Input:    N(FlagHex|FlagSign|FlagNeg, "ff"),
			ExpectGo: `wat.Num{Hex|Sign|Neg, "ff"}`,
			Expect:   "-0xff",
		},
		{
			Name:     "UnsignedFloat",
			Input:    N(FlagFloat, "123", "456"),
			ExpectGo: `wat.Num{Float, "123", "456"}`,
			Expect:   "123.456",
		},
		{
			Name:     "PosFloat",
			Input:    N(FlagFloat|FlagSign, "123", "456"),
			ExpectGo: `wat.Num{Float|Sign, "123", "456"}`,
			Expect:   "+123.456",
		},
		{
			Name:     "NegFloat",
			Input:    N(FlagFloat|FlagSign|FlagNeg, "123", "456"),
			ExpectGo: `wat.Num{Float|Sign|Neg, "123", "456"}`,
			Expect:   "-123.456",
		},
		{
			Name:     "NegHexFloat",
			Input:    N(FlagFloat|FlagHex|FlagSign|FlagNeg, "123", "456"),
			ExpectGo: `wat.Num{Float|Hex|Sign|Neg, "123", "456"}`,
			Expect:   "-0x123.456",
		},
		{
			Name:     "NegFloatWithUnsignedExp",
			Input:    N(FlagFloat|FlagSign|FlagNeg, "123", "456", "789"),
			ExpectGo: `wat.Num{Float|Sign|Neg, "123", "456", "789"}`,
			Expect:   "-123.456e789",
		},
		{
			Name:     "NegFloatWithPosExp",
			Input:    N(FlagFloat|FlagSign|FlagNeg|FlagExpSign, "123", "456", "789"),
			ExpectGo: `wat.Num{Float|Sign|Neg|ExpSign, "123", "456", "789"}`,
			Expect:   "-123.456e+789",
		},
		{
			Name:     "NegFloatWithNegExp",
			Input:    N(FlagFloat|FlagSign|FlagNeg|FlagExpSign|FlagExpNeg, "123", "456", "789"),
			ExpectGo: `wat.Num{Float|Sign|Neg|ExpSign|ExpNeg, "123", "456", "789"}`,
			Expect:   "-123.456e-789",
		},
		{
			Name:     "NegHexFloatWithUnsignedExp",
			Input:    N(FlagFloat|FlagHex|FlagSign|FlagNeg, "123", "456", "789"),
			ExpectGo: `wat.Num{Float|Hex|Sign|Neg, "123", "456", "789"}`,
			Expect:   "-0x123.456p789",
		},
		{
			Name:     "NegHexFloatWithPosExp",
			Input:    N(FlagFloat|FlagHex|FlagSign|FlagNeg|FlagExpSign, "123", "456", "789"),
			ExpectGo: `wat.Num{Float|Hex|Sign|Neg|ExpSign, "123", "456", "789"}`,
			Expect:   "-0x123.456p+789",
		},
		{
			Name:     "NegHexFloatWithNegExp",
			Input:    N(FlagFloat|FlagHex|FlagSign|FlagNeg|FlagExpSign|FlagExpNeg, "123", "456", "789"),
			ExpectGo: `wat.Num{Float|Hex|Sign|Neg|ExpSign|ExpNeg, "123", "456", "789"}`,
			Expect:   "-0x123.456p-789",
		},
		{
			Name:     "UnsignedInf",
			Input:    N(FlagFloat | FlagInf),
			ExpectGo: `wat.Num{Float|Inf}`,
			Expect:   "inf",
		},
		{
			Name:     "PosInf",
			Input:    N(FlagFloat | FlagInf | FlagSign),
			ExpectGo: `wat.Num{Float|Inf|Sign}`,
			Expect:   "+inf",
		},
		{
			Name:     "NegInf",
			Input:    N(FlagFloat | FlagInf | FlagSign | FlagNeg),
			ExpectGo: `wat.Num{Float|Inf|Sign|Neg}`,
			Expect:   "-inf",
		},
		{
			Name:     "UnsignedNaN",
			Input:    N(FlagFloat | FlagNaN),
			ExpectGo: `wat.Num{Float|NaN}`,
			Expect:   "nan",
		},
		{
			Name:     "PosNaN",
			Input:    N(FlagFloat | FlagNaN | FlagSign),
			ExpectGo: `wat.Num{Float|NaN|Sign}`,
			Expect:   "+nan",
		},
		{
			Name:     "NegNaN",
			Input:    N(FlagFloat | FlagNaN | FlagSign | FlagNeg),
			ExpectGo: `wat.Num{Float|NaN|Sign|Neg}`,
			Expect:   "-nan",
		},
		{
			Name:     "UnsignedAcanonicalNaN",
			Input:    N(FlagFloat|FlagNaN|FlagAcanonical|FlagHex, "deadbeef"),
			ExpectGo: `wat.Num{Float|NaN|Acanonical|Hex, "deadbeef"}`,
			Expect:   "nan:0xdeadbeef",
		},
		{
			Name:     "PosAcanonicalNaN",
			Input:    N(FlagFloat|FlagNaN|FlagAcanonical|FlagHex|FlagSign, "deadbeef"),
			ExpectGo: `wat.Num{Float|NaN|Acanonical|Hex|Sign, "deadbeef"}`,
			Expect:   "+nan:0xdeadbeef",
		},
		{
			Name:     "NegAcanonicalNaN",
			Input:    N(FlagFloat|FlagNaN|FlagAcanonical|FlagHex|FlagSign|FlagNeg, "deadbeef"),
			ExpectGo: `wat.Num{Float|NaN|Acanonical|Hex|Sign|Neg, "deadbeef"}`,
			Expect:   "-nan:0xdeadbeef",
		},
	}

	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
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
