package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Store struct {
	updateChannel   chan<- AuctionDto
	auctions        map[string]*AuctionDto
	unprocessedBids map[string]*[]BidDto
	mutex           sync.Mutex
}

func NewStore(updateChannel chan<- AuctionDto) *Store {
	store := &Store{
		updateChannel:   updateChannel,
		auctions:        make(map[string]*AuctionDto),
		unprocessedBids: make(map[string]*[]BidDto),
	}

	go store.monitor()

	return store
}

func (s *Store) TryProcess(bidDto BidDto) error {
	s.mutex.Lock()

	curr, found := s.auctions[bidDto.ProductId]

	if !found {
		bids, found := s.unprocessedBids[bidDto.ProductId]
		if !found {
			newBids := []BidDto{bidDto}
			s.unprocessedBids[bidDto.ProductId] = &newBids
		} else {
			bidDtos := append(*bids, bidDto)
			s.unprocessedBids[bidDto.ProductId] = &bidDtos
		}
		s.mutex.Unlock()
		return nil
	}

	updated, err := updateAuction(*curr, bidDto)
	if err != nil {
		s.mutex.Unlock()
		return err
	}

	s.auctions[bidDto.ProductId] = &updated

	s.updateChannel <- updated

	fmt.Printf("Updated auction: %s with bid %s\n", bidDto.ProductId, bidDto.BidId)

	s.mutex.Unlock()
	return nil
}

func (s *Store) Store(dto ProductDto) error {
	if dto.State != StateSubmitted {
		// Event was sent by this backend type, can be ignored
		return nil
	}

	s.mutex.Lock()

	_, found := s.auctions[dto.ProductId]

	if found {
		s.mutex.Unlock()
		return errors.New("Product already present")
	}

	auction := AuctionDto{Product: dto}

	bids, found := s.unprocessedBids[dto.ProductId]
	if found {
		for _, bid := range *bids {
			auction, _ = updateAuction(auction, bid)
		}
		delete(s.unprocessedBids, dto.ProductId)
	}

	s.auctions[dto.ProductId] = &auction

	s.updateChannel <- auction

	fmt.Printf("Stored Product: %s\n", dto.ProductId)

	s.mutex.Unlock()
	return nil
}

func updateAuction(auctionDto AuctionDto, bidDto BidDto) (AuctionDto, error) {
	if auctionDto.Product.TimeStamp > bidDto.TimeStamp {
		return auctionDto, errors.New("invalid bid, sent before Product creation")
	}

	if auctionDto.Product.State != StateSubmitted && auctionDto.Product.State != StateBidsPlaced {
		return auctionDto, errors.New("invalid bid, Product no longer available")
	}

	if auctionDto.Product.OwnerId == bidDto.BidderId {
		return auctionDto, errors.New("invalid bid, owner can't place bids")
	}

	lastVal := auctionDto.Product.StartingValue
	lastInteraction := auctionDto.Product.TimeStamp
	if auctionDto.Product.State == StateBidsPlaced {
		lastVal = auctionDto.LastBid.Value
		lastInteraction = auctionDto.LastBid.TimeStamp
	}

	bid, _ := strconv.ParseFloat(bidDto.Value, 64)
	val, _ := strconv.ParseFloat(lastVal, 64)

	if lastInteraction+auctionDto.Product.WaitDuration < bidDto.TimeStamp {
		return auctionDto, errors.New("bid not accepted, auction is closed")
	}

	if bid <= val {
		return auctionDto, errors.New("bid not accepted, low value")
	}

	auctionDto.LastBid = bidDto
	auctionDto.Product.State = StateBidsPlaced
	return auctionDto, nil
}

func (s *Store) monitor() {
	for range time.Tick(time.Second * 10) {
		s.closeAuctions()
	}
}

func (s *Store) closeAuctions() {
	s.mutex.Lock()

	for _, auction := range s.auctions {
		if auction.Product.State == StateSubmitted {
			lastInteractionPossible := auction.Product.TimeStamp + auction.Product.WaitDuration
			if lastInteractionPossible < time.Now().Unix() {
				fmt.Printf("Canceled auction for Product %s since no bid have been placed", auction.Product.ProductId)
				auction.Product.State = StateCancelled

				s.updateChannel <- *auction
			}
		} else if auction.Product.State == StateBidsPlaced {
			lastInteractionPossible := auction.LastBid.TimeStamp + auction.Product.WaitDuration
			if lastInteractionPossible < time.Now().Unix() {
				fmt.Printf("Closed auction for Product %s since no bid has been placed in the last %d seconds", auction.Product.ProductId, auction.Product.WaitDuration)
				auction.Product.State = StateAuctionDone

				s.updateChannel <- *auction
			}
		}
	}

	s.mutex.Unlock()
}
