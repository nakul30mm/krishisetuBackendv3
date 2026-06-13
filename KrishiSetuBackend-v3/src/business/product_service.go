package business

import (
	"errors"
	"krishisetu-backend/models"
	"krishisetu-backend/src/dto"
)

type ProductService struct {
	productRepo ProductRepository
	userRepo    UserRepository
}

func NewProductService(productRepo ProductRepository, userRepo UserRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

func (s *ProductService) CreateProduct(userID uint, req dto.CreateProductDTO) (*models.Product, error) {
	product := models.Product{
		Name:             req.Name,
		Category:         req.Category,
		Description:      req.Description,
		Price:            req.Price,
		Unit:             req.Unit,
		Quantity:         req.Quantity,
		LocationPincode:  req.Location.Pincode,
		LocationDistrict: req.Location.District,
		LocationState:    req.Location.State,
		LocationCity:     req.Location.City,
		LocationArea:     req.Location.Location,
		SellerID:         userID,
		IsAvailable:      true,
		Image1:           req.Image1,
		Image2:           req.Image2,
		Image3:           req.Image3,
	}

	if err := s.productRepo.Create(&product); err != nil {
		return nil, errors.New("Failed to create product")
	}

	return &product, nil
}

func (s *ProductService) FormatProductResponse(p models.Product) map[string]interface{} {
	return map[string]interface{}{
		"id":          p.ID,
		"name":        p.Name,
		"category":    p.Category,
		"description": p.Description,
		"price":       p.Price,
		"unit":        p.Unit,
		"quantity":    p.Quantity,
		"location": map[string]interface{}{
			"pincode":  p.LocationPincode,
			"district": p.LocationDistrict,
			"state":    p.LocationState,
			"city":     p.LocationCity,
			"location": p.LocationArea,
		},
		"image1":      p.Image1,
		"image2":      p.Image2,
		"image3":      p.Image3,
		"seller": map[string]interface{}{
			"id":     p.Seller.ID,
			"name":   p.Seller.FullName,
			"rating": p.Seller.AverageRating,
		},
		"created_at": p.CreatedAt,
	}
}

func (s *ProductService) GetAllProducts(userID uint, search, sort, minPrice, maxPrice string) ([]map[string]interface{}, error) {
	var userPincode string
	if userID != 0 {
		user, err := s.userRepo.FindByID(userID)
		if err == nil && user != nil {
			userPincode = user.Pincode
		}
	}

	products, err := s.productRepo.GetFiltered(search, sort, minPrice, maxPrice, userPincode)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	for _, p := range products {
		response = append(response, s.FormatProductResponse(p))
	}

	return response, nil
}

func (s *ProductService) GetMyProducts(userID uint, search string) ([]map[string]interface{}, error) {
	products, err := s.productRepo.GetBySellerID(userID)
	if err != nil {
		return nil, err
	}

	// Filter by search locally or via DB (GORM handled search, let's keep logic identical)
	var response []map[string]interface{}
	for _, p := range products {
		response = append(response, s.FormatProductResponse(p))
	}

	return response, nil
}

func (s *ProductService) UpdateProduct(userID uint, productID uint, req dto.UpdateProductDTO) (*models.Product, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("Product not found")
	}

	if product.SellerID != userID {
		return nil, errors.New("You can only edit your own products")
	}

	product.Name = req.Name
	product.Category = req.Category
	product.Description = req.Description
	product.Price = req.Price
	product.Unit = req.Unit
	product.Quantity = req.Quantity
	product.LocationPincode = req.Location.Pincode
	product.LocationDistrict = req.Location.District
	product.LocationState = req.Location.State
	product.LocationCity = req.Location.City
	product.LocationArea = req.Location.Location
	product.Image1 = req.Image1
	product.Image2 = req.Image2
	product.Image3 = req.Image3

	if err := s.productRepo.Update(product); err != nil {
		return nil, errors.New("Failed to update product")
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(userID uint, productID uint) error {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return errors.New("Product not found")
	}

	if product.SellerID != userID {
		return errors.New("You can only delete your own products")
	}

	if err := s.productRepo.Delete(productID); err != nil {
		return errors.New("Failed to delete product")
	}

	return nil
}
