package service

type promotionCommissionPlanService struct {
	*BaseService
}

func NewPromotionCommissionPlanService(params ...any) *promotionCommissionPlanService {
	return &promotionCommissionPlanService{
		BaseService: newService(params...),
	}
}
