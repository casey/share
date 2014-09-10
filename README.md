publish
=======

A content-addressible storage service. Coming soon!

[Test instance here, eventually.](http://rodarmor-publish.appspot.com)


API
---

`PUT /KEY DATA` Puts new data on the server. KEY must be the SHA-256 hash of DATA, which is the body of the request.

`GET /KEY` Will return the given data.


Quirks
------

Lots!
