package hyperlocal

import (
	"context"

	"github.com/venom90/shiprocket-go/courier"
	internalclient "github.com/venom90/shiprocket-go/internal/client"
	"github.com/venom90/shiprocket-go/orders"
	"github.com/venom90/shiprocket-go/pickupaddress"
	"github.com/venom90/shiprocket-go/shipment"
)

type Service struct {
	orders          *orders.Service
	couriers        *courier.Service
	shipments       *shipment.Service
	pickupAddresses *pickupaddress.Service
}

func NewService(client *internalclient.Client) *Service {
	return &Service{
		orders:          orders.NewService(client),
		couriers:        courier.NewService(client),
		shipments:       shipment.NewService(client),
		pickupAddresses: pickupaddress.NewService(client),
	}
}

func (s *Service) CreateOrder(ctx context.Context, request *orders.CreateCustomOrderRequest) (*orders.CustomOrderResponse, error) {
	return s.orders.CreateCustomOrder(ctx, request)
}

func (s *Service) ListOrders(ctx context.Context, params *orders.OrdersListParams) (*orders.OrdersListResponse, error) {
	return s.orders.GetOrdersWithParams(ctx, params)
}

func (s *Service) GetOrderDetails(ctx context.Context, request *orders.GetOrderDetailsRequest) (orders.OrderDetailResponse, error) {
	return s.orders.GetOrderDetails(ctx, request)
}

func (s *Service) ExportOrders(ctx context.Context, request *orders.ExportOrdersRequest) (*orders.ExportOrdersResponse, error) {
	return s.orders.ExportOrders(ctx, request)
}

func (s *Service) AssignAWB(ctx context.Context, request *courier.AssignAWBRequest) (*courier.AssignAWBResponse, error) {
	return s.couriers.AssignAWB(ctx, request)
}

func (s *Service) CheckServiceability(ctx context.Context, params *courier.ServiceabilityParams) (*courier.ServiceabilityResponse, error) {
	if params == nil {
		params = &courier.ServiceabilityParams{}
	}
	if params.IsNewHyperlocal == nil {
		isHyperlocal := true
		params.IsNewHyperlocal = &isHyperlocal
	}
	return s.couriers.CheckServiceability(ctx, params)
}

func (s *Service) TrackByAWB(ctx context.Context, request *shipment.TrackByAWBRequest) (*shipment.TrackingResponse, error) {
	return s.shipments.TrackByAWB(ctx, request)
}

func (s *Service) TrackByAWBs(ctx context.Context, request *shipment.TrackByAWBsRequest) (shipment.MultiTrackingResponse, error) {
	return s.shipments.TrackByAWBs(ctx, request)
}

func (s *Service) TrackByShipmentID(ctx context.Context, request *shipment.TrackByShipmentIDRequest) (*shipment.TrackingResponse, error) {
	return s.shipments.TrackByShipmentID(ctx, request)
}

func (s *Service) TrackByOrder(ctx context.Context, request *shipment.TrackByOrderRequest) (shipment.OrderTrackingResponse, error) {
	return s.shipments.TrackByOrder(ctx, request)
}

func (s *Service) ListPickupAddresses(ctx context.Context) (*pickupaddress.ListResponse, error) {
	return s.pickupAddresses.List(ctx)
}

func (s *Service) CreatePickupAddress(ctx context.Context, request *pickupaddress.CreateRequest) (*pickupaddress.CreateResponse, error) {
	return s.pickupAddresses.Create(ctx, request)
}
