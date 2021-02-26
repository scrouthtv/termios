#!/bin/sh
curl https://jonasjacek.github.io/colors/data.json | jq '.[] | {name: .name, id: .colorId} | join("=")' \
	| sed -e 's/^"/Color256/g' -e 's/=/ = Color{Spectrum8, /g' -e 's/"$/, 0, 0}/g'

