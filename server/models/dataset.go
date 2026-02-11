package models

type DatasetCreateReq struct {
	Icon        string `json:"icon" validate:"required,emoji"`
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type DatasetUpdateReq struct {
	ID          uint   `json:"id" validate:"required"`
	Icon        string `json:"icon" validate:"emoji"`
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type DatasetInfo struct {
	ID          uint   `json:"id"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     uint   `json:"owner_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type DatasetListResp struct {
	Total    int64         `json:"total"`
	Datasets []DatasetInfo `json:"datasets"`
}

type DatasetInfoReq struct {
	ID uint `param:"dataset_id" validate:"required"`
}
