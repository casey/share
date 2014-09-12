share
=====

Data sharing service.

[Test instance here.](http://rodarmor-share.appspot.com)


API
---

`PUT /KEY DATA` Puts new data on the server. KEY must be the SHA-256 hash of DATA, which is the body of the request. If supplied, EXT

`GET /KEY` Will return the given data.

`GET /KEY.EXT` Will return the given data with a content type appropriate for the given EXT. If EXT is "sniff", it will attempt to automatically detect the MIME type of the data.
