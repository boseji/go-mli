// store.go - Storage Handler test
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

// Storage Handler test
package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	TEST_FILE = "test.csv"
)

func Test_storeGoroutine(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T, c chan string)
	}{
		{
			name: "Positive test Sending message",
			fn: func(t *testing.T, c chan string) {
				s := fmt.Sprintf("%q,\"Test\",\"Data\"",
					time.Now().Format(time.RFC3339))
				c <- s
				time.Sleep(100 * time.Millisecond)
				content, err := os.ReadFile(TEST_FILE)
				if err != nil {
					t.Fatal("Failed to read the generated log file")
				}
				if !bytes.Equal(content, []byte(STORE_HEADER+s)) {
					t.Fatalf("Expected: \n%q\n Got:\n %q", STORE_HEADER+s,
						string(content))
				}
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			c := make(chan string, 2)
			wg.Add(1)
			os.Remove(TEST_FILE)
			go storeGoroutine(c, ctx, &wg, TEST_FILE)
			time.Sleep(100 * time.Millisecond)
			tt.fn(t, c)
			time.Sleep(100 * time.Millisecond)
			cancel()
			close(c)
			wg.Wait()
			os.Remove(TEST_FILE)
		})

	}
}

func Test_recordGoroutine(t *testing.T) {
	tests := []struct {
		name     string
		fnBefore func(t *testing.T, c chan string)
		msg      string
		fnAfter  func(t *testing.T, c chan string)
	}{
		{
			name: "Negative test with a filled channel",
			fnBefore: func(t *testing.T, c chan string) {
				c <- "Test 1"
				c <- "Test 2"
			},
			msg:     "Test3",
			fnAfter: func(t *testing.T, c chan string) {},
		},
		{
			name:     "Test the Normal operation",
			fnBefore: func(t *testing.T, c chan string) {},
			msg:      "Test1",
			fnAfter: func(t *testing.T, c chan string) {
				s := <-c
				if s != "Test1" {
					t.Fatalf("Failed to get the correct value: Got = %s Expected = %s",
						s, "Test1")
				}
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			c := make(chan string, 2)
			wg.Add(1)
			tt.fnBefore(t, c)
			go recordGoroutine(c, ctx, &wg, tt.msg, STORE_WAIT)
			time.Sleep(STORE_WAIT * 3)
			cancel()
			tt.fnAfter(t, c)
			close(c)
			wg.Wait()
		})

	}
}

func Test_getRecorder(t *testing.T) {
	tests := []struct {
		name     string
		t        time.Duration
		doRecord func(t *testing.T, rec recorderFn)
		verify   func(t *testing.T, c chan string)
	}{
		{
			name: "Working Record",
			t:    STORE_WAIT * 2,
			doRecord: func(t *testing.T, rec recorderFn) {
				rec("Test1", "Test2")
			},
			verify: func(t *testing.T, c chan string) {
				s := <-c
				if !strings.Contains(s, "\"Test1\",\"Test2\"") {
					t.Fatalf("failed to find sub-string \n expected : %s\n got %s",
						"\"Test1\",\"Test2\"", s)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			c := make(chan string, 2)
			rec := getRecorder(c, ctx, &wg, tt.t)
			tt.doRecord(t, rec)
			time.Sleep(STORE_WAIT * 3)
			cancel()
			tt.verify(t, c)
			close(c)
			wg.Wait()
		})
	}
}
