# vit-kedqr

Tool to encrypt an ed25519extended private key and build/display a qr code with encryted data.
Used/needed for Catalyst project.

## Installation
- Download the latest release from releases page [https://github.com/input-output-hk/vit-kedqr/releases](https://github.com/input-output-hk/vit-kedqr/releases)
- or if you have [Rust](https://www.rust-lang.org/learn/get-started) installed use: `cargo install --git github.com/input-output-hk/vit-kedqr`
- if you want to build from the source use ie: `cargo build --bin vit-kedqr`

### Generating QR Codes with encrypted data
1. Generate an ed25519extended private key by using [jörmungandr](https://github.com/input-output-hk/jormungand) **jcli** binary
- `jcli key generate --type ed25519extended test-key.sk`
and you will get a new file `test-key.sk` containing the bech32 encoded key, ie: `ed25519e_sk14rwkgpmmg5s29e4k8m4mny324lj4rv8x9tqg0tn5khlfqzgjt9ftj90u642j2skwraddf2qd88eqv8wv3a463mshgmz9dxtvthjswgqvcdwty`

2. Use `vit-kedqr` binary to build your encrypted qr code with your provided pin code (4-digit number is expected for now)
- `vit-kedqr --input=test-key.sk --pin=1234 --output=qr-test-key.png` will output the qr code on the provided .png file **qr-test-key.png**.
- `vit-kedqr --input=test-key.sk --pin=1234` will do the same as the previous command, but insted of file the output will be printed out on the console ie:
```

█████████████████████████████████████████████████████████████████████
█████████████████████████████████████████████████████████████████████
████ ▄▄▄▄▄ ███▄█ ▀▄█▀▀▄█ ▄█ ▀ █▄█▀ ▀▀▄▀█▀▄▀▄█▀▄█▀▀▀▄  █  █ ▄▄▄▄▄ ████
████ █   █ █   ▀▀▀▀██ ▄▄█▀▄▄ ▀█▄▀▀▄▀ ▄█▄▀▀█▀▄▀ ▀▄██▄▀▀ ▄ █ █   █ ████
████ █▄▄▄█ █ ▀█▀▄▀▀▄▄█▄▄▀▀█▀▀▀▄  ▄▄▄ ▄█ ▀▀█ █▀ ██ ██▄▄ ▄██ █▄▄▄█ ████
████▄▄▄▄▄▄▄█ ▀ ▀ █▄▀▄▀ ▀ ▀ █ █▄▀ █▄█ █▄█▄█▄█ █ ▀ █▄▀▄█ ▀▄█▄▄▄▄▄▄▄████
████▄▀ ▄ ▄▄██   ▀▀▄█ ▀▄▄ █▄▄█▀▀█▄  ▄ ▄▄ ▄ ▄ ▀▄▀██ ▄▀ ██▀▄█  ▄▄▄▀█████
████▄▀▀▀ ▄▄ █▄▀██▄▄ ▀ ▄▀▄ ▄   █▀▄█▀▄▀▄▄▄▄▀▄█▀▄▀▄▀█▄█▄ ▄▄▄ ▀ ▀▄▄█ ████
██████ ▄ █▄██ █  ▄ ▀▄█ █▄██▀▄▀▀▀█▄█▀▀ ▄▄▄█▄▄▀█▀ █▀▄▄ ██▀▄▀▀█▀▄▄ ▄████
█████▄▀▀▄▄▄ █ ▄█▄▄█▄▄▄▄▄ ▀▄ █▄▀▄▄█ ▀█▀▄█ ▄▄▀█▄▀██▀▄█▄▄ █▄██▀▀█▄▀ ████
███████▄█▄▄██▀ █▀ ▄▀ █▀▀ ▀  ▀▀▀▀█ ▀█ ▀▄▄██▄▀ ▄ ▄ ██ ▄ █▄▄▀ ██ ▄▀▄████
█████▄██▄▀▄██▀█▄  █▀██ ▄ ▄ ▀ ▄ ▀▄▀  ▀█▄ █▄▄▀  ▀   ▄▄▄█▄ ▄█▀▀  ██▄████
████ █ █▀█▄ ▄█ ▀ ▀ ▀██▄  █▀▄▄▄▄ ██▀▀▀ █▀ █▄█▀▄ ███▄▄ ▄▄▀▄▄ █▀▀▄▀█████
██████▄█▄ ▄▀██▀███▀▄ ▄█▀▄  ██▄█▄▀▄ ▀██ ▀▄  █▀ █ ▀▄ ▀▄ ▄██ ▀▀▀▀▄▀ ████
████▄▀ ▄█▄▄▀▀▄▀▄ ▄▄▀█▄█ █ █▀▄▀▀▄ ▄ ▀▄▄ ▀▄▀  █▀█▀▀▄ █▄██▀▄▄ █▀█▄▄▄████
████▄███▄▀▄ ▀▄▄▄██▀█▄ ▀ ▀▀█▀█▄▀▄▄▄▀▀█▄█▀▄▀█▄ ▀  ▀▄█ ▄▄ ▄█ ▀▄▀ █  ████
█████▀▀  ▄▄▄ ▄█▀██▄ █ ▀ ▄▀█▄▀▀▀▄ ▄▄▄  █▀ ▄█ █▄ ▀█ ▄▀▄ ██ ▄▄▄ █▄▄▄████
████ ▀▀  █▄█ ██▄ ▄ ▄ ▄▀ █ ▄█▀ ▀  █▄█ ▀ ▄▄▄▄█▀▄██▀█▄▀▄█▄  █▄█  ▄▄ ████
████▀█▀▄  ▄▄▄▀▀ ▄█▄▀▀ ▄▀▄▄▀▄▀▀▀▄  ▄▄▄ █▄▄█▄ █▀  █▀▄█ ▄▄▄  ▄ ▄█▄▄▄████
████▄▄█  ▄▄▄▄█▄▄██▄█▀▀▀█▄▄   ▀▀ ▄  ▀▄█▄▀ ▄▄▄███ █▀▄▀▄ ▄▄▄▄▄ ▄ ▄  ████
████▀▄█▀▄▄▄▄▀█  ▄▀▄▀█▄▀█▀▄▀▀▀▄▀▄▀▀▄▄▄██▀  █▄█  ▀█▀█▀ ██▄▀█▄▀▄▀▄  ████
█████▄▀▀▀█▄▀█▀ ▀██▀▀▄███ ██▄▀▀█▄▄█▄▄▄ ▄▀▄█ ▄▀██▀▀▀▄▀▄ ▄▄▄▀▄ ▄▀█▄▀████
████ ▄█▀  ▄█ █ ▄ ▄   ▀█▄▀ █▄ ▀█▄▄▀▄▄█▀▄▀▄ █ ▀█ ▀▀ ▄█ ▄█▄▄█▄▄▄█▄█▄████
████▄▄▄ █▄▄█  ▀ ▀▄▄▀▀█ █▄▄█▄  █▄▀▀▄▄▄▀▄ ▄█▄▄▀▄▀ ▀▀▄▄▄▄▄▄▀▄▄█▄█▄▄ ████
████ █ ▄█ ▄▀▀▀   █▀▄▄▀█▄▄ ▄▄▀▀▀ ▄▀▄▀█▀██ ▄█ █▄ █▀▀█ ▄▀█ ▄  ▄▄  ▄▄████
████ ▄██  ▄ █▄▀▀  ▀▄▄▀▀▄ ▀▄▄ ▄█ ▄▄█▄▄███▄▄▀ ▀   ▀▀█▀▄█▄█▄▄▄▄█▀█▄ ████
████▄ ██▄▄▄▄ ▀▀▀ █ ▀ ▀▄█▀▀▄▀▀▄▀▄▀▄▀▄███▄ ▀██▀█ ██▄██▄▄▄▄█▄▄▀▄▀▄ ▄████
████▀▀ ▄ ▄▄  █ ▀▀▀▄▀ ▀█▀█▄▄█ ▄▀ ██▄▀▀ ▄▄▄█▄█▀▄▀ ▀ ▄▀▄▄▄▄▄▀▄ ██▄█▄████
████▄▄▄▄██▄▄ ▀▀▀ █▀██ ▄ ▀▄▀▀▀▀▀▀ ▄▄▄ █▄▄▄█▄ ▀█▀█▀█▄█▄▄█  ▄▄▄    ▄████
████ ▄▄▄▄▄ █▀▄▄ ▄ ▀█▄█ ▄█▄ ▀▀▀██ █▄█ ▄▄█  ▄ █▄▀██   ▄▄▄▀ █▄█ ▄█▀▄████
████ █   █ █   █▄  ▄█▀▀███ ▀▀▀▀ ▄  ▄ ██▄███▀ █▀▄ ▀██▄▀█      ▄▄█ ████
████ █▄▄▄█ █▄ ▀ ▀ ▀ █ ▀█  ██▀▄  ▀██   ▄ █▄▄█  ▀  █▄█▄█▄▀▀▀▀██ ▄█ ████
████▄▄▄▄▄▄▄█▄█▄▄▄▄████▄██▄██▄▄▄███▄██▄██▄▄▄▄█▄███▄█▄▄█▄▄███▄██▄▄▄████
█████████████████████████████████████████████████████████████████████
▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
```

3. Scan the qr code with the Catalyst voting app.


### CLI Usage
```
vit-kedqr 0.1.0
QCode CLI toolkit

USAGE:
    vit-kedqr [OPTIONS] --input <input> --pin <pin>

FLAGS:
    -h, --help       Prints help information
    -V, --version    Prints version information

OPTIONS:
    -i, --input <input>      Path to file containing ed25519extended bech32 value
    -o, --output <output>    Path to file to save qr code output, if not provided console output will be attempted
    -p, --pin <pin>          Pin code. 4-digit number is used on Catalyst
```

#### Warning

The command-line arguments `-input`, `-output`, and `-pin`, are no longer valid. Instead, use the short or long form, as described in the help documentation.

To view the help documentation, execute `vit-kedqr -h` or `vit-kedqr --help` from the command-line.
