uni - find and display Unicode characters
=========================================

Usage
-----

    Usage:
      uni [-n] <search>    search for codepoints with names matching <search>
      uni [-n] /regex/     search for codepoints with names matching regular expression /regex/
      uni [-p] U+<xxxx>    display codepoint U+<xxxx>
      uni [-c] <string>    display each codepoint in <string>
      uni -x <hex>         decode UTF-8 string from <hex> and display codepoints if valid

    Other flags:
      -8                   display UTF-8 sequences alongside codepoints
      -16                  display UTF-16 sequences alongside codepoints


Examples
--------

Search for codepoints containing a word in their name:

    % uni dog
    ⺨	U+2EA8 	CJK RADICAL DOG
    ⽝	U+2F5D 	KANGXI RADICAL DOG
    🌭	U+1F32D	HOT DOG
    🐕	U+1F415	DOG
    🐶	U+1F436	DOG FACE
    🦮	U+1F9AE	GUIDE DOG

Search for codepoints whose name matches a regular expression:

    % uni /^snow/
    ☃	U+2603 	SNOWMAN
    ⛄	U+26C4 	SNOWMAN WITHOUT SNOW
    ❄	U+2744 	SNOWFLAKE
    🏂	U+1F3C2	SNOWBOARDER
    🏔	U+1F3D4	SNOW CAPPED MOUNTAIN

Display a specific codepoint:

    % uni U+1f98a
    🦊	U+1F98A	FOX FACE

Decode a hexadecimal string as a sequence of codepoints:

    % uni -x 6b69cc81207475cc9bcca3
    k	U+006B  (6B)      	LATIN SMALL LETTER K
    i	U+0069  (69)      	LATIN SMALL LETTER I
    ◌́	U+0301  (CC 81)   	COMBINING ACUTE ACCENT
    	U+0020  (20)      	SPACE
    t	U+0074  (74)      	LATIN SMALL LETTER T
    u	U+0075  (75)      	LATIN SMALL LETTER U
    ◌̛	U+031B  (CC 9B)   	COMBINING HORN
    ◌̣	U+0323  (CC A3)   	COMBINING DOT BELOW

(Vietnamese for "characters".)

Display the codepoints which make up a string:

    % uni -c 🏳️‍🌈
    🏳	U+1F3F3	WAVING WHITE FLAG
    	U+FE0F 	VARIATION SELECTOR-16
    	U+200D 	ZERO WIDTH JOINER
    🌈	U+1F308	RAINBOW


Acknowledgements
----------------

`uni` draws heavy inspiration from the Perl utility [App::Uni](https://metacpan.org/pod/App::Uni).
