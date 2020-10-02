# msgapi

## 1) Install

```
make install
```

## 2) Config

Create `~/.config/msgapi/msgapirc` with mode `0600`:

```shell
USER=yOurDbUsER
PASS=YoURdBpAss
DB=msgapi # db name
PORT=3777
HOST=localhost # db host
```

## 3) Setup DB

```
make db
```

## 4) Run

```
msgapi
```

# TODO

* initrc script
* basic auth or jwt
* dockerize

