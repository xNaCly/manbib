## issues

- https://github.com/jgm/pandoc/issues/8852 => toc generation not possible, workaround: see bottom

## Section numbers

1.  Executable programs or shell commands
2.  System calls (functions provided by the kernel)
3.  Library calls (functions within program libraries)
4.  Special files (usually found in /dev)
5.  File formats and conventions, e.g. /etc/passwd
6.  Games
7.  Miscellaneous (including macro packages and conventions), e.g. man(7), groff(7), man-pages(7)
8.  System administration commands (usually only for root)
9.  Kernel routines [Non standard]

## converting page to html

```bash
zcat /usr/share/man/man1/man.1.gz | pandoc --from man --to markdown | pandoc --toc --from markdown --to html5 --template template
```
