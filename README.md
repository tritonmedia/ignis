<p align="center">
  <img width=500 height=300 src="https://78.media.tumblr.com/202569ad66cd8041e0e1bd601dafd1a9/tumblr_p5a4fmbp0M1rblqwco2_1280.png" alt="Ignis in a car" />
</p>

<p align="center">
  <b>ignis</b>
</p>

<p align="center"> A Telegram Bot frontend to Triton Media.</p>


## What is this?

Ignis a friendly, helpful, Telegram bot that will assist users in creating media on the Triton Media platform over Telegram.

## Installation / Hacking


1. Install it (setup a GOPATH)

```bash
if [[ -z "$GOPATH" ]]; then
  echo "Creating a temporary GOPATH at ~/go, add this to your ~/.<shell>rc to make it permanent:"
  echo '  export GOPATH="$HOME/go"'
  GOPATH="$HOME/go"
  mkdir -p "$GOPATH"
fi

mkdir -p "$GOPATH/tritonmedia"; cd "$GOPATH/tritonmedia"
git clone git@github.com:tritonmedia/ignis ignis; cd "ignis"
make dep
make
```

2. Configure it via the config.yaml file in `./config` (see `./config/config.example.yaml`)

3. Run it

```bash
# in the ignis directory
./bin/ignis
```