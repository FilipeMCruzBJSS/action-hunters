package main

type BidDto struct {
	ProductId string `json:"product_id"`
	Value     string `json:"value"`
	BidderId  string `json:"bidder_id"`
	BidId     string `json:"bid_id"`
	TimeStamp int64  `json:"timeStamp"`
}

type ProductDto struct {
	ProductId     string `json:"product_id"`
	StartingValue string `json:"starting_value"`
	OwnerId       string `json:"owner_id"`
	Description   string `json:"description"`
	State         string `json:"state"`
	TimeStamp     int64  `json:"timeStamp"`
	WaitDuration  int64  `json:"wait_duration"`
}

const (
	StateSubmitted   = "submitted"
	StateBidsPlaced  = "bids_placed"
	StateCancelled   = "canceled"
	StateAuctionDone = "auction_done"
)

type AuctionDto struct {
	Product ProductDto `json:"product"`
	LastBid BidDto     `json:"last_id"`
}
