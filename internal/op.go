// INTERLOCK | https://github.com/usbarmory/interlock
// Copyright (c) WithSecure Corporation
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

// +build linux

package interlock

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"syscall"

	"golang.org/x/term"
)

var opPattern = regexp.MustCompile("^(lock|unlock|derive)(:.+)?$")

func Op(op string) (err error) {
	var cmd string
	var arg string
	var key []byte

	m := opPattern.FindStringSubmatch(op)

	if len(m) != 3 {
		return errors.New("invalid operation")
	}

	cmd = m[1]

	if m[2] != "" {
		if m[2] == ":" {
			return errors.New("invalid operation")
		}

		arg = m[2][1:]
	}

	switch cmd {
	case "unlock":
		if arg == "" {
			return errors.New("invalid operation")
		}

		key, err = promptPassword(false)

		if err != nil {
			return
		}

		err = unlock(arg, string(key))
	case "lock":
		err = lock()
	case "derive":
		var derivedKey string

		if conf.authHSM == nil {
			return errors.New("HSM is required for key derivation")
		}

		if arg == "" {
			key, err = promptPassword(true)

			if err != nil {
				return
			}
		} else {
			key, err = hex.DecodeString(arg)
		}

		if err != nil {
			return
		}

		derivedKey, err = deriveKey(string(key))

		if err != nil {
			return
		}

		fmt.Println(derivedKey)
	}

	return
}

func promptPassword(confirm bool) ([]byte, error) {
	fmt.Print("Password: ")
	password, _ := term.ReadPassword(int(syscall.Stdin))

	if confirm {
		fmt.Printf("\nConfirm password: ")
		confirmation, _ := term.ReadPassword(int(syscall.Stdin))

		if !bytes.Equal(password, confirmation) {
			fmt.Println()
			//lint:ignore ST1005 we want capitalization and punctuation in this error
			return nil, errors.New("Password mismatch!")
		}
	}

	fmt.Println()

	return password, nil
}
