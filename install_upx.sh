#!/usr/bin/env bash

# UPX Installation Script
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

# https://github.com/upx/upx/releases/latest
release_file=$(curl -s https://api.github.com/repos/upx/upx/releases | \
grep browser_download_url | grep amd64_linux | head -n 1 | \
cut -d '"' -f 4)
wget -c "$release_file"

extract=$(echo "$release_file" | cut -d "/" -f 9)
fol=$(tar -tf "./${extract}" | grep upx | head -n 1)
tar -xvf "./${extract}" --wildcards *upx -C .
sudo mv "${fol}upx" /usr/bin
rm "./${extract}"
rm -rf "./${fol}"
