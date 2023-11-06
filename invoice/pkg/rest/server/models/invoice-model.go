package models

type Invoice struct {
	Id int64 `json:"id,omitempty"`

	Amount float32 `json:"amount,omitempty"`
}
