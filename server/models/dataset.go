package models

type DatasetCreateReq struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type DatasetUpdateReq struct {
	Name        string `json:"name" validate:"omitempty,min=1,max=100"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

type DatasetResp struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     uint   `json:"owner_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type DatasetListResp struct {
	Datasets []DatasetResp `json:"datasets"`
}

type DatasetCreateResp struct {
	ID uint `json:"id"`
}

type DatasetUpdateResp struct {
	ID uint `json:"id"`
}

type DatasetDeleteResp struct {
	ID uint `json:"id"`
}
