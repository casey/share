Share
=====

Service for sharing.

Share only shares data that can be used freely by anyone.

So that it may do so, you must supply with every PUT request the license that covers the entirety of the contents of the PUT request in a header named "License".

Share will decline to store data under a license other than "Anyone may do anything with this."

[Test instance here.](http://rodarmor-share.appspot.com)

To avoid hosing the GAE free-tier datastore storage quota, PUTs are arbitrarily limited to 128 bytes.


API
---

`PUT /KEY` Share the body of the PUT request. The SHA-256 hash of the body must match KEY.

`GET /KEY` Get the shared data whose SHA-256 hash is equal to KEY.

`GET /KEY.EXT` Same as above, but sets the Content-Type header appropriately for the given EXT. If EXT is "sniff", Share will try to guess a Content-Type according to the [WHATWG MIME Sniffing standard](http://mimesniff.spec.whatwg.org).
