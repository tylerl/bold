# Bold (command-line utility)

Colors matching lines of stdin. Sort of like grep, but all the lines are printed
and the matching lines are printed using the chosen foreground and background
color.

## Usage

    bold [-fg <color>] [-bg <color>] <regex>

**Valid colors include:**\
`(dark|light)-{black,red,green,yellow,blue,magenta,cyan,white}`\
Use "none" to explicitly select no color.

Regex can be specified as `-re <regex>` or simply as the last argument.

Examples:

    ps | bold bash
    ps | bold -re bash -fg black -bg light-white

Options:

    -bg string
        Background color
    -fg string
        Highlight color (default "light-red")
    -re string
        Regex to match
