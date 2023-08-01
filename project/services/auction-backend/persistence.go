package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type Store struct {
	auctions        map[string]*AuctionDto
	unprocessedBids map[string]*[]BidDto
	mutex           sync.Mutex
}

func NewStore() *Store {
	return &Store{
		auctions:        make(map[string]*AuctionDto),
		unprocessedBids: make(map[string]*[]BidDto),
	}
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

	fmt.Printf("Updated auction: %s with bid %s\n", bidDto.ProductId, bidDto.BidId)

	s.mutex.Unlock()
	return nil
}

func (s *Store) Store(dto ProductDto) error {
	if dto.State != StateSubmitted {
		return errors.New("invalid state")
	}

	s.mutex.Lock()

	_, found := s.auctions[dto.ProductId]

	if found {
		s.mutex.Unlock()
		return errors.New("product already present")
	}

	auction := AuctionDto{product: dto}

	bids, found := s.unprocessedBids[dto.ProductId]
	if found {
		for _, bid := range *bids {
			auction, _ = updateAuction(auction, bid)
		}
		delete(s.unprocessedBids, dto.ProductId)
	}

	s.auctions[dto.ProductId] = &auction

	fmt.Printf("Stored product: %s\n", dto.ProductId)

	s.mutex.Unlock()
	return nil
}

func (s *Store) Get(productId string) (AuctionDto, error) {
	s.mutex.Lock()

	curr, found := s.auctions[productId]

	if !found {
		s.mutex.Unlock()
		return AuctionDto{}, errors.New("product not present")
	}

	fmt.Printf("Fetched product: %s\n", productId)

	s.mutex.Unlock()
	return *curr, nil
}

func updateAuction(auctionDto AuctionDto, bidDto BidDto) (AuctionDto, error) {
	if auctionDto.product.TimeStamp > bidDto.TimeStamp {
		return auctionDto, errors.New("invalid bid, sent before product creation")
	}

	if auctionDto.product.State != StateSubmitted && auctionDto.product.State != StateBidsPlaced {
		return auctionDto, errors.New("invalid bid, product no longer available")
	}

	if auctionDto.product.OwnerId == bidDto.BidderId {
		return auctionDto, errors.New("invalid bid, owner can't place bids")
	}

	lastVal := auctionDto.product.StartingValue
	if auctionDto.product.State == StateBidsPlaced {
		lastVal = auctionDto.lastBid.Value
	}

	bid, _ := strconv.ParseFloat(bidDto.Value, 64)
	val, _ := strconv.ParseFloat(lastVal, 64)

	if bid <= val {
		return auctionDto, errors.New("bid not accepted, low value")
	}

	auctionDto.lastBid = bidDto
	auctionDto.product.State = StateBidsPlaced
	return auctionDto, nil
}
