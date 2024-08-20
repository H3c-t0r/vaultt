// Code generated by "enumer -type=KeyUsage -trimprefix=KeyUsage -transform=snake"; DO NOT EDIT.

package logical

import (
	"fmt"
)

const _KeyUsageName = "encryptdecryptsignverifywrapunwrapgenerate_random"

var _KeyUsageIndex = [...]uint8{0, 7, 14, 18, 24, 28, 34, 49}

func (i KeyUsage) String() string {
	i -= 1
	if i < 0 || i >= KeyUsage(len(_KeyUsageIndex)-1) {
		return fmt.Sprintf("KeyUsage(%d)", i+1)
	}
	return _KeyUsageName[_KeyUsageIndex[i]:_KeyUsageIndex[i+1]]
}

var _KeyUsageValues = []KeyUsage{1, 2, 3, 4, 5, 6, 7}

var _KeyUsageNameToValueMap = map[string]KeyUsage{
	_KeyUsageName[0:7]:   1,
	_KeyUsageName[7:14]:  2,
	_KeyUsageName[14:18]: 3,
	_KeyUsageName[18:24]: 4,
	_KeyUsageName[24:28]: 5,
	_KeyUsageName[28:34]: 6,
	_KeyUsageName[34:49]: 7,
}

// KeyUsageString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func KeyUsageString(s string) (KeyUsage, error) {
	if val, ok := _KeyUsageNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to KeyUsage values", s)
}

// KeyUsageValues returns all values of the enum
func KeyUsageValues() []KeyUsage {
	return _KeyUsageValues
}

// IsAKeyUsage returns "true" if the value is listed in the enum definition. "false" otherwise
func (i KeyUsage) IsAKeyUsage() bool {
	for _, v := range _KeyUsageValues {
		if i == v {
			return true
		}
	}
	return false
}
