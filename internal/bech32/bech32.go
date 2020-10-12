package bech32

import (
	"fmt"
	"strings"
)

const (
	dataCharList           = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
	separatorChar          = '1'
	separatorLength        = 1
	checksumLength         = 6
	hrpMinLength           = 1
	hrpMaxLength           = 83
	hrpCharMinBound        = byte(33)
	hrpCharMaxBound        = byte(126)
	encodedStringMinLength = hrpMinLength + separatorLength + checksumLength // 8
)

const (
	OTHERCASE = 0
	UPPERCASE = 1
	LOWERCASE = 2
)

type uint5 = uint8

// DecodeAndConvert decodes a bech32 encoded string and converts it to base256 encoded bytes,
// returning the human-readable part and the data part excluding the checksum.
func DecodeAndConvert(bech string) (string, []byte, error) {
	hrp, data, err := Decode(bech)
	if err != nil {
		return "", nil, err
	}

	converted, err := ConvertBits(data, 5, 8, false)
	if err != nil {
		return "", nil, err
	}

	return hrp, converted, nil
}

// ConvertAndEncode converts a base256 encoded byte string to base32 encoded byte string
// and encodes the byte slice into a bech32 string with the human-readable part hrp
func ConvertAndEncode(hrp string, data []byte) (string, error) {
	converted, err := ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", err
	}

	return Encode(hrp, converted)
}

// Decode bech32 encoded string, returning the human-readable
// part and the data part excluding the checksum.
func Decode(bech string) (string, []uint5, error) {
	if len(bech) < encodedStringMinLength {
		return "", nil, fmt.Errorf("invalid bech32 string length %d", len(bech))
	}

	sepPos := strings.LastIndexByte(bech, separatorChar)
	if sepPos < hrpMinLength {
		return "", nil, fmt.Errorf("invalid bech32 string separator position %d", sepPos)
	}
	if sepPos+checksumLength >= len(bech) {
		return "", nil, fmt.Errorf("invalid bech32 string data length %d", (len(bech) - sepPos - 1))
	}

	hrp := bech[:sepPos]
	charCase, err := checkHrp(hrp)
	if err != nil {
		return "", nil, err
	}

	switch charCase {
	case UPPERCASE:
		hrp = strings.ToLower(bech[:sepPos])
	case LOWERCASE, OTHERCASE:
		hrp = bech[:sepPos]
	}

	// datapart
	for i := sepPos + 1; i < len(bech); i++ {
		if bech[i] < hrpCharMinBound || bech[i] > hrpCharMaxBound {
			return hrp, nil, fmt.Errorf("invalid character in string: '%c'", bech[i])
		}
	}

	// The characters must be either all lowercase or all uppercase.
	lower := strings.ToLower(bech)
	upper := strings.ToUpper(bech)
	if bech != lower && bech != upper {
		return hrp, nil, fmt.Errorf("mixed-case strings not allowed")
	}

	// work with the lowercase string from now on.
	bech = lower
	data := bech[sepPos+1:]

	decoded, err := toBytes(data)
	if err != nil {
		return hrp, nil, fmt.Errorf("failed converting data to bytes: %v", err)
	}

	if !verifyChecksum(hrp, decoded) {
		moreInfo := ""
		checksum := bech[len(bech)-6:]
		expected, err := toChars(createChecksum(hrp, decoded[:len(decoded)-6]))
		if err == nil {
			moreInfo = fmt.Sprintf("Expected %v, got %v.", expected, checksum)
		}
		return hrp, nil, fmt.Errorf("checksum failed. " + moreInfo)
	}

	// exclude the last 6 bytes (checksum)
	return hrp, decoded[:len(decoded)-6], nil
}

// Encode a byte slice into a bech32 string with the hrp.
// Note that the bytes must each encode 5 bits (base32).
func Encode(hrp string, data []uint5) (string, error) {
	// Calculate the checksum of the data and append it at the end.
	checksum := createChecksum(hrp, data)
	combined := append(data, checksum...)

	// The resulting bech32 string is the concatenation of the hrp, the
	// separator 1, data and checksum. Everything after the separator is
	// represented using the specified charset.
	dataChars, err := toChars(combined)
	if err != nil {
		return "", fmt.Errorf("unable to convert data bytes to chars: %v", err)
	}

	return hrp + string(separatorChar) + dataChars, nil
}

// toBytes converts each character in the string 'chars' to the value of the
// index of the correspoding character in 'dataCharList'.
func toBytes(data string) ([]uint5, error) {
	decoded := make([]uint5, 0, len(data))
	for i := 0; i < len(data); i++ {
		index := strings.IndexByte(dataCharList, data[i])
		if index < 0 {
			return nil, fmt.Errorf("invalid character not part of dataCharList: %v", data[i])
		}
		decoded = append(decoded, uint5(index))
	}

	return decoded, nil
}

// toChars converts the byte slice 'data' to a string where each byte in 'data'
// encodes the index of a character in 'dataCharList'.
func toChars(data []uint5) (string, error) {
	result := make([]byte, 0, len(data))
	for _, b := range data {
		if int(b) >= len(dataCharList) {
			return "", fmt.Errorf("invalid data byte: %v", b)
		}
		result = append(result, dataCharList[b])
	}

	return string(result), nil
}

