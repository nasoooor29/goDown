package types

import "fmt"

type SiteError struct {
	*SiteData
	ErrMsg string
}

func (e *SiteError) Error() string {
	return e.ErrMsg
}

func CreateCustomSiteError(sm string) error {
	return &SiteError{
		ErrMsg: sm,
	}
}

func FeatureIsNotSupported(data *SiteData, feature SiteFeature) error {
	return &SiteError{
		SiteData: data,
		ErrMsg:   fmt.Sprintf("feature: %v. is not supported, on site: %v.", feature, data.Url),
	}
}
