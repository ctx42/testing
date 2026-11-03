// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mocker

import "unicode"

// wideRunes are the runes Unicode 15.0 gives the East Asian Width "W"
// (wide) or "F" (fullwidth) property, which take two columns when shown.
// The ranges are every rune whose golang.org/x/text/width v0.43.0 Kind is
// EastAsianWide or EastAsianFullwidth, merged where consecutive; regenerate
// them when [unicode.Version] moves past 15.0.
var wideRunes = &unicode.RangeTable{
	R16: []unicode.Range16{
		{Lo: 0x1100, Hi: 0x115F, Stride: 1},
		{Lo: 0x231A, Hi: 0x231B, Stride: 1},
		{Lo: 0x2329, Hi: 0x232A, Stride: 1},
		{Lo: 0x23E9, Hi: 0x23EC, Stride: 1},
		{Lo: 0x23F0, Hi: 0x23F0, Stride: 1},
		{Lo: 0x23F3, Hi: 0x23F3, Stride: 1},
		{Lo: 0x25FD, Hi: 0x25FE, Stride: 1},
		{Lo: 0x2614, Hi: 0x2615, Stride: 1},
		{Lo: 0x2648, Hi: 0x2653, Stride: 1},
		{Lo: 0x267F, Hi: 0x267F, Stride: 1},
		{Lo: 0x2693, Hi: 0x2693, Stride: 1},
		{Lo: 0x26A1, Hi: 0x26A1, Stride: 1},
		{Lo: 0x26AA, Hi: 0x26AB, Stride: 1},
		{Lo: 0x26BD, Hi: 0x26BE, Stride: 1},
		{Lo: 0x26C4, Hi: 0x26C5, Stride: 1},
		{Lo: 0x26CE, Hi: 0x26CE, Stride: 1},
		{Lo: 0x26D4, Hi: 0x26D4, Stride: 1},
		{Lo: 0x26EA, Hi: 0x26EA, Stride: 1},
		{Lo: 0x26F2, Hi: 0x26F3, Stride: 1},
		{Lo: 0x26F5, Hi: 0x26F5, Stride: 1},
		{Lo: 0x26FA, Hi: 0x26FA, Stride: 1},
		{Lo: 0x26FD, Hi: 0x26FD, Stride: 1},
		{Lo: 0x2705, Hi: 0x2705, Stride: 1},
		{Lo: 0x270A, Hi: 0x270B, Stride: 1},
		{Lo: 0x2728, Hi: 0x2728, Stride: 1},
		{Lo: 0x274C, Hi: 0x274C, Stride: 1},
		{Lo: 0x274E, Hi: 0x274E, Stride: 1},
		{Lo: 0x2753, Hi: 0x2755, Stride: 1},
		{Lo: 0x2757, Hi: 0x2757, Stride: 1},
		{Lo: 0x2795, Hi: 0x2797, Stride: 1},
		{Lo: 0x27B0, Hi: 0x27B0, Stride: 1},
		{Lo: 0x27BF, Hi: 0x27BF, Stride: 1},
		{Lo: 0x2B1B, Hi: 0x2B1C, Stride: 1},
		{Lo: 0x2B50, Hi: 0x2B50, Stride: 1},
		{Lo: 0x2B55, Hi: 0x2B55, Stride: 1},
		{Lo: 0x2E80, Hi: 0x2E99, Stride: 1},
		{Lo: 0x2E9B, Hi: 0x2EF3, Stride: 1},
		{Lo: 0x2F00, Hi: 0x2FD5, Stride: 1},
		{Lo: 0x2FF0, Hi: 0x2FFB, Stride: 1},
		{Lo: 0x3000, Hi: 0x303E, Stride: 1},
		{Lo: 0x3041, Hi: 0x3096, Stride: 1},
		{Lo: 0x3099, Hi: 0x30FF, Stride: 1},
		{Lo: 0x3105, Hi: 0x312F, Stride: 1},
		{Lo: 0x3131, Hi: 0x318E, Stride: 1},
		{Lo: 0x3190, Hi: 0x31E3, Stride: 1},
		{Lo: 0x31F0, Hi: 0x321E, Stride: 1},
		{Lo: 0x3220, Hi: 0x3247, Stride: 1},
		{Lo: 0x3250, Hi: 0x4DBF, Stride: 1},
		{Lo: 0x4E00, Hi: 0xA48C, Stride: 1},
		{Lo: 0xA490, Hi: 0xA4C6, Stride: 1},
		{Lo: 0xA960, Hi: 0xA97C, Stride: 1},
		{Lo: 0xAC00, Hi: 0xD7A3, Stride: 1},
		{Lo: 0xF900, Hi: 0xFAFF, Stride: 1},
		{Lo: 0xFE10, Hi: 0xFE19, Stride: 1},
		{Lo: 0xFE30, Hi: 0xFE52, Stride: 1},
		{Lo: 0xFE54, Hi: 0xFE66, Stride: 1},
		{Lo: 0xFE68, Hi: 0xFE6B, Stride: 1},
		{Lo: 0xFF01, Hi: 0xFF60, Stride: 1},
		{Lo: 0xFFE0, Hi: 0xFFE6, Stride: 1},
	},
	R32: []unicode.Range32{
		{Lo: 0x16FE0, Hi: 0x16FE4, Stride: 1},
		{Lo: 0x16FF0, Hi: 0x16FF1, Stride: 1},
		{Lo: 0x17000, Hi: 0x187F7, Stride: 1},
		{Lo: 0x18800, Hi: 0x18CD5, Stride: 1},
		{Lo: 0x18D00, Hi: 0x18D08, Stride: 1},
		{Lo: 0x1AFF0, Hi: 0x1AFF3, Stride: 1},
		{Lo: 0x1AFF5, Hi: 0x1AFFB, Stride: 1},
		{Lo: 0x1AFFD, Hi: 0x1AFFE, Stride: 1},
		{Lo: 0x1B000, Hi: 0x1B122, Stride: 1},
		{Lo: 0x1B132, Hi: 0x1B132, Stride: 1},
		{Lo: 0x1B150, Hi: 0x1B152, Stride: 1},
		{Lo: 0x1B155, Hi: 0x1B155, Stride: 1},
		{Lo: 0x1B164, Hi: 0x1B167, Stride: 1},
		{Lo: 0x1B170, Hi: 0x1B2FB, Stride: 1},
		{Lo: 0x1F004, Hi: 0x1F004, Stride: 1},
		{Lo: 0x1F0CF, Hi: 0x1F0CF, Stride: 1},
		{Lo: 0x1F18E, Hi: 0x1F18E, Stride: 1},
		{Lo: 0x1F191, Hi: 0x1F19A, Stride: 1},
		{Lo: 0x1F200, Hi: 0x1F202, Stride: 1},
		{Lo: 0x1F210, Hi: 0x1F23B, Stride: 1},
		{Lo: 0x1F240, Hi: 0x1F248, Stride: 1},
		{Lo: 0x1F250, Hi: 0x1F251, Stride: 1},
		{Lo: 0x1F260, Hi: 0x1F265, Stride: 1},
		{Lo: 0x1F300, Hi: 0x1F320, Stride: 1},
		{Lo: 0x1F32D, Hi: 0x1F335, Stride: 1},
		{Lo: 0x1F337, Hi: 0x1F37C, Stride: 1},
		{Lo: 0x1F37E, Hi: 0x1F393, Stride: 1},
		{Lo: 0x1F3A0, Hi: 0x1F3CA, Stride: 1},
		{Lo: 0x1F3CF, Hi: 0x1F3D3, Stride: 1},
		{Lo: 0x1F3E0, Hi: 0x1F3F0, Stride: 1},
		{Lo: 0x1F3F4, Hi: 0x1F3F4, Stride: 1},
		{Lo: 0x1F3F8, Hi: 0x1F43E, Stride: 1},
		{Lo: 0x1F440, Hi: 0x1F440, Stride: 1},
		{Lo: 0x1F442, Hi: 0x1F4FC, Stride: 1},
		{Lo: 0x1F4FF, Hi: 0x1F53D, Stride: 1},
		{Lo: 0x1F54B, Hi: 0x1F54E, Stride: 1},
		{Lo: 0x1F550, Hi: 0x1F567, Stride: 1},
		{Lo: 0x1F57A, Hi: 0x1F57A, Stride: 1},
		{Lo: 0x1F595, Hi: 0x1F596, Stride: 1},
		{Lo: 0x1F5A4, Hi: 0x1F5A4, Stride: 1},
		{Lo: 0x1F5FB, Hi: 0x1F64F, Stride: 1},
		{Lo: 0x1F680, Hi: 0x1F6C5, Stride: 1},
		{Lo: 0x1F6CC, Hi: 0x1F6CC, Stride: 1},
		{Lo: 0x1F6D0, Hi: 0x1F6D2, Stride: 1},
		{Lo: 0x1F6D5, Hi: 0x1F6D7, Stride: 1},
		{Lo: 0x1F6DC, Hi: 0x1F6DF, Stride: 1},
		{Lo: 0x1F6EB, Hi: 0x1F6EC, Stride: 1},
		{Lo: 0x1F6F4, Hi: 0x1F6FC, Stride: 1},
		{Lo: 0x1F7E0, Hi: 0x1F7EB, Stride: 1},
		{Lo: 0x1F7F0, Hi: 0x1F7F0, Stride: 1},
		{Lo: 0x1F90C, Hi: 0x1F93A, Stride: 1},
		{Lo: 0x1F93C, Hi: 0x1F945, Stride: 1},
		{Lo: 0x1F947, Hi: 0x1F9FF, Stride: 1},
		{Lo: 0x1FA70, Hi: 0x1FA7C, Stride: 1},
		{Lo: 0x1FA80, Hi: 0x1FA88, Stride: 1},
		{Lo: 0x1FA90, Hi: 0x1FABD, Stride: 1},
		{Lo: 0x1FABF, Hi: 0x1FAC5, Stride: 1},
		{Lo: 0x1FACE, Hi: 0x1FADB, Stride: 1},
		{Lo: 0x1FAE0, Hi: 0x1FAE8, Stride: 1},
		{Lo: 0x1FAF0, Hi: 0x1FAF8, Stride: 1},
		{Lo: 0x20000, Hi: 0x3FFFF, Stride: 1},
	},
}

// lineWidth returns the number of columns the line takes when shown: a tab
// moves to the next multiple of [tabWidth], a wide rune takes two columns,
// and a combining mark none.
func lineWidth(line string) int {
	var width int
	for _, r := range line {
		switch {
		case r == '\t':
			width += tabWidth - width%tabWidth
		case unicode.In(r, unicode.Mn, unicode.Me):
			// A combining mark sits on the rune before it.
		case unicode.Is(wideRunes, r):
			width += 2
		default:
			width++
		}
	}
	return width
}
