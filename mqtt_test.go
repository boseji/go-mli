// mqtt.go - MQTT Functionality
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

// MQTT Functionality
package main

import (
	"context"
	"os"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// dummyRecorderFn is a mock for the recorderFn in purely logging type function
func dummyRecorderFn(t *testing.T) recorderFn {
	return func(s1, s2 string) {
		t.Logf("Parameters: %q, %q", s1, s2)
	}
}

// addDummyCallbacks configures the MQTT options to have callback functions
// associated with the test cases.
func addDummyCallbacks(t *testing.T, opts *mqtt.ClientOptions) {
	opts.SetDefaultPublishHandler(
		func(client mqtt.Client, msg mqtt.Message) {
			t.Logf("[MQTT] Received message: %q\n from topic: %q\n",
				msg.Payload(), msg.Topic())
		})
	opts.SetOnConnectHandler(
		func(client mqtt.Client) {
			t.Logf("[MQTT] Connected to Broker\n")
		})
	opts.SetConnectionLostHandler(
		func(client mqtt.Client, err error) {
			t.Logf("[MQTT] Connect lost: %v\n", err)
		})
}

func Test_setupMQTT(t *testing.T) {
	type args struct {
		m cfg
	}
	tests := []struct {
		name    string
		args    args
		checkFn func(t *testing.T, opts *mqtt.ClientOptions)
	}{
		{
			name: "Basic Config",
			args: args{
				m: cfg{
					ADDR:     ":1883",
					ClientID: "go-mli-mqtt-test",
				},
			},
			checkFn: func(t *testing.T, opts *mqtt.ClientOptions) {
				if opts.ClientID != "go-mli-mqtt-test" {
					t.Errorf("failed to get correct client ID : %q",
						opts.ClientID)
				}
				if opts.Servers[0].String() != "tcp://127.0.0.1:1883" {
					t.Errorf("failed to get correct Server : %q",
						opts.Servers[0].String())
				}
				if len(opts.Username) != 0 {
					t.Errorf("failed to get correct username: %q",
						opts.Username)
				}
			},
		},
		{
			name: "Config with Passwords",
			args: args{
				m: cfg{
					ADDR:     ":1883",
					ClientID: "go-mli-mqtt-test",
					Username: "test-user",
					Password: "test-Password",
				},
			},
			checkFn: func(t *testing.T, opts *mqtt.ClientOptions) {
				if opts.Username != "test-user" {
					t.Errorf("failed to get correct username: %q",
						opts.Username)
				}
				if opts.Password != "test-Password" {
					t.Errorf("failed to get correct password: %q",
						opts.Password)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, cancel := context.WithCancel(context.Background())
			rec := dummyRecorderFn(t)
			got := setupMQTT(tt.args.m, cancel, rec)
			tt.checkFn(t, got)
		})
	}
}

func Test_connectMQTT(t *testing.T) {
	type args struct {
		optsFn func(t *testing.T) *mqtt.ClientOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		cleanFs []string
	}{
		{
			name: "Basic connection to test.mosquitto.org",
			args: args{
				optsFn: func(t *testing.T) *mqtt.ClientOptions {
					opts := mqtt.NewClientOptions()
					opts.AddBroker("mqtt://test.mosquitto.org:1883")
					opts.SetClientID("go-mli-testing-connectMQTT")
					addDummyCallbacks(t, opts)
					return opts
				},
			},
		},
		{
			name: "Basic Test using setupMQTT",
			args: args{
				optsFn: func(t *testing.T) *mqtt.ClientOptions {
					_, cancel := context.WithCancel(context.Background())
					rec := dummyRecorderFn(t)
					got := setupMQTT(cfg{
						ADDR:     "mqtt://test.mosquitto.org:1883",
						ClientID: "go-mli-testing-connectMQTT",
					}, cancel, rec)
					return got
				},
			},
		},
		{
			name: "Auth Test MQTT, unencrypted, authenticated",
			args: args{
				optsFn: func(t *testing.T) *mqtt.ClientOptions {
					_, cancel := context.WithCancel(context.Background())
					rec := dummyRecorderFn(t)
					got := setupMQTT(cfg{
						ADDR:     "mqtt://test.mosquitto.org:1884",
						ClientID: "go-mli-testing-connectMQTT",
						Username: "rw",
						Password: "readwrite",
					}, cancel, rec)
					return got
				},
			},
		},
		{
			name: "Negative Test - wrong URI",
			args: args{
				optsFn: func(t *testing.T) *mqtt.ClientOptions {
					_, cancel := context.WithCancel(context.Background())
					rec := dummyRecorderFn(t)
					got := setupMQTT(cfg{
						ADDR:     "mqtt://test.mosquittoi.org:1884",
						ClientID: "go-mli-testing-connectMQTT",
						Username: "rw",
						Password: "readwrite",
					}, cancel, rec)
					return got
				},
			},
			wantErr: true,
		},
		{
			name: "Negative Test - wrong username",
			args: args{
				optsFn: func(t *testing.T) *mqtt.ClientOptions {
					_, cancel := context.WithCancel(context.Background())
					rec := dummyRecorderFn(t)
					got := setupMQTT(cfg{
						ADDR:     "mqtt://test.mosquitto.org:1884",
						ClientID: "go-mli-testing-connectMQTT",
						Username: "rw1",
						Password: "readwrite",
					}, cancel, rec)
					return got
				},
			},
			wantErr: true,
		},
		// {
		// 	name: "Secure Test MQTT, encrypted, unauthenticated",
		// 	args: args{
		// 		optsFn: func(t *testing.T) *mqtt.ClientOptions {
		// 			_, cancel := context.WithCancel(context.Background())
		// 			rec := dummyRecorderFn(t)

		// 			got := setupMQTT(cfg{
		// 				ADDR:     "tcps://test.mosquitto.org:8883",
		// 				ClientID: "go-mli-testing-connectMQTT",
		// 				//CAFile:   "./others/mosquitto.org.crt",
		// 			}, cancel, rec)
		// 			return got
		// 		},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := connectMQTT(tt.args.optsFn(t))
			if (err != nil) != tt.wantErr {
				t.Errorf("connectMQTT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			// Safe Disconnect
			err = disconnectMQTT(client, 10)
			if err != nil {
				t.Errorf("disconnection failed: %v", err)
			}
			// Clean files
			for _, fl := range tt.cleanFs {
				os.Remove(fl)
			}
		})
	}
}
