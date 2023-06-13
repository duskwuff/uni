package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/rangetable"
	"golang.org/x/text/unicode/runenames"
)

var (
	flagC  = flag.Bool("c", false, "")
	flagN  = flag.Bool("n", false, "")
	flagP  = flag.Bool("p", false, "")
	flagX  = flag.Bool("x", false, "")
	flag8  = flag.Bool("8", false, "")
	flag16 = flag.Bool("16", false, "")
)

var specialNames = map[rune]string{
	0x00: "NUL",
	0x01: "SOH (start of heading)",
	0x02: "STX (start of text)",
	0x03: "ETX (end of text)",
	0x04: "EOT (end of transmission)",
	0x05: "ENQ (enquiry)",
	0x06: "ACK (acknowledgement)",
	0x07: "BEL (bell)",
	0x08: "BS (backspace)",
	0x09: "HT (horizontal tab)",
	0x0a: "LF (line feed)",
	0x0b: "VT (vertical tab)",
	0x0c: "FF (form feed)",
	0x0d: "CR (carriage return)",
	0x0e: "SO (shift out)",
	0x0f: "SI (shift in)",
	0x10: "DLE (data link escape)",
	0x11: "DC1 (device control 1 / xon)",
	0x12: "DC2 (device control 2)",
	0x13: "DC3 (device control 3 / xoff)",
	0x14: "DC4 (device control 4)",
	0x15: "NAK (negative acknowledgement)",
	0x16: "SYN (synchronous idle)",
	0x17: "ETB (end of transmission block)",
	0x18: "CAN (cancel)",
	0x19: "EM (end of medium)",
	0x1a: "SUB (substitute)",
	0x1b: "ESC (escape)",
	0x1c: "FS (file separator)",
	0x1d: "GS (group separator)",
	0x1e: "RS (record separator)",
	0x1f: "US (unit separator)",
	0x7f: "DEL (delete)",
}

func usagefn() {
	fmt.Fprint(flag.CommandLine.Output(), ("Usage:\n" +
		"  uni [-n] <search>    search for codepoints with names matching <search>\n" +
		"  uni [-n] /regex/     search for codepoints with names matching regular expression /regex/\n" +
		"  uni [-p] U+<xxxx>    display codepoint U+<xxxx>\n" +
		"  uni [-c] <string>    display each codepoint in <string>\n" +
		"  uni -x <hex>         decode UTF-8 string from <hex> and display codepoints if valid\n" +
		"\n" +
		"Other flags:\n" +
		"  -8                   display UTF-8 sequences alongside codepoints\n" +
		"  -16                  display UTF-16 sequences alongside codepoints\n" +
		""),
	)
}

func countSet(flags ...bool) int {
	n := 0
	for i := range flags {
		if flags[i] {
			n++
		}
	}
	return n
}

func isASCII(x string) bool {
	for i := range x {
		if x[i] < 0x20 || x[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func main() {
	flag.Usage = usagefn
	flag.Parse()

	if countSet(*flagC, *flagN, *flagP, *flagX) > 1 {
		fmt.Fprintf(os.Stderr, "uni: -c / -n / -p / -x are all mutually exclusive")
		os.Exit(1)
	}

	arg := strings.Join(flag.Args(), " ")

	switch {
	case *flagC:
		doString(arg)
	case *flagN:
		doName(arg)
	case *flagP:
		for _, arg := range flag.Args() {
			doCodepoint(arg)
		}
	case *flagX:
		doHex(arg)

	case strings.HasPrefix(arg, "/"): // regex search
		doName(arg)

	case strings.HasPrefix(arg, "U+"): // codepoint reference
		for _, arg := range flag.Args() {
			doCodepoint(arg)
		}

	case len(arg) == 1 || !isASCII(arg):
		doString(arg)

	case len(arg) == 0:
		flag.Usage()
		os.Exit(1)

	default:
		doName(arg)
	}
}

func doCodepoint(p string) {
	if strings.HasPrefix(p, "U+") {
		cp, err := strconv.ParseInt(p[2:], 16, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to parse %q: %s\n", p, err)
			os.Exit(1)
		}
		showPoint(rune(cp))
	} else {
		fmt.Fprintf(os.Stderr, "unable to parse codepoint %q\n", p)
	}
}

func doHex(x string) {
	b, err := hex.DecodeString(strings.ReplaceAll(x, " ", ""))
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid hex input: %s\n", err)
		os.Exit(1)
	}

	*flag8 = true
	doBytes(b)
}

func doString(s string) {
	doBytes([]byte(s))
}

func doBytes(b []byte) {
	for i := 0; i < len(b); {
		r, n := utf8.DecodeRune(b[i:])
		if r == utf8.RuneError && n == 1 {
			fmt.Printf("\tU+????  (%02X)\n", b[i])
		} else {
			showPoint(r)
		}
		i += n
	}
}

func doName(search string) {
	var re *regexp.Regexp
	var err error
	if strings.HasPrefix(search, "/") {
		re, err = regexp.Compile(`(?i)` + strings.Trim(search, "/"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	} else {
		re, err = regexp.Compile(`(?i)\b` + regexp.QuoteMeta(search) + `\b`)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	rangetable.Visit(
		rangetable.Assigned(unicode.Version),
		func(r rune) {
			if re.MatchString(runenames.Name(r)) {
				showPoint(r)
			}
		},
	)
}

func showPoint(r rune) {
	if r > unicode.MaxRune {
		fmt.Printf("\tU+%X is not a valid Unicode character\n", r)
		return
	}

	ch := ""
	switch {
	case unicode.Is(unicode.Variation_Selector, r):
		// included in Mn but not what we wanted
	case unicode.Is(unicode.Mn, r): // nonspacing combining marks, e.g. combining acute
		ch = "\u25cc" + string(r) // U+25CC DOTTED CIRCLE
	case unicode.IsPrint(r):
		ch = string(r)
	}

	pt := fmt.Sprintf("%04X", r)

	name := specialNames[r]
	if name == "" {
		name = runenames.Name(r)
	}

	units := ""
	switch {
	case *flag8:
		ub := strings.Builder{}
		ub.WriteString(" (")
		for i, b := range []byte(string(r)) {
			if i > 0 {
				ub.WriteString(" ")
			}
			ub.WriteString(fmt.Sprintf("%02X", b))
		}
		ub.WriteString(")")
		units = fmt.Sprintf("%-11s", ub.String()) // len("XX XX XX XX") = 11

	case *flag16:
		if r < 0x10000 {
			units = fmt.Sprintf(" (%04X)     ", r) // padded to match width of long form
		} else {
			w1 := 0xD800 + (r >> 10) - 0x40
			w2 := 0xDC00 + (r & 0x3ff)
			units = fmt.Sprintf(" (%04X %04X)", w1, w2)
		}
	}

	fmt.Printf("%s\tU+%-5s%s\t%s\n", ch, pt, units, name)
}
