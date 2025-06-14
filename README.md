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

## Install `go-mli`

```sh
go install github.com/boseji/go-mli@latest
```

## कार्यविधि - Usage

### `upx` क्रमादेश

`UPX` - (नवीनतम संस्करण) संक्षिप्त करने वाला क्रमादेश।

```sh
# https://github.com/upx/upx/releases/latest
export release_file=$(curl -s https://api.github.com/repos/upx/upx/releases | \
grep browser_download_url | grep amd64_linux | head -n 1 | \
cut -d '"' -f 4)
wget -c "$release_file"

export extract=$(echo $release_file | cut -d "/" -f 9)
export fol=$(tar -tf "./${extract}" | grep upx | head -n 1)
tar -xvf "./${extract}" --wildcards *upx -C .
sudo mv "${fol}upx" /usr/bin
rm "./${extract}"
rm -rf "./${fol}"
```

### निर्माण क्रमादेश - Build Executables Command

```sh
make build
```

### परीक्षा क्रमादेश - Test the Code

```sh
make test
```

## Name `go-mli`

`go-` prefix is the designate this as a Golang project.

The `mli => MLI` expands as `Mqtt Logging command lIne`.

*Naming convention* of this application as per the comment [here](https://www.reddit.com/r/golang/comments/r3as15/comment/hma99nc/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button) and a [Go playbook here](https://go.dev/play/p/MNfRtvAn0Po).

## License

This project is released under the GNU General Public License v2. See the [LICENSE](./LICENSE.txt) file for details.

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

