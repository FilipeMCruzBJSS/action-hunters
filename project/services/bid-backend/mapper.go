package main

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
)

func Verify(dto InputBidDto) (OutputBidDto, error) {

	if dto.ProductId == "" {
		return OutputBidDto{}, errors.New("missing product id property")
	}

	if dto.BidderId == "" {
		return OutputBidDto{}, errors.New("missing bidder id property")
	}

	if dto.Value == "" {
		return OutputBidDto{}, errors.New("missing value property")
	}

	_, err := strconv.ParseFloat(dto.Value, 64)
	if err != nil {
		return OutputBidDto{}, errors.New("invalid value")
	}

	return OutputBidDto{
		ProductId: dto.ProductId,
		Value:     dto.Value,
		BidderId:  dto.BidderId,
		BidId:     uuid.New().String(),
	}, nil
}
