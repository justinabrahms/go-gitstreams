# go-gitstreams ![travis-ci build status](https://api.travis-ci.org/justinabrahms/go-gitstreams.png)

go-gitstreams is the Go portion of http://gitstreams.com/, a service
providing daily digests of GitHub activity. This repository is
generally responsible for querying information from the backing MySQL
datastore, formatting activity into aggregate format, and emailing
them out.

The original database schema was generated by Django's ORM. The
syncing piece of GitStreams is currently in Python and closed
source. An annotated schema can be found in `schema.sql`.

I'm currently interested in receiving contributions around fleshing
out templates, reducing the repetition in the various \*_render
methods, suggestions on testing and generally making the code more
idiomatic Go. Any suggestions on this front are happily
accepted. Anything *NOT* on listed will probably also be
accepted. Until I've moved away from Django as my ORM for the
closed-source portion, the schema (unfortunately) can't be greatly
reorganized, though new things can be added.