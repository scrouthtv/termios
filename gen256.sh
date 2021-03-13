#!/bin/sh
#curl https://jonasjacek.github.io/colors/data.json | jq '.[] | {name: .name, id: .colorId} | join("=")' \
#	| sed -e 's/^"/Color256/g' -e 's/=/ = Color{Spectrum8, /g' -e 's/"$/, 0, 0}/g'

curl https://jonasjacek.github.io/colors/data.json | jq '.[] | {name: .name, r: .rgb.r, g: .rgb.g, b: .rgb.b} | join(" ")' \
	| awk '{print "\tcoldef256[&Color256" $1 "] = RGB{" $2 ", " $3 ", " $4 "}"}' | tr -d '"'
	#| sed -e 's/^"/Color256/g' -e 's/=/ = Color{Spectrum8, /g' -e 's/"$/, 0, 0}/g'
