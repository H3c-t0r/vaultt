// generated by stringer -type=segmentKind; DO NOT EDIT

package protocol

import "fmt"

const (
	_segmentKind_name_0 = "skInvalidskRequestskReply"
	_segmentKind_name_1 = "skError"
)

var (
	_segmentKind_index_0 = [...]uint8{0, 9, 18, 25}
	_segmentKind_index_1 = [...]uint8{0, 7}
)

func (i segmentKind) String() string {
	switch {
	case 0 <= i && i <= 2:
		return _segmentKind_name_0[_segmentKind_index_0[i]:_segmentKind_index_0[i+1]]
	case i == 5:
		return _segmentKind_name_1
	default:
		return fmt.Sprintf("segmentKind(%d)", i)
	}
}
