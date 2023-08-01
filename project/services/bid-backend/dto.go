package main

type InputBidDto struct {
	ProductId string `json:"product_id"`
	Value     string `json:"value"`
	BidderId  string `json:"bidder_id"`
}

type OutputBidDto struct {
	ProductId string `json:"product_id"`
	Value     string `json:"value"`
	BidderId  string `json:"bidder_id"`
	BidId     string `json:"bid_id"`
}
