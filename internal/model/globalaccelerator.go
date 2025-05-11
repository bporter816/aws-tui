package model

import (
	gaTypes "github.com/aws/aws-sdk-go-v2/service/globalaccelerator/types"
)

type (
	GlobalAcceleratorAccelerator gaTypes.Accelerator
	GlobalAcceleratorListener    gaTypes.Listener
)
