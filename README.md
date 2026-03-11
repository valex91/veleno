## Veleno

As a learning excercise for golang I created Veleno, a MITM Proxy.

- Generate a local certificate with the shell script, which then add it to your local certs
- Poison the local DNS redirecting traffic to the listener
- Proxy traffic from the listener to an arbitrary destination
