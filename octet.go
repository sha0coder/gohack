package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func ip2octet(ip string) [4]int {
	oc := [4]int{0, 0, 0, 0}
	var err error
	for i, o := range strings.Split(ip, ".") {
		oc[i], err = strconv.Atoi(o)
		check(err, "bad ip")
	}
	return oc
}

func octet2ip(oc [4]int) string {
	var ip bytes.Buffer
	for i, o := range oc {
		ip.WriteString(strconv.Itoa(o))
		if i < 3 {
			ip.WriteByte('.')
		}
	}

	return ip.String()
}

func equalOctets(oc1 [4]int, oc2 [4]int) bool {
	return oc1[0] == oc2[0] &&
		oc1[1] == oc2[1] &&
		oc1[2] == oc2[2] &&
		oc1[3] == oc2[3]
}

func incOctet(oc [4]int) [4]int {
	if oc[3] < 255 {
		oc[3]++
		return oc
	}

	if oc[2] < 255 {
		oc[2]++
		oc[3] = 0
		return oc
	}

	if oc[1] < 255 {
		oc[1]++
		oc[2] = 0
		oc[3] = 0
		return oc
	}

	if oc[0] < 255 {
		oc[0]++
		oc[3] = 0
		oc[2] = 0
		oc[1] = 0
		return oc
	}

	fmt.Println("wrong ip range")
	return oc
}
