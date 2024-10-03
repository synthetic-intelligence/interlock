// INTERLOCK | https://github.com/usbarmory/interlock
// Copyright (c) WithSecure Corporation
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

package interlock

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type jsonObject map[string]interface{}

func parseRequest(r *http.Request) (j jsonObject, err error) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		return
	}

	if conf.Debug {
		log.Printf("%s", body)
	}

	d := json.NewDecoder(strings.NewReader(string(body[:])))
	d.UseNumber()

	err = d.Decode(&j)

	if err != nil {
		return
	}

	return
}

func (j jsonObject) String() (s string) {
	b, err := json.Marshal(j)

	if err != nil {
		log.Print(err)
		return
	}

	s = string(b)

	return
}

func validateRequest(req jsonObject, reqAttrs []string) error {
	for i := 0; i < len(reqAttrs); i++ {
		var ok bool

		args := strings.Split(reqAttrs[i], ":")

		if len(args) != 2 {
			return errors.New("unknown validation argument")
		}

		key := args[0]
		kind := args[1]

		if _, ok = req[key]; !ok {
			return fmt.Errorf("missing attribute %s", key)
		}

		switch kind {
		case "s":
			_, ok = req[key].(string)
		case "b":
			_, ok = req[key].(bool)
		case "n":
			_, ok = req[key].(json.Number)
		case "a":
			_, ok = req[key].([]interface{})
		case "i":
			_, ok = req[key]
		default:
			return errors.New("unknown validation kind")
		}

		if !ok {
			return fmt.Errorf("invalid attribute %s (%s)", key, kind)
		}
	}

	return nil
}
