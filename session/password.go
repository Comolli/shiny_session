package session

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"strings"
)

const (
	PasswordOK                = iota // Password passes our rules.
	PasswordTooShort                 // Password is too short.
	PasswordIsAName                  // Password is one of the predefined names.
	PasswordWasCompromised           // Password was found in a list of compromised passwords.
	PasswordFoundInDictionary        // Password was found in a dictionary.
	PasswordRepetitive               // Password consists of just repetetive characters.
	PasswordSequential               // Password consists of a simple sequence.
)

func initPasswords() {
	uncompress := func(compressed string) []string {
		decoded, _ := base64.StdEncoding.DecodeString(strings.Replace(compressed, "\n", "", -1))
		reader, _ := gzip.NewReader(bytes.NewReader(decoded))
		uncompressed, _ := ioutil.ReadAll(reader)
		return strings.Split(string(uncompressed), "\n")
	}

	dictionary = uncompress(dictionaryCompressed)
	commonPasswords = uncompress(commonPasswordsCompressed)
}
func ReasonablePassword(password string, names []string) int {
	if len(password) < 8 {
		return PasswordTooShort
	}
	for _, word := range names {
		if strings.ToLower(password) == strings.ToLower(word) {
			return PasswordIsAName
		}
	}
	for _, word := range commonPasswords {
		if password == word {
			return PasswordWasCompromised
		}
	}
	for _, word := range dictionary {
		if password == word {
			return PasswordFoundInDictionary
		}
	}
	var first rune
	for index, ch := range password {
		if index == 0 {
			first = ch
		} else {
			if ch != first {
				first = 0
				break
			}
		}
	}
	if first != 0 {
		return PasswordRepetitive
	}
	for _, sequence := range []string{
		"qwertyuiop",
		"qwertzuiopü",
		"azertyuiop",
		"asdfghjklöä",
		"qsdfghjklm",
		//"yxcvbnm", Too short.
		//"zxcvbnm",
		"01234567890",
		"abcdefghijklmnopqrstuvwxyz",
	} {
		if strings.Contains(sequence, strings.ToLower(password)) {
			return PasswordSequential
		}
	}
	return PasswordOK
}
