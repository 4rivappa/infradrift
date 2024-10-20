package handler

type State struct {
	Version          int        `json:"version"`
	TerraformVersion string     `json:"terraform_version"`
	Resources        []Resource `json:"resources"`
}

type Resource struct {
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Provider  string     `json:"provider"`
	Instances []Instance `json:"instances"`
}

type Instance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

type RequiredIdsResource struct {
	MainFileIds []string
	CommandId   string
}
