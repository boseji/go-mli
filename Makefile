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
# SPDX: GPL-3.0-or-later
#
#   go-mli - Boseji's Golang MQTT Logging command line
#   Copyright (C) 2024 by Abhijit Bose (aka. Boseji)
#
#   This program is free software: you can redistribute it and/or modify 
#   it under the terms of the GNU General Public License as published by the 
#   Free Software Foundation, either version 3 of the License, or 
#   (at your option) any later version.
#
#   This program is distributed in the hope that it will be useful, 
#   but WITHOUT ANY WARRANTY; without even the implied warranty 
#   of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. 
#   See the GNU General Public License for more details.
#
#   You should have received a copy of the GNU General Public License along 
#   with this program. If not, see <https://www.gnu.org/licenses/>.
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