package api

const (
	ConstraintDistinctProperty  = "distinct_property"
	ConstraintDistinctHosts     = "distinct_hosts"
	ConstraintRegex             = "regexp"
	ConstraintVersion           = "version"
	ConstraintSemver            = "semver"
	ConstraintSetContains       = "set_contains"
	ConstraintSetContainsAll    = "set_contains_all"
	ConstraintSetContainsAny    = "set_contains_any"
	ConstraintAttributeIsSet    = "is_set"
	ConstraintAttributeIsNotSet = "is_not_set"
)

// Constraint is used to serialize a job placement constraint.
type Constraint struct {
	LTarget string
	RTarget string
	Operand string
}

// NewConstraint generates a new job placement constraint.
func NewConstraint(left, operand, right string) *Constraint {
	return &Constraint{
		LTarget: left,
		RTarget: right,
		Operand: operand,
	}
}
