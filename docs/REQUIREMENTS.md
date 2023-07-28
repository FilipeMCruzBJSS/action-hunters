# Requirements

This document introduces the requirements of this project.
It is divided by sprints.

## First Sprint

Due date: 03/08/23

### Functional Requirements

Two major roles were defined: bidder and seller.

These are the basic use cases for each role:

- As a bidder I want to bid on a specific product;
- As a bidder I want to see the available products and their current bid;
- As a seller I want to publish a new product to be auctioned.

### Non-Functional Requirements

- Must use go;
- Must have a microservice architecture;
- Must use kafka;
- Must be possible to replicate each microservice to improve the platform performance and deal with usage spikes.

## Second Sprint

Due date: TBD

### Functional Requirements

A new role is defined, the system.

- As the system I should ensure the bidder has the liquidity necessary to keep bids up;
- As the system I should, from time to time collect metrics about the business, e.g.:
  - Amount of money that flows thought the platform;
  - Products with the highest bids;
  - Sellers that got payed the most for their products.
- As a seller i want to see how are the bids on my products.
