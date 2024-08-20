// main.go - Main Program
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

// Main Program
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
)

const (
	version = "0.0.3"
)

func main() {
	var wg sync.WaitGroup
	var cfg cfg

	// Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("\n go-mli Boseji's Golang MQTT Logging command line")
	fmt.Println("--------------------------------------------------")
	fmt.Println(" Version: " + version)
	fmt.Println()
	defer fmt.Println() // Clearing the Last line

	// Define Flags
	cfgFile := flag.String("config", "config.json",
		"JSON File containing the Configuration.")
	ver := flag.Bool("v", false, "Version number of the program")
	flag.Parse()

	log.Println("[main] Flag Processed: ", flag.Parsed())

	// Print version only
	if *ver {
		log.Println("[main] Program Version :" + version)
		return
	}

	// Get the Config File
	configFile, err := filepath.Abs(*cfgFile)
	if err != nil {
		log.Fatalf("[main][ERROR] Failed to get the Path for %q\n", *cfgFile)
	}

	// Check if the Config File exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Fatalf("[main][ERROR] Configuration file does not exists.")
	}

	// Load Configuration
	err = cfg.Load(configFile)
	if err != nil {
		log.Fatalf("[main][ERROR] Failed to load the Config file %q -\n%v",
			configFile, err)
	}

	// Load Message
	log.Println("[main] `go-mli` Boseji's Golang MQTT Logging command line")
	log.Println("[main] Configuration Loaded -", configFile)
	log.Println("[main] Present Configuration: \n", cfg)

	// Create the Handlers
	loggingFile := time.Now().Format("log-2006-01-02T15-04-05.csv")
	logChan := make(chan string, len(cfg.Topics)*2)
	recFn := getRecorder(logChan, ctx, &wg, STORE_WAIT*2)

	// Handle Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	exitChan := make(chan struct{})
	isError := false

	// Start the Storage Process
	wg.Add(1)
	go storeGoroutine(logChan, ctx, &wg, loggingFile)

	// Ctrl+C Go Routine
	wg.Add(1) // For the Ctrl+C Go Routine
	go func() {
		defer wg.Done()
		select {
		case <-signalChan: // first signal, cancel context
			cancel()
			log.Println()
			log.Println("[INTR] Intercepted a Ctrl+C..")
			break
		case <-ctx.Done():
			log.Println("[INTR][ERROR] Intercepted Error..")
			signal.Stop(signalChan)
			isError = true
			break
		}
		signal.Stop(signalChan)
		close(exitChan) // For Exit
	}()

	// Create MQTT Connection
	mqttOptions := setupMQTT(cfg, cancel, recFn)
	client, err := connectMQTT(mqttOptions)
	if err != nil {
		log.Printf("[main][ERROR] Failed to connect to the MQTT Broker - \n %v", err)
		cancel()
	}

	// Subscribe to the desired topics
	for _, topic := range cfg.Topics {
		err = subscribeMQTT(client, topic)
		if err != nil {
			log.Printf("[main][ERROR] Failed to subscribe to %q\n %v\n",
				topic, err)
			cancel()
		}
	}

	// Wait for Exit with SIGINT or SIGKILL
	<-exitChan
	log.Println("[main] Closing connection..")
	err = disconnectMQTT(client, 20)
	if err != nil {
		log.Printf("[main][ERROR] Failed to close MQTT Connection : %v\n", err)
		cancel()
	}

	// Error in exit
	if isError {
		log.Println("[main][ERROR] Exiting due an Error or failure.")
		os.Exit(1)
	}

	// Wait for Every GoRoutine to Terminate
	wg.Wait()

	// Just to Satisfy Make
	log.Println("[main] Program Terminated Normally.")
	os.Exit(0)
}
