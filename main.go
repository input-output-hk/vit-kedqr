//$(which go) run $0 $@; exit $?

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/input-output-hk/vit_ked_qr/internal/bech32"
	"github.com/input-output-hk/vit_ked_qr/internal/ked"
	"github.com/skip2/go-qrcode"
)

/*
jcli key to-bytes
ed25519e_sk14rwkgpmmg5s29e4k8m4mny324lj4rv8x9tqg0tn5khlfqzgjt9ftj90u642j2skwraddf2qd88eqv8wv3a463mshgmz9dxtvthjswgqvcdwty
a8dd64077b4520a2e6b63eebb9922aafe551b0e62ac087ae74b5fe9009125952b915fcd5552542ce1f5ad4a80d39f2061dcc8f6ba8ee1746c456996c5de50720
*/

const (
	PIN_LENGTH = 4             // 4-digit numer
	KEY_HRP    = "ed25519e_sk" // bech32 hrp of ed25519extended secret key
	KEY_LENGTH = 64            // bytes of ed25519extended secret key
)

func main() {
	var (
		input  = flag.String("input", "", "path to file containing ed25519extended bech32 value")
		output = flag.String("output", "", "path to file to save qr code output, if not provided console output will be attempted")
		pin    = flag.String("pin", "", "Pin code. 4-digit number is used on Catalyst")
	)
	flag.Parse()

	/* PIN/PASSWORD check*/
	password := []byte(strings.TrimSpace(*pin))
	if len(password) != PIN_LENGTH {
		fmt.Fprintf(os.Stderr, "%s: needs %d digits, bud %d provided\n", "pin", PIN_LENGTH, len(password))
		os.Exit(1)
	}
	// convert pin/password digits from string to int value
	for i := range password {
		s := string(password[i])
		d, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: needs %d digits, but non digits detected\n", "pin", PIN_LENGTH)
			os.Exit(1)
		}
		password[i] = byte(d)
	}

	/* OUTPUT checks */
	*output = strings.TrimSpace(*output)

	/* INPUT checks */
	*input = strings.TrimSpace(*input)
	if *input == "" {
		fmt.Fprintf(os.Stderr, "flag needs an argument: -%s\n", "input")
		os.Exit(1)
	}

	// read input file data
	fileData, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error reading file %s - [%v]\n", "input", *input, err.Error())
		os.Exit(1)
	}

	// convert file data to string
	inputData := strings.TrimSpace(string(fileData))

	// bech32 decoding and convert to base256 (from base32)
	hrp, data, err := bech32.DecodeAndConvert(inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error %s\n", "bech32.DecodeAndConvert", err.Error())
		os.Exit(1)
	}

	// check hrp and key length expectations
	switch {
	case hrp != KEY_HRP:
		fmt.Fprintf(os.Stderr, "%s: error, expected %s, got %s\n", "HRP", KEY_HRP, hrp)
		os.Exit(1)
	case len(data) != KEY_LENGTH:
		fmt.Fprintf(os.Stderr, "%s: error, expected length %d, got %d\n", "KEY", KEY_LENGTH, len(data))
		os.Exit(1)
	}

	// encrypt the data
	encData, err := ked.Encrypt(password, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error, %s\n", "ked.Encrypt", err.Error())
		os.Exit(1)
	}

	// Prepare the qr code with hex encoded data
	qrData := hex.EncodeToString(encData)
	qrc, err := qrcode.New(qrData, qrcode.Medium)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error, %s\n", "qrcode.New", err.Error())
		os.Exit(1)
	}

	if *output != "" {
		err = qrc.WriteFile(-2, *output) // write to outfut file
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error, %s\n", "qrc.WriteFile", err.Error())
			fmt.Printf("\n%s\n", qrc.ToSmallString(false)) // output to console if write to file fails
		}
	} else {
		fmt.Printf("\n%s\n", qrc.ToSmallString(false)) // output to console in no output provided
	}
	/* we are done and all the data are outputed */

	/*********************************************************************/
	/* Perform some reverse checks, encoded data -> original bech32, jic */
	/*********************************************************************/

	// read hex string data used to build the qr code
	qrEncData, err := hex.DecodeString(qrData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error, %s\n", "hex.DecodeString", err.Error())
		os.Exit(1)
	}

	// decrypt the data
	qrDecData, err := ked.Decrypt(password, qrEncData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error, %s\n", "ked.Decrypt", err.Error())
		os.Exit(1)
	}

	// convert to bech32
	qrInputData, err := bech32.ConvertAndEncode(KEY_HRP, qrDecData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error, %s\n", "bech32.ConvertAndEncode", err.Error())
		os.Exit(1)
	}

	// check if the original bech32 ed25519extended key matches the decrypted output
	if inputData != qrInputData {
		fmt.Fprintf(os.Stderr, "Encryption - Decryption checks resulted in missmatch error. The outputs may be corrupted/wrong\n")
		os.Exit(1)
	}
}
