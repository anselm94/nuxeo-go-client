package internal

import (
	"mime"
	"net/url"
	"strconv"

	"resty.dev/v3"
)

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

func GetStreamContentLengthFrom(res *resty.Response) int {
	contentLength := res.Header().Get(HeaderContentLength)
	if valContentLength, err := strconv.Atoi(contentLength); err == nil {
		return valContentLength
	}
	return -1
}

func GetStreamContentTypeFrom(res *resty.Response) string {
	return res.Header().Get(HeaderContentType)
}

func GetStreamFilenameFrom(res *resty.Response) string {
	// parse content-disposition header and get filename
	contentDisposition := res.Header().Get(HeaderContentDisposition)
	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return ""
	}
	return params["filename"]
}
