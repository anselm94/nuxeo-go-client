package internal

import "net/url"

func MergeUrlValues(paramA, paramB url.Values) url.Values {
	if paramA == nil && paramB == nil {
		return nil
	}
	if paramA == nil {
		return paramB
	}
	if paramB == nil {
		return paramA
	}
	merged := url.Values{}
	for k, vs := range paramA {
		for _, v := range vs {
			merged.Add(k, v)
		}
	}
	for k, vs := range paramB {
		for _, v := range vs {
			merged.Add(k, v)
		}
	}
	return merged
}
