package jsonutil

import (
	"fmt"
)

// MergeOption defines option of merge function
type MergeOption int

const (
	// CombineOption will make variable with same name be an array
	CombineOption MergeOption = iota
	// OverwriteOption will let variable replace the before one
	OverwriteOption
	// DeepOverwriteOption will try to combine array and map
	// both array: deep overwrite element one by one
	// both map: deep overwrite element one by one
	// others: just replace the before one
	DeepOverwriteOption
	// ConflictOption will return error if conflict occurred
	ConflictOption
)

// VarFilter defines a function to check whether var will be filtered
type VarFilter func(key string) bool

// MergeWithFilter merges filtered variables of VariableMaps objs into VariableMap dst by differernt options
func MergeWithFilter(dst VariableMap, opt MergeOption, isNew bool, filter VarFilter, objs ...VariableMap) (VariableMap, error) {
	vs := dst
	if dst == nil {
		vs = NewVariableMap("", nil)
	} else if isNew {
		vs = dst.Copy()
	}
	for _, obj := range objs {
		if obj == nil {
			continue
		}
		for k, v := range obj.to() {
			if filter != nil {
				if filter(k) {
					continue
				}
			}
			dv, ok := vs.Get(k)
			switch opt {
			case CombineOption:
				var array VariableArray
				if !ok {
					array = NewVariableArray(k, nil)
				} else {
					arr, isArr := dv.(VariableArray)
					if !isArr {
						arr = NewVariableArray(k, nil)
						arr.Append(dv)
					}
					array = arr
				}
				array.Append(v)
				vs.Set(k, array)
			case OverwriteOption:
				vs.Set(k, v)
			case DeepOverwriteOption:
				vs.Set(k, deepOverwrite(dv, v))
			case ConflictOption:
				if ok {
					return nil, fmt.Errorf("variable %v has been defined, variables: %v", k, vs)
				}
				vs.Set(k, v)
			default:
				return nil, fmt.Errorf("unknown merge option, use conflict option, variable %v has been defined", k)
			}
		}
	}
	return vs, nil
}

// Merge merges VariableMaps objs into VariableMap dst by differernt options
func Merge(dst VariableMap, opt MergeOption, isNew bool, objs ...VariableMap) (VariableMap, error) {
	return MergeWithFilter(dst, opt, isNew, nil, objs...)
}

// IsConflict check whether objs are conflict with dst
func IsConflict(dst VariableMap, objs ...VariableMap) bool {
	for _, obj := range objs {
		if obj == nil {
			continue
		}
		for k := range obj.to() {
			if _, ok := dst.Get(k); ok {
				return true
			}
		}
	}
	return false
}

func deepOverwrite(dst, src Variable) Variable {
	if dst == nil || dst.Type() == NullType {
		return src
	}
	if src == nil || src.Type() == NullType {
		return dst
	}
	dstvm, dstok := dst.(VariableMap)
	srcvm, srcok := src.(VariableMap)
	if dstok && srcok {
		return deepOverwriteMap(dstvm, srcvm)
	}
	dstva, dstok := dst.(VariableArray)
	srcva, srcok := src.(VariableArray)
	if dstok && srcok {
		return deepOverwriteArray(dstva, srcva)
	}
	return src
}

func deepOverwriteMap(dst, src VariableMap) VariableMap {
	for k, v := range src.to() {
		dv, _ := dst.Get(k)
		dst.Set(k, deepOverwrite(dv, v))
	}
	return dst
}

func deepOverwriteArray(dst, src VariableArray) VariableArray {
	for i, v := range src.to() {
		if i >= dst.Len() {
			dst.Append(v)
		} else {
			dv := dst.Get(i)
			dst.Set(i, deepOverwrite(dv, v))
		}
	}
	return dst
}
