package main

type InputProductDto struct {
	StartingValue string `json:"starting_value"`
	OwnerId       string `json:"owner_id"`
	Description   string `json:"description"`
}

type OutputProductDto struct {
	ProductId     string `json:"product_id"`
	StartingValue string `json:"starting_value"`
	OwnerId       string `json:"owner_id"`
	Description   string `json:"description"`
	State         string `json:"state"`
}

const StateSubmitted = "submitted"
