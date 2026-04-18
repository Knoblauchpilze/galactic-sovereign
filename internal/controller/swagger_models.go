package controller

type ToolkitErrorDoc struct {
	Value   int    `json:"Value"`
	Message string `json:"Message"`
	Cause   any    `json:"Cause"`
}

type PlanetRequestDoc struct {
	Player string `json:"player"`
	Name   string `json:"name"`
}

type PlanetResponseDoc struct {
	Id        string `json:"id"`
	Player    string `json:"player"`
	Name      string `json:"name"`
	Homeworld bool   `json:"homeworld"`
	CreatedAt string `json:"createdAt"`
}

type PlanetResourceResponseDoc struct {
	Planet    string  `json:"planet"`
	Resource  string  `json:"resource"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type PlanetResourceProductionResponseDoc struct {
	Planet     string  `json:"planet"`
	Building   *string `json:"building,omitempty"`
	Resource   string  `json:"resource"`
	Production int     `json:"production"`
}

type PlanetResourceStorageResponseDoc struct {
	Planet   string `json:"planet"`
	Resource string `json:"resource"`
	Storage  int    `json:"storage"`
}

type PlanetBuildingResponseDoc struct {
	Planet    string `json:"planet"`
	Building  string `json:"building"`
	Level     int    `json:"level"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type BuildingActionRequestDoc struct {
	Planet   string `json:"planet"`
	Building string `json:"building"`
}

type BuildingActionResponseDoc struct {
	Id           string `json:"id"`
	Planet       string `json:"planet"`
	Building     string `json:"building"`
	CurrentLevel int    `json:"currentLevel"`
	DesiredLevel int    `json:"desiredLevel"`
	CreatedAt    string `json:"createdAt"`
	CompletedAt  string `json:"completedAt"`
}

type FullPlanetResponseDoc struct {
	PlanetResponseDoc
	Resources       []PlanetResourceResponseDoc           `json:"resources"`
	Productions     []PlanetResourceProductionResponseDoc `json:"productions"`
	Storages        []PlanetResourceStorageResponseDoc    `json:"storages"`
	Buildings       []PlanetBuildingResponseDoc           `json:"buildings"`
	BuildingActions []BuildingActionResponseDoc           `json:"buildingActions"`
}

type PlayerRequestDoc struct {
	ApiUser  string `json:"api_user"`
	Universe string `json:"universe"`
	Name     string `json:"name"`
}

type PlayerResponseDoc struct {
	Id        string `json:"id"`
	ApiUser   string `json:"api_user"`
	Universe  string `json:"universe"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type UniverseRequestDoc struct {
	Name string `json:"name"`
}

type UniverseResponseDoc struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type ResourceResponseDoc struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type BuildingResponseDoc struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type BuildingCostResponseDoc struct {
	Building string  `json:"building"`
	Resource string  `json:"resource"`
	Cost     int     `json:"cost"`
	Progress float64 `json:"progress"`
}

type BuildingResourceProductionResponseDoc struct {
	Building string  `json:"building"`
	Resource string  `json:"resource"`
	Base     int     `json:"base"`
	Progress float64 `json:"progress"`
}

type BuildingResourceStorageResponseDoc struct {
	Building string  `json:"building"`
	Resource string  `json:"resource"`
	Base     int     `json:"base"`
	Scale    float64 `json:"scale"`
	Progress float64 `json:"progress"`
}

type FullBuildingResponseDoc struct {
	BuildingResponseDoc
	Costs       []BuildingCostResponseDoc               `json:"costs"`
	Productions []BuildingResourceProductionResponseDoc `json:"productions"`
	Storages    []BuildingResourceStorageResponseDoc    `json:"storages"`
}

type FullUniverseResponseDoc struct {
	UniverseResponseDoc
	Resources []ResourceResponseDoc     `json:"resources"`
	Buildings []FullBuildingResponseDoc `json:"buildings"`
}
