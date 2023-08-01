package main

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
)

func Verify(dto InputProductDto) (OutputProductDto, error) {

	if dto.Description == "" {
		return OutputProductDto{}, errors.New("missing description property")
	}

	if dto.OwnerId == "" {
		return OutputProductDto{}, errors.New("missing owner id property")
	}

	if dto.StartingValue == "" {
		return OutputProductDto{}, errors.New("missing starting value property")
	}

	_, err := strconv.ParseFloat(dto.StartingValue, 64)
	if err != nil {
		return OutputProductDto{}, errors.New("invalid starting value")
	}

	return OutputProductDto{
		ProductId:     uuid.New().String(),
		StartingValue: dto.StartingValue,
		OwnerId:       dto.OwnerId,
		Description:   dto.Description,
		State:         StateSubmitted,
	}, nil
}
