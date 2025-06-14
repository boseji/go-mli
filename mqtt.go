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

// MQTT Functionality
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// setupMQTT setups up the options for the MQTT connection from the
// supplied configuration and returns the same.
func setupMQTT(m cfg,
	cancel context.CancelFunc,
	rec recorderFn) *mqtt.ClientOptions {
	// From (https://www.emqx.com/en/blog/how-to-use-mqtt-in-golang)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(m.ADDR)
	opts.SetClientID(m.ClientID)
	// If Username is available
	if len(m.Username) > 0 {
		opts.SetUsername(m.Username)
		opts.SetPassword(m.Password)
	}
	// If CA files are available
	if len(m.CAFile) > 0 {
		certPool := x509.NewCertPool()
		ca, err := os.ReadFile(m.CAFile)
		if err != nil {
			log.Fatalln(err.Error())
		}
		certPool.AppendCertsFromPEM(ca)
		opts.SetTLSConfig(&tls.Config{
			RootCAs: certPool,
		})
	}

	// Set Callbacks
	opts.SetDefaultPublishHandler(
		func(client mqtt.Client, msg mqtt.Message) {
			log.Printf("[MQTT] Received message: %q from topic: %q\n",
				msg.Payload(), msg.Topic())
			// Send for Record
			rec(msg.Topic(), string(msg.Payload()))
		})
	opts.SetOnConnectHandler(
		func(client mqtt.Client) {
			log.Println("[MQTT] Connected to Broker")
		})
	opts.SetConnectionLostHandler(
		func(client mqtt.Client, err error) {
			log.Printf("[MQTT] Connect lost: %v", err)
			cancel()
		})

	// Clean Sessions for each run
	opts.SetCleanSession(true)

	return opts
}

// connectMQTT creates the MQTT Client using the supplied MQTT options
// and returns the same upon connection.
func connectMQTT(opts *mqtt.ClientOptions) (mqtt.Client, error) {
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to mqtt:\n %v", token.Error())
	}
	return client, nil
}

// disconnectMQTT helps to disconnect the client within a given time period
// supplied in number of milliseconds.
func disconnectMQTT(client mqtt.Client, ms uint) error {
	if client == nil {
		return fmt.Errorf("no client")
	}
	if client.IsConnected() {
		client.Disconnect(ms)
		log.Println("[MQTT] Disconnected.")
		return nil
	}
	return fmt.Errorf("mqtt not connected")
}

// subscribeMQTT helps to create subscription to the supplied topics
// for the client.
func subscribeMQTT(client mqtt.Client, topic string) error {
	if client == nil {
		return fmt.Errorf("no client")
	}
	token := client.Subscribe(topic, 1, nil)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to %q:\n %v",
			topic, token.Error())
	}
	log.Printf("[MQTT] Subscribed to topic: %q\n", topic)
	return nil
}
