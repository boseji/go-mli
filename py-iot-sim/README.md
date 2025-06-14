>
> ॐ भूर्भुवः स्वः
> 
> तत्स॑वि॒तुर्वरे॑ण्यं॒
> 
> भर्गो॑ दे॒वस्य॑ धीमहि।
> 
> धियो॒ यो नः॑ प्रचो॒दया॑त्॥
> 
# बोसजी के द्वारा रचित गो-मिल तन्त्राक्ष्।

> एक एम.क्यू.टी.टी अधिलेख को प्रचालेखन करने वाला तन्त्राक्ष्।

यह गो-क्रमादेश आधारित एम.क्यू.टी.टी अधिलेख में प्रचालेखन का तन्त्राक्ष् है।
यह किसी विशेष एम.क्यू.टी.टी संग्राहक से कई विषयों को प्रचालेखन करने की उपियोगिता देता है।

रचित प्रचालेखन विधि `csv` प्रारूप के समान है।

***एक रचनात्मक भारतीय उत्पाद।***

## `go-mli` Boseji's Golang MQTT Logging command line

Easy to use Golang based MQTT Command line logger.
This allows to log multiple topics from a particular MQTT Broker.

The generated log is akin to `csv` format.

---

## Python IoT MQTT Simulator - Part of the `go-mli` Project

This Python tool simulates multiple IoT devices publishing synthetic sensor data to an MQTT broker. It supports configurable data generation patterns such as linear, noise, sinusoidal, square, triangle, and sawtooth waveforms. The simulation runs concurrently across devices and parameters.

---

## Features

* MQTT client simulator for IoT use-cases
* Supports multiple virtual devices
* Each device can emit multiple parameter streams
* Parameter generators include:

  * Linear
  * Noise (Gaussian)
  * Sinusoidal
  * Square wave
  * Triangle wave
  * Sawtooth wave
* Fully configurable via JSON
* Parallel publishing using `asyncio`
* Verbose logging option

---

## Getting Started

### Installation

```bash
pip install paho-mqtt numpy
```

### Running the Program

```bash
python iot_mqtt_simulator.py --config config.json
```

If `config.json` is not specified or does not exist, it will be generated with default values.

---

## Configuration File (`config.json`)

```json
{
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
  "verbose": true
}
```

### Parameter Types

* **linear**: gradually increasing or decreasing value

  ```json
  {"type": "linear", "start": 20, "slope": 0.1}
  ```
* **noise**: random value around a mean with standard deviation

  ```json
  {"type": "noise", "mean": 5.0, "stddev": 0.1}
  ```
* **sin**: sinusoidal value
* **square**: alternating +amplitude and -amplitude
* **triangle**: triangle waveform
* **sawtooth**: sawtooth waveform

---

## MQTT Message Format

Each message published to MQTT looks like:

```json
{
  "device": "device1",
  "param": "Temperature1",
  "value": 23.4,
  "ts": 1718033043.2251
}
```

### MQTT Topic Format

```
<iot-topic-prefix>/device/<device-id>/<param>
```

Example:

```
iot/demo/device/device1/Temperature1
```

---

## License

This project is released under the GNU General Public License v2. See the [LICENSE](../LICENSE.txt) file for details.

Sources: <https://github.com/boseji/go-mli>

`go-mli` - Boseji's Golang MQTT Logging command line.

Copyright (C) 2024 by Abhijit Bose (aka. Boseji)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License version 2 only
as published by the Free Software Foundation.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

SPDX-License-Identifier: `GPL-2.0-only`

Full Name: `GNU General Public License v2.0 only`

Please visit <https://spdx.org/licenses/GPL-2.0-only.html> for details.
