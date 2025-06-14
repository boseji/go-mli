# Makefile
#
# go-mli - Boseji's Golang MQTT Logging command line
#
# Sources
# -------
# https://github.com/boseji/go-mli
#
# License
# -------
#
#   go-mli - Boseji's Golang MQTT Logging command line
#   Copyright (C) 2024 by Abhijit Bose (aka. Boseji)
#
#   This program is free software: you can redistribute it and/or modify
#   it under the terms of the GNU General Public License version 2 only
#   as published by the Free Software Foundation.
#
#   This program is distributed in the hope that it will be useful,
#   but WITHOUT ANY WARRANTY; without even the implied warranty of
#   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. 
#
#   You should have received a copy of the GNU General Public License
#   along with this program. If not, see <https://www.gnu.org/licenses/>.
#
#  SPDX-License-Identifier: GPL-2.0-only
#  Full Name: GNU General Public License v2.0 only
#  Please visit <https://spdx.org/licenses/GPL-2.0-only.html> for details.
#

GOFILES  := cfg.go store.go mqtt.go main.go

run:
	go mod tidy
	go run ${GOFILES}

build: test
	go mod tidy
# '-' Helps to Ignore Errors (https://stackoverflow.com/a/2670143)
	-mkdir build 2> /dev/null
	GOOS="windows/amd64" & go build -o build/mli_windows_x64.exe ${GOFILES}
	upx --best build/mli_windows_x64.exe
	GOOS="linux/amd64" & go build -o build/mli_linux_x64 ${GOFILES}
	upx --best build/mli_linux_x64

test:
	go mod tidy
# '-' Helps to Ignore Errors (https://stackoverflow.com/a/2670143)
	-mkdir build 2> /dev/null
	go test -v -race

clean:
	rm -rf build