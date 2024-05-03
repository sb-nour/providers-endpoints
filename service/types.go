package service

type ProviderRegions struct {
	Provider string
	Regions  Regions
}

type Regions struct {
	Storage map[string]string `json:"storage"`
	Compute map[string]string `json:"compute"`
}
