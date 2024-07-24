// store.go - Storage Handler
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
// SPDX: GPL-3.0-or-later
//
//   go-mli - Boseji's Golang MQTT Logging command line
//   Copyright (C) 2024 by Abhijit Bose (aka. Boseji)
//
//   This program is free software: you can redistribute it and/or modify
//   it under the terms of the GNU General Public License as published by the
//   Free Software Foundation, either version 3 of the License, or
//   (at your option) any later version.
//
//   This program is distributed in the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without even the implied warranty
//   of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
//   See the GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License along
//   with this program. If not, see <https://www.gnu.org/licenses/>.
//

// Storage Handler
package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
)

const (
	// Wait time for the Storage Goroutine Loop
	STORE_WAIT = 10 * time.Millisecond
	// File permissions for the Log file
	STORE_PERM = 0644
	// Header for Log File
	STORE_HEADER = "Time Stamp,Topic,Data\n"
)

// storeGoroutine is a Go process that waits for a record to be generated
// then it writes the same into the supplied filename.
func storeGoroutine(c <-chan string,
	ctx context.Context, wg *sync.WaitGroup, storeFile string) {
	// Exit with Signalling Completion
	defer wg.Done()
	// Check for Files and Write the Header
	if _, err := os.Stat(storeFile); os.IsNotExist(err) {
		log.Println("[Store] log file does not exists creating one")
		err := os.WriteFile(storeFile,
			[]byte(STORE_HEADER),
			STORE_PERM)
		if err != nil {
			log.Println("[Store] Could not initialize the log file:\n ", err)
			return
		}
	}

	// Process Loop
	for {
		// Channel Receiver
		select {

		case <-ctx.Done():
			log.Println("[Store] Cancel detected")
			return

		case s, ok := <-c:
			if !ok {
				log.Println("[Store] Channel Close detected")
				return
			}
			log.Printf("[Store] Got # %s\n", s)
			f, err := os.OpenFile(storeFile,
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println("[Store] failed to open file:\n ", err)
			}
			defer f.Close()
			if _, err := f.WriteString(s); err != nil {
				log.Println("[Store] failed to write data:\n ", err)
			}

		default:
			time.Sleep(STORE_WAIT)

		}
	}
}
