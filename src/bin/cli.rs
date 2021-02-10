use chain_crypto::bech32::Bech32;
use chain_crypto::{Ed25519Extended, SecretKey};
use vit_kedqr::KeyQrCode;
use std::{
    convert::From,
    error::Error,
    fmt,
    fs::OpenOptions,
    io::{BufRead, BufReader},
    num::ParseIntError,
    path::PathBuf,
    str::FromStr,
};
use structopt::StructOpt;

/// QCode CLI toolkit
#[derive(Debug, PartialEq, StructOpt)]
#[structopt(rename_all = "kebab-case")]
pub struct QRcodeApp {
    #[structopt(
        long = "input",
        parse(from_os_str),
        about = "path to file containing ed25519extended bech32 value"
    )]
    input: PathBuf,
    #[structopt(
        long = "output",
        parse(from_os_str),
        about = "path to file to save qr code output, if not provided console output will be attempted"
    )]
    output: Option<PathBuf>,
    #[structopt(
        long = "pin",
        parse(try_from_str),
        about = "Pin code. 4-digit number is used on Catalyst"
    )]
    pin: QRPin,
}

impl QRcodeApp {
    pub fn exec(self) -> Result<(), Box<dyn Error>> {
        let QRcodeApp { input, output, pin } = self;
        // open input key and parse it
        let key_file = OpenOptions::new()
            .create(false)
            .read(true)
            .write(false)
            .append(false)
            .open(&input)
            .expect("Could not open input file.");

        let mut reader = BufReader::new(key_file);
        let mut key_str = String::new();
        let _key_len = reader
            .read_line(&mut key_str)
            .expect("Could not read input file.");
        let sk = key_str.trim_end().to_string();

        let secret_key: SecretKey<Ed25519Extended> =
            SecretKey::try_from_bech32_str(&sk).expect("Malformed secret key.");
        // use parsed pin from args
        let pwd = pin.password;
        // generate qrcode with key and parsed pin
        let qr = KeyQrCode::generate(secret_key.clone(), &[pwd.0, pwd.1, pwd.2, pwd.3]);
        // process output
        // EVERYTHING IS WORKING UP TIL HERE
        // TODO: process output file to path when given, output to stdout when path is None.
        match output {
            Some(path) => {
                qr.write_svg(path).unwrap();
            }
            None => {
                // FIXME:
                // to match functionality with the go version,
                // when no path is given, the output is printed
                // to stdout.
            }
        }
        Ok(())
    }
}

#[derive(Debug, PartialEq)]
pub struct QRPin {
    password: (u8, u8, u8, u8),
}

#[derive(Debug)]
pub struct BadPinError {}

impl Error for BadPinError {}

impl fmt::Display for BadPinError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "The PIN must consist of 4 digits.")
    }
}

impl From<ParseIntError> for BadPinError {
    fn from(_error: ParseIntError) -> Self {
        BadPinError {}
    }
}

impl FromStr for QRPin {
    type Err = BadPinError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s.chars().count() {
            n if n == 4 => {
                let mut c = s.chars();
                let mut pwd = (0u8, 0u8, 0u8, 0u8);
                match c.next() {
                    Some(v) if v.is_ascii_digit() => {
                        pwd.0 = v.to_digit(10).unwrap() as u8;
                    }
                    _ => return Err(BadPinError {}),
                }
                match c.next() {
                    Some(v) if v.is_ascii_digit() => {
                        pwd.1 = v.to_digit(10).unwrap() as u8;
                    }
                    _ => return Err(BadPinError {}),
                }
                match c.next() {
                    Some(v) if v.is_ascii_digit() => {
                        pwd.2 = v.to_digit(10).unwrap() as u8;
                    }
                    _ => return Err(BadPinError {}),
                }
                match c.next() {
                    Some(v) if v.is_ascii_digit() => {
                        pwd.3 = v.to_digit(10).unwrap() as u8;
                    }
                    _ => return Err(BadPinError {}),
                }
                Ok(QRPin { password: pwd })
            }
            _ => Err(BadPinError {}),
        }
    }
}

pub fn main() {
    QRcodeApp::from_args().exec().unwrap_or_else(report_error)
}

// same as in JCli
fn report_error(error: Box<dyn Error>) {
    eprintln!("{}", error);
    let mut source = error.source();
    while let Some(sub_error) = source {
        eprintln!("  |-> {}", sub_error);
        source = sub_error.source();
    }
    std::process::exit(1)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parse_pin_successfully() {
        for (pin, pwd) in &[
            ("0000", (0, 0, 0, 0)),
            ("1123", (1, 1, 2, 3)),
            ("0002", (0, 0, 0, 2)),
        ] {
            let qr_pin = QRPin::from_str(pin).unwrap();
            assert_eq!(qr_pin, QRPin { password: *pwd })
        }
    }
    #[test]
    fn pins_that_do_not_satisfy_length_reqs_return_error() {
        for bad_pin in &["", "1", "11", "111", "11111"] {
            let qr_pin = QRPin::from_str(bad_pin);
            assert!(qr_pin.is_err(),)
        }
    }

    #[test]
    fn pins_that_do_not_satisfy_content_reqs_return_error() {
        for bad_pin in &["    ", " 111", "llll", "000u"] {
            let qr_pin = QRPin::from_str(bad_pin);
            assert!(qr_pin.is_err(),)
        }
    }
}