// ConvertBits converts a byte slice where each byte is encoding fromBits bits,
// to a byte slice where each byte is encoding toBits bits.
func ConvertBits(data []byte, fromBits, toBits uint8, pad bool) ([]byte, error) {
	if fromBits < 1 || fromBits > 8 || toBits < 1 || toBits > 8 {
		return nil, fmt.Errorf("only bit groups between 1 and 8 allowed")
	}

	// The final bytes, each byte encoding toBits bits.
	var regrouped []byte

	// Keep track of the next byte we create and how many bits we have
	// added to it out of the toBits goal.
	nextByte := byte(0)
	filledBits := uint8(0)

	for _, b := range data {

		// Discard unused bits.
		b = b << (8 - fromBits)

		// How many bits remaining to extract from the input data.
		remFromBits := fromBits
		for remFromBits > 0 {
			// How many bits remaining to be added to the next byte.
			remToBits := toBits - filledBits

			// The number of bytes to next extract is the minimum of
			// remFromBits and remToBits.
			toExtract := remFromBits
			if remToBits < toExtract {
				toExtract = remToBits
			}

			// Add the next bits to nextByte, shifting the already
			// added bits to the left.
			nextByte = (nextByte << toExtract) | (b >> (8 - toExtract))

			// Discard the bits we just extracted and get ready for
			// next iteration.
			b = b << toExtract
			remFromBits -= toExtract
			filledBits += toExtract

			// If the nextByte is completely filled, we add it to
			// our regrouped bytes and start on the next byte.
			if filledBits == toBits {
				regrouped = append(regrouped, nextByte)
				filledBits = 0
				nextByte = 0
			}
		}
	}

	// We pad any unfinished group if specified.
	if pad && filledBits > 0 {
		nextByte = nextByte << (toBits - filledBits)
		regrouped = append(regrouped, nextByte)
		filledBits = 0
		nextByte = 0
	}

	// Any incomplete group must be <= 4 bits, and all zeroes.
	if filledBits > 0 && (filledBits > 4 || nextByte != 0) {
		return nil, fmt.Errorf("invalid incomplete group")
	}

	return regrouped, nil
}

// ----------------------------------------------------------------------------

// createChecksum from the human-readable and data parts
func createChecksum(hrp string, data []uint5) []uint5 {
	checksum := []uint5{0, 0, 0, 0, 0, 0}
	mod := polymod(append(append(expandHrp(hrp), data...), checksum...)) ^ 1
	for i := 0; i < 6; i++ {
		checksum[i] = uint5((mod >> (5 * (5 - uint32(i))))) & 31
	}

	return checksum
}

// verifyChecksum from the human-readable and data parts
func verifyChecksum(hrp string, data []uint5) bool {
	return polymod(append(expandHrp(hrp), data...)) == 1
}

// checkHrp validity. Returns the case of the HRP, if any.
func checkHrp(hrp string) (int, error) {
	hrpCharCase := -1

	// len [1-83]
	if len(hrp) < hrpMinLength || len(hrp) > hrpMaxLength {
		return hrpCharCase, fmt.Errorf("invalid hrp : hrp=%s", hrp)
	}

	hasLower := false
	hasUpper := false

	for i := 0; i < len(hrp); i++ {
		b := hrp[i]
		// b [33 - 126]
		if b < hrpCharMinBound || b > hrpCharMaxBound {
			return hrpCharCase, fmt.Errorf("invalid character : hrp[%d]=%s", i, string(b))
		}
		if b >= 'a' && b <= 'z' {
			hasLower = true
		} else if b >= 'A' && b <= 'Z' {
			hasUpper = true
		}
		if hasLower && hasUpper {
			return hrpCharCase, fmt.Errorf("mix case : hrp=%s", hrp)
		}
	}

	switch {
	case hasUpper && !hasLower:
		hrpCharCase = UPPERCASE
	case !hasUpper && hasLower:
		hrpCharCase = LOWERCASE
	case !hasUpper && !hasLower:
		hrpCharCase = OTHERCASE
	case hasUpper && hasLower:
		// Unreachable
	}

	return hrpCharCase, nil
}

// expandHrp into 5bit-bytes
func expandHrp(hrp string) []uint5 {
	exp := make([]uint5, 0, len(hrp)*2+1)
	for i := 0; i < len(hrp); i++ {
		exp = append(exp, uint5(hrp[i]>>5)) // Part 1: shift down 5 bits
	}
	exp = append(exp, 0) // Separator: 0 byte
	for i := 0; i < len(hrp); i++ {
		exp = append(exp, uint5(hrp[i]&31)) // Part 2: zero top 3 bits
	}

	return exp
}

// polymod calculation, BIP 173
func polymod(values []uint5) uint32 {
	var (
		gen = [...]uint32{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3} // magic generator
		top uint8
		chk uint32 = 1
	)
	for _, v := range values {
		top = uint8(chk >> 25)
		chk = (chk&0x1ffffff)<<5 ^ uint32(v)
		for i := 0; i < 5; i++ {
			if (top>>uint8(i))&1 == 1 {
				chk = chk ^ gen[i]
			}
		}
	}

	return chk
}
