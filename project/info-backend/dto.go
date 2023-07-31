package main

type BidDto struct {
	ProductId string `json:"product_id"`
	Value     string `json:"value"`
	BidderId  string `json:"bidder_id"`
	BidId     string `json:"bid_id"`
	TimeStamp string `json:"timestamp"`
}
