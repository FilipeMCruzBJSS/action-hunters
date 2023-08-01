package main

type BidDto struct {
	ProductId string `json:"product_id"`
	Value     string `json:"value"`
	BidderId  string `json:"bidder_id"`
	BidId     string `json:"bid_id"`
	TimeStamp string `json:"timestamp"`
}

type ProductDto struct {
	ProductId     string `json:"product_id"`
	StartingValue string `json:"starting_value"`
	OwnerId       string `json:"owner_id"`
	Description   string `json:"description"`
	State         string `json:"state"`
	TimeStamp     string `json:"timestamp"`
}

const (
	StateSubmitted   = "submitted"
	StateBidsPlaced  = "bids_placed"
	StateCancelled   = "canceled"
	StateAuctionDone = "auction_done"
)

type AuctionDto struct {
	product ProductDto
	lastBid BidDto
}
