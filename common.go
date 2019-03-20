package ginsta

func getBestDisplayResource(resources []displayResource) string {
	var lastWidth, lastHeight int
	var item string
	for _, resource := range resources {
		if resource.Src == "" {
			continue
		}

		if resource.ConfigWidth > lastWidth ||
			resource.ConfigHeight > lastHeight {
			item = resource.Src
		}
	}

	return item
}

type displayResource struct {
	Src          string `json:"src"`
	ConfigWidth  int    `json:"config_width"`
	ConfigHeight int    `json:"config_height"`
}
