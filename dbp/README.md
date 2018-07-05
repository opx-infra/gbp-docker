# dbp

A simplified wrapper for `gbp-docker`.

## Usage

```bash
$ dbp src/
```

Build artifacts are found in `pool/stretch-amd64/src/`.

## Building against additional package sources

```bash
$ export EXTRA_SOURCES="
deb     http://deb.openswitch.net/stretch unstable opx opx-non-free
deb-src http://deb.openswitch.net/stretch unstable opx"
$ dbp src/
```
