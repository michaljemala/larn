Larn is a roguelike game written by Noah Morgan in 1986.
This is a Go port of Larn, by David Symonds (dsymonds@golang.org).

Progress
--------
I have transliterated each original C source file.
I have done the minimal changes required to get it to build.
Some chunks have been stubbed out instead of ported (e.g. loading/saving),
but a substantial enough part works that you can play it.

Roadmap
-------
I will next be debugging, fixing and filling out the stubbed functions.

Next, I will abstract enough of the code to turn it into a web based
version, probably hosted on App Engine; you can be playing on one
web browser, sign in on another, and resume the same game. There will
be some kind of scoreboard too.

License
-------
The license of the original Larn by Noah Morgan is unclear. He posted
it to a public discussion group without stating the license; 1986 was
an era before licenses were "a thing". Given that, and given that I plan
to replace all the original C code, I am licensing this under the
[BSD 3-Clause Licence](http://www.opensource.org/licenses/bsd-license.php).
