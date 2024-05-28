// Code generated by "enumer -type=MountMigrationStatus -trimprefix=MigrationStatus -transform=kebab"; DO NOT EDIT.

package vault

import (
	"fmt"
)

const _MountMigrationStatusName = "in-progresssuccessfailure"

var _MountMigrationStatusIndex = [...]uint8{0, 11, 18, 25}

func (i MountMigrationStatus) String() string {
	if i < 0 || i >= MountMigrationStatus(len(_MountMigrationStatusIndex)-1) {
		return fmt.Sprintf("MountMigrationStatus(%d)", i)
	}
	return _MountMigrationStatusName[_MountMigrationStatusIndex[i]:_MountMigrationStatusIndex[i+1]]
}

var _MountMigrationStatusValues = []MountMigrationStatus{0, 1, 2}

var _MountMigrationStatusNameToValueMap = map[string]MountMigrationStatus{
	_MountMigrationStatusName[0:11]:  0,
	_MountMigrationStatusName[11:18]: 1,
	_MountMigrationStatusName[18:25]: 2,
}

// MountMigrationStatusString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func MountMigrationStatusString(s string) (MountMigrationStatus, error) {
	if val, ok := _MountMigrationStatusNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to MountMigrationStatus values", s)
}

// MountMigrationStatusValues returns all values of the enum
func MountMigrationStatusValues() []MountMigrationStatus {
	return _MountMigrationStatusValues
}

// IsAMountMigrationStatus returns "true" if the value is listed in the enum definition. "false" otherwise
func (i MountMigrationStatus) IsAMountMigrationStatus() bool {
	for _, v := range _MountMigrationStatusValues {
		if i == v {
			return true
		}
	}
	return false
}
