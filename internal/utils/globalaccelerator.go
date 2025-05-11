package utils

import (
	"fmt"
	gaTypes "github.com/aws/aws-sdk-go-v2/service/globalaccelerator/types"
	"strconv"
	"strings"
)

func FormatGlobalAcceleratorPortRanges(ports []gaTypes.PortRange) string {
	var parts []string
	for _, v := range ports {
		if v.FromPort == nil && v.ToPort == nil {
			continue
		}
		if v.FromPort == nil {
			parts = append(parts, strconv.Itoa(int(*v.ToPort)))
		} else if v.ToPort == nil {
			parts = append(parts, strconv.Itoa(int(*v.FromPort)))
		} else {
			l, r := *v.FromPort, *v.ToPort
			if l == r {
				parts = append(parts, strconv.Itoa(int(*v.FromPort)))
			} else {
				parts = append(parts, fmt.Sprintf("%v-%v", l, r))
			}
		}
	}
	return strings.Join(parts, ", ")
}
