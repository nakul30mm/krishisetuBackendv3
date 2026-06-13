package business

import (
	"errors"
	"krishisetu-backend/models"
	"krishisetu-backend/src/dto"
)

type MarketplaceRequestService struct {
	requestRepo MarketplaceRequestRepository
	productRepo ProductRepository
}

func NewMarketplaceRequestService(
	requestRepo MarketplaceRequestRepository,
	productRepo ProductRepository,
) *MarketplaceRequestService {
	return &MarketplaceRequestService{
		requestRepo: requestRepo,
		productRepo: productRepo,
	}
}

func (s *MarketplaceRequestService) CreateMarketplaceRequest(buyerID uint, req dto.CreateMarketplaceRequestDTO) (*models.MarketplaceRequest, error) {
	product, err := s.productRepo.FindByID(req.ProductID)
	if err != nil {
		return nil, errors.New("Product not found")
	}

	if product.SellerID == buyerID {
		return nil, errors.New("You cannot buy your own product")
	}

	marketplaceRequest := models.MarketplaceRequest{
		ProductID:      req.ProductID,
		BuyerID:        buyerID,
		SellerID:       product.SellerID,
		RequestedPrice: req.RequestedPrice,
		Quantity:       req.Quantity,
		Status:         "PENDING",
	}

	if err := s.requestRepo.Create(&marketplaceRequest); err != nil {
		return nil, errors.New("Failed to create buy request")
	}

	return &marketplaceRequest, nil
}

func (s *MarketplaceRequestService) GetSentMarketplaceRequests(buyerID uint, search string) ([]models.MarketplaceRequest, error) {
	return s.requestRepo.GetByBuyerID(buyerID, search)
}

func (s *MarketplaceRequestService) GetReceivedMarketplaceRequests(sellerID uint, search string) ([]models.MarketplaceRequest, error) {
	return s.requestRepo.GetBySellerID(sellerID, search)
}

func (s *MarketplaceRequestService) UpdateMarketplaceRequestStatus(sellerID uint, reqID uint, action string) (*models.MarketplaceRequest, error) {
	var status string
	// if action == "approve" {
	// 	status = "APPROVED"
	// } else if action == "reject" {
	// 	status = "REJECTED"
	// } else {
	// 	return nil, errors.New("Invalid action")
	// }

	switch action {
	case "approve":
		status = "APPROVED"
	case "rejected":
		status = "REJECTED"
	default:
		return nil, errors.New("Invalid action")
	}

	req, err := s.requestRepo.FindByID(reqID)
	if err != nil {
		return nil, errors.New("Request not found")
	}

	if req.SellerID != sellerID {
		return nil, errors.New("Only the product seller can map this action")
	}

	if req.Status != "PENDING" {
		return nil, errors.New("Request is already " + req.Status)
	}

	req.Status = status
	if err := s.requestRepo.Update(req); err != nil {
		return nil, err
	}

	return req, nil
}

func (s *MarketplaceRequestService) DeleteMarketplaceRequest(buyerID uint, reqID uint) error {
	req, err := s.requestRepo.FindByID(reqID)
	if err != nil {
		return errors.New("Request not found")
	}

	if req.BuyerID != buyerID {
		return errors.New("Only the buyer can delete this request")
	}

	if req.Status != "PENDING" {
		return errors.New("Only pending requests can be deleted")
	}

	return s.requestRepo.Delete(reqID)
}

func (s *MarketplaceRequestService) ConfirmMarketplaceTransaction(buyerID uint, reqID uint, req dto.ConfirmTransactionDTO) (*models.MarketplaceRequest, error) {
	marketplaceRequest, err := s.requestRepo.FindByID(reqID)
	if err != nil {
		return nil, errors.New("Request not found")
	}

	if marketplaceRequest.BuyerID != buyerID {
		return nil, errors.New("Only the buyer can confirm this transaction")
	}

	if marketplaceRequest.Status != "APPROVED" {
		return nil, errors.New("Transaction status can only be confirmed for approved requests")
	}

	marketplaceRequest.TransactionStatus = req.Status
	if err := s.requestRepo.Update(marketplaceRequest); err != nil {
		return nil, err
	}

	return marketplaceRequest, nil
}

func (s *MarketplaceRequestService) UpdateMarketplaceRequest(buyerID uint, reqID uint, req dto.UpdateMarketplaceRequestDTO) (*models.MarketplaceRequest, error) {
	marketplaceRequest, err := s.requestRepo.FindByID(reqID)
	if err != nil {
		return nil, errors.New("Request not found")
	}

	if marketplaceRequest.BuyerID != buyerID {
		return nil, errors.New("Only the buyer can edit this request")
	}

	if marketplaceRequest.Status != "PENDING" {
		return nil, errors.New("Only pending requests can be edited")
	}

	marketplaceRequest.Quantity = req.Quantity
	if err := s.requestRepo.Update(marketplaceRequest); err != nil {
		return nil, errors.New("Failed to update request")
	}

	return marketplaceRequest, nil
}
