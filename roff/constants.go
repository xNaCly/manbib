package roff

type Token struct {
	Kind    int
	Line    int
	Pos     int
	Content string
}

const (
	BOLD            = iota + 1 // B
	BOLDITALIC                 // BI
	BOLDROMAN                  // BR
	EXAMPLESTART               // EX
	EXAMPLEEND                 // EE
	ITALIC                     // I
	ITALICBOLD                 // IB
	ITALICROMAN                // IR
	INDENTPARAGRAPH            // IP
	LEFTPARAGRAPH              // LP
	MAILTOSTART                // MT
	MAILTOEND                  // ME
	CMDOPTION                  // OP
	PARAGRAPH                  // P
	ROMANBOLD                  // RB
	ROMANITALIC                // RI
	SMALL                      // SM
	SMALLBOLD                  // SB
	SECTIONHEADING             // SH
	SYNOPSISSTART              // SY
	SYNOPSISEND                // YS
	TITLEHEADING               // TH
	TAGGEDPARAGRAPH            // TP
	URLSTART                   // UR
	URLEND                     // UE
	SPACINGLINE                // .
	TEXT                       // anything else
	NEWLINE                    // \n
)
