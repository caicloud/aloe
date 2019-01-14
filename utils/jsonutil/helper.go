package jsonutil

import (
	"fmt"
)

type MergeOption int

const (
	// CombineOption will make variable with same name be an array
	CombineOption MergeOption = iota
	// OverwriteOption will let variable replace the before one
	OverwriteOption
	// ConflictOption will return error if conflict occured
	ConflictOption
)

func Merge(dst VariableMap, opt MergeOption, isNew bool, objs ...VariableMap) (VariableMap, error) {
	vs := dst
	if isNew {
		vs = dst.Copy()
	}
	for _, obj := range objs {
		if obj == nil {
			continue
		}
		for k, v := range obj.to() {
			dv, ok := vs.Get(k)
			if !ok {
				vs.Set(k, v)
				continue
			}
			switch opt {
			case CombineOption:
				arr, isArr := dv.(VariableArray)
				if !isArr {
					arr = NewVariableArray(k, nil)
					arr.Append(dv)
				}
				arr.Append(v)
				vs.Set(k, arr)
			case OverwriteOption:
				vs.Set(k, v)
			case ConflictOption:
				return nil, fmt.Errorf("variable %v has been defined", k)
			default:
				return nil, fmt.Errorf("unknown merge option, use conflict option, variable %v has been defined", k)
			}
		}
	}
	return vs, nil
}

func IsConflict(dst VariableMap, objs ...VariableMap) bool {
	for _, obj := range objs {
		for k := range obj.to() {
			if _, ok := dst.Get(k); ok {
				return true
			}
		}
	}
	return false
}
