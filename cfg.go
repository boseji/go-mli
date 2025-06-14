// cfg.go - Configuration File Handler
//
//     ॐ भूर्भुवः स्वः
//     तत्स॑वि॒तुर्वरे॑ण्यं॒
//    भर्गो॑ दे॒वस्य॑ धीमहि।
//   धियो॒ यो नः॑ प्रचो॒दया॑त्॥
//
//
// बोसजी के द्वारा रचित गो-मिल तन्त्राक्ष्
// ============================
//
// यह गो-क्रमादेश आधारित एम.क्यू.टी.टी अधिलेख में प्रचालेखन का तन्त्राक्ष् है।
//
// एक रचनात्मक भारतीय उत्पाद।
//
// go-mli - Boseji's Golang MQTT Logging command line
//
// Easy to use Golang based MQTT Command line logger.
//
// Sources
// -------
// https://github.com/boseji/go-mli
//
// License
// -------
//
//   go-mli - Boseji's Golang MQTT Logging command line
//   Copyright (C) 2024 by Abhijit Bose (aka. Boseji)
//
//   This program is free software: you can redistribute it and/or modify
//   it under the terms of the GNU General Public License version 2 only
//   as published by the Free Software Foundation.
//
//   This program is distributed in the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without even the implied warranty of
//   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
//
//   You should have received a copy of the GNU General Public License
//   along with this program. If not, see <https://www.gnu.org/licenses/>.
//
//  SPDX-License-Identifier: GPL-2.0-only
//  Full Name: GNU General Public License v2.0 only
//  Please visit <https://spdx.org/licenses/GPL-2.0-only.html> for details.
//

// Configuration File Handler
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// cfg stores Configuration for MQTT and Topics needed for logging.
type cfg struct {
	ADDR           string
	Username       string
	Password       string
	CAFile         string
	ClientID       string
	ClientCertFile string
	ClientKeyFile  string
	Topics         []string
}

// Load helps to read the supplied JSON file and fill up the configuration.
func (m *cfg) Load(Filename string) error {
	bs, err := os.ReadFile(Filename)
	if err != nil {
		return fmt.Errorf("failed to load file %q :\n %v", Filename, err)
	}

	err = json.Unmarshal(bs, m)
	if err != nil {
		return fmt.Errorf("failed to process the file %q :\n %v", Filename, err)
	}

	return nil
}

// Save helps to save back the configuration into the supplied JSON file.
func (m *cfg) Save(Filename string) error {
	bs, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to encode back configuration:\n %v", err)
	}

	err = os.WriteFile(Filename, bs, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %q :\n %v", Filename, err)
	}

	return nil
}

// String implements the Stringer interface to print out the configuration.
func (m cfg) String() string {
	bs, _ := json.MarshalIndent(m, "", "  ")
	return string(bs)
}

// writeTemplate helps to create the default JSON file with
// dummy configuration as a starter.
func writeTemplate(Filename string) error {
	m := &cfg{
		ADDR:           "tcp://192.168.0.0:1883",
		Username:       "Username Here",
		Password:       "Password Here",
		CAFile:         "/path/to/ca.crt-optional",
		ClientID:       "go-mli-demo",
		ClientCertFile: "/path/to/user.client.crt-optional",
		ClientKeyFile:  "/path/to/user.client.key-optional",
		Topics: []string{
			"demo",
			"d1",
			"Sensor1/Temp",
			"Sensor1/Humidity",
		},
	}

	return m.Save(Filename)
}
