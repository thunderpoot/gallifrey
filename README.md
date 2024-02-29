# ![gallifrey_small](https://github.com/thunderpoot/gallifrey/assets/54200401/19631f67-84fb-46cf-aa49-95973e104033) Gallifrey

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

A simple tool for signing/verifying arbitrary data with ED25519 keys, written in Go. Takes data from `STDIN`. Generates new private/public keys `private_key.pem` and `public_key.pem` if they don't already exist in the working directory.

### Installation
1. Clone this repository
2. Compile the binary

```
$ go build gallifrey.go
```

### Usage
```
$ ./gallifrey <mode> [arguments]
Modes:
        sign
        verify <publicKey> <signature>
```

```
$ echo "foo bar" | ./gallifrey sign
Signature: w/eZA0DlOUn6mUgYvY7To6aNNGg+v7C9PnfWtlCiooT1yeKaCrK5jiDG3Au6y7q/s2rowRlJ8mU+Ad3ALr/2Bw==
Public Key: wAooJAYrCTI01ZL0LnOmL2ZZE9xHEuvgmL5GT5Kmjuc=
```

<sub>_You can set the key and signature as environment variables..._</sub>

```
$ echo "foo bar" > foo.txt
$ ./gallifrey verify $key $sig < foo.txt
### SIGNATURE OK ###
```
