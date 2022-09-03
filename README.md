# Neobank

## resources

```
https://medium.com/@annapeterson89/whats-the-point-of-golang-pointers-everything-you-need-to-know-ac5e40581d4d

https://medium.com/goingogo/why-use-testmain-for-testing-in-go-dafb52b406bc

https://stackoverflow.com/questions/38172661/what-is-the-meaning-of-and

```

## psql debugging

1. For NO KEY UPDATE -> Mark the current transaction as update (acquire lock) + tell psql that the UPDATE won't affect the key so it can open up
2. Code on how to monitor lock -> [Monitor Lock](https://wiki.postgresql.org/wiki/Lock_Monitoring)

## paseto vs jwt

- lesser flexibility on ciper suite and algo selection (known insecure algo still valid)
- trivial forgery
  - set algo to 'none'
  - change asymmetric also to symmetric algo (hacker sign token with the public key hence server will validate it as true)

## token generation shortcut

```
openssl rand -hex 64 | head -c 32 (128 char -> 32 char)
pwgen -s 15 1 # generate random string of length 15 (-s to reduce entropy)
```
