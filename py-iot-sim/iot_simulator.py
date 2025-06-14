#!/usr/bin/env python3

# iot_mqtt_simulator.py - Python based Simulator Part of the `mli` project
#
# go-mli - Boseji's Golang MQTT Logging command line
#
# IoT Simulator Program
# ----------------------
#
# A script to simulate multiple IoT devices publishing synthetic sensor data to an MQTT broker.
# Supports waveforms like linear, noise, sinusoidal, square, triangle, and sawtooth.
#
# This program is free software; you can redistribute it and/or modify it
# under the terms of the GNU General Public License as published by
# the Free Software Foundation; version 2 of the License only.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA
#
# Dependencies:
# - paho-mqtt
# - numpy
# Install via: pip install paho-mqtt numpy
#
# This program is dependent on `paho-mqtt` and `numpy`
# Please install them before running this program with:
#
# pip install paho-mqtt numpy
#

import argparse
import json
import os
import random
import string
import time
import asyncio
import numpy as np
from typing import Dict, Any
from paho.mqtt import client as mqtt_client

DEFAULT_CONFIG_FILE = "config.json"

def generate_client_id(prefix="client"):
    """
    Generate a unique MQTT client ID.
    :param prefix: Optional prefix to use.
    :return: A unique client ID string.
    """
    return prefix + "_" + ''.join(random.choices(string.ascii_letters + string.digits, k=8))

def default_config() -> Dict[str, Any]:
    """
    Provide default configuration used to create a new config file if one is not specified or found.
    :return: Dictionary representing default settings.
    """
    return {
        "mqtt": {
            "uri": "test.mosquitto.org",
            "port": 1883,
            "username": "",
            "password": "",
            "client_id": "",
            "topic_prefix": "iot/demo"
        },
        "devices": {
            "device1": {
                "Temperature1": {"type": "linear", "start": 20, "slope": 0.1},
                "Voltage1": {"type": "noise", "mean": 3.3, "stddev": 0.05},
                "Sin1": {"type": "sin", "amplitude": 1.0, "frequency": 0.1},
                "Square1": {"type": "square", "amplitude": 1.0, "frequency": 0.2},
                "Saw1": {"type": "sawtooth", "amplitude": 1.0, "frequency": 0.2},
                "Tri1": {"type": "triangle", "amplitude": 1.0, "frequency": 0.2}
            }
        },
        "publish_interval": 2,
        "verbose": True
    }

def load_config(path: str) -> Dict[str, Any]:
    """
    Load configuration from JSON file or create one using defaults if missing.
    :param path: Path to the config JSON file.
    :return: Configuration dictionary.
    """
    if not os.path.exists(path):
        print(f"Config file not found at '{path}'. Creating default.")
        with open(path, "w") as f:
            json.dump(default_config(), f, indent=4)
    with open(path, "r") as f:
        return json.load(f)

def get_value_generator(param_config):
    """
    Return a function to generate parameter values based on specified waveform.
    :param param_config: Dictionary specifying waveform and parameters.
    :return: Callable function producing a value.
    """
    t = 0
    ptype = param_config["type"]
    amp = param_config.get("amplitude", 1.0)
    freq = param_config.get("frequency", 1.0)

    if ptype == "linear":
        start = param_config.get("start", 0)
        slope = param_config.get("slope", 1)
        def linear():
            nonlocal t
            val = start + slope * t
            t += 1
            return round(val, 2)
        return linear

    elif ptype == "noise":
        mean = param_config.get("mean", 0)
        stddev = param_config.get("stddev", 1)
        return lambda: round(random.gauss(mean, stddev), 3)

    elif ptype == "sin":
        def sin():
            nonlocal t
            val = amp * np.sin(2 * np.pi * freq * t)
            t += 1
            return round(val, 3)
        return sin

    elif ptype == "square":
        def square():
            nonlocal t
            val = amp * np.sign(np.sin(2 * np.pi * freq * t))
            t += 1
            return round(val, 3)
        return square

    elif ptype == "sawtooth":
        def saw():
            nonlocal t
            val = amp * (2 * (t * freq - np.floor(t * freq + 0.5)))
            t += 1
            return round(val, 3)
        return saw

    elif ptype == "triangle":
        def tri():
            nonlocal t
            val = amp * (2 * abs(2 * (t * freq - np.floor(t * freq + 0.5))) - 1)
            t += 1
            return round(val, 3)
        return tri

    else:
        raise ValueError(f"Unknown generator type: {ptype}")

async def simulate_device(device_id, parameters, mqtt_cfg, interval, verbose=False):
    """
    Simulate a device sending multiple parameter values to an MQTT broker at regular intervals.
    :param device_id: Unique identifier for the device.
    :param parameters: Dictionary of parameter configurations.
    :param mqtt_cfg: MQTT connection configuration.
    :param interval: Time in seconds between publishes.
    :param verbose: Enable detailed logging if True.
    """
    client_id = mqtt_cfg["client_id"] or generate_client_id(device_id)
    client = mqtt_client.Client(client_id=client_id)

    if mqtt_cfg["username"]:
        client.username_pw_set(mqtt_cfg["username"], mqtt_cfg["password"])

    def on_connect(client, userdata, flags, rc):
        if verbose:
            print(f"[{device_id}] Connected to MQTT broker with result code {rc}")

    client.on_connect = on_connect
    client.connect(mqtt_cfg["uri"], mqtt_cfg["port"])
    client.loop_start()

    generators = {
        param: get_value_generator(cfg) for param, cfg in parameters.items()
    }

    while True:
        for param, generator in generators.items():
            value = generator()
            topic = f"{mqtt_cfg['topic_prefix']}/device/{device_id}/{param}"
            payload = json.dumps({"device": device_id, "param": param, "value": value, "ts": time.time()})
            client.publish(topic, payload)
            if verbose:
                print(f"[{device_id}] Published to {topic}: {payload}")
        await asyncio.sleep(interval)

async def main():
    """
    Entry point: parses config file, loads devices and MQTT settings, and starts simulation tasks.
    """
    parser = argparse.ArgumentParser(description="IoT MQTT Device Simulator")
    parser.add_argument("--config", default=DEFAULT_CONFIG_FILE, help="Path to config JSON file")
    args = parser.parse_args()

    config = load_config(args.config)
    mqtt_cfg = config["mqtt"]
    devices = config["devices"]
    interval = config.get("publish_interval", 2)
    verbose = config.get("verbose", False)

    tasks = []
    for device_id, param_cfgs in devices.items():
        task = simulate_device(device_id, param_cfgs, mqtt_cfg, interval, verbose)
        tasks.append(asyncio.create_task(task))

    await asyncio.gather(*tasks)

if __name__ == "__main__":
    import sys
    import platform

    try:
        if platform.system() == "Windows" and sys.version_info < (3, 8):
            # Windows event loop policy fix for Python < 3.8
            asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())

        loop = asyncio.get_event_loop()
        loop.run_until_complete(main())

    except KeyboardInterrupt:
        print("\nSimulation stopped by user.")
