package product_service

import (
	"context"
	product_service "testLmd/services/product-service/proto/product-service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.nhat.io/grpcmock"
	"google.golang.org/grpc/codes"
)

type TestCase struct {
	scenario         string
	request          interface{}
	expectedResponse interface{}
	response         interface{}
	responseOut      interface{}
	expectedMsgError string
	MsgError         string
	ErrorStatus      codes.Code
	method           string
}

const (
	Success                  = "success"
	Error                    = "error"
	reservationProduct       = "/pb.ProductService/ReservationProduct"
	cancelReservationProduct = "/pb.ProductService/CancelReservationProduct"
	getAllProducts           = "/pb.ProductService/GetAllProducts"
)

func createServiceServerMock(testCase TestCase) grpcmock.ServerMockerWithContextDialer {
	opts := grpcmock.RegisterService(product_service.RegisterProductServiceServer)
	if testCase.MsgError != "" {
		return grpcmock.MockServerWithBufConn(opts, func(s *grpcmock.Server) {
			s.ExpectUnary(testCase.method).
				WithPayload(testCase.request).
				ReturnError(testCase.ErrorStatus, testCase.MsgError)
		})
	}

	return grpcmock.MockServerWithBufConn(opts, func(s *grpcmock.Server) {
		s.ExpectUnary(testCase.method).
			WithPayload(testCase.request).
			Return(testCase.response)
	})
}

func TestConcert(t *testing.T) {
	t.Parallel()
	testCases := []TestCase{
		{
			scenario: Success,
			method:   reservationProduct,
			request:  &product_service.ReservationProductRequest{WarehouseId: 1, Id: []int64{2, 3, 5, 7, 11, 13}},
			expectedResponse: &product_service.ReservationProductResponse{ProductId: map[string]string{
				"Товар номер 1 с id - 2":  " Поставлен в резерв",
				"Товар номер 2 с id - 3":  " Поставлен в резерв",
				"Товар номер 3 с id - 5":  " Поставлен в резерв",
				"Товар номер 4 с id - 7":  " Поставлен в резерв",
				"Товар номер 5 с id - 11": " Поставлен в резерв",
				"Товар номер 6 с id - 13": " Поставлен в резерв",
			}},
			response: &product_service.ReservationProductResponse{ProductId: map[string]string{
				"Товар номер 1 с id - 2":  " Поставлен в резерв",
				"Товар номер 2 с id - 3":  " Поставлен в резерв",
				"Товар номер 3 с id - 5":  " Поставлен в резерв",
				"Товар номер 4 с id - 7":  " Поставлен в резерв",
				"Товар номер 5 с id - 11": " Поставлен в резерв",
				"Товар номер 6 с id - 13": " Поставлен в резерв",
			}},
			responseOut: &product_service.ReservationProductResponse{},
		},
		{
			scenario: Success,
			method:   cancelReservationProduct,
			request:  &product_service.CancelReservationProductRequest{WarehouseId: 1, Id: []int64{2, 3, 5, 7, 11, 13}},
			expectedResponse: &product_service.CancelReservationProductResponse{ProductId: map[string]string{
				"Товар номер 1 с id - 2":  " снят с резерва",
				"Товар номер 2 с id - 3":  " снят с резерва",
				"Товар номер 3 с id - 5":  " снят с резерва",
				"Товар номер 4 с id - 7":  " снят с резерва",
				"Товар номер 5 с id - 11": " снят с резерва",
				"Товар номер 6 с id - 13": " снят с резерва",
			}},
			response: &product_service.CancelReservationProductResponse{ProductId: map[string]string{
				"Товар номер 1 с id - 2":  " снят с резерва",
				"Товар номер 2 с id - 3":  " снят с резерва",
				"Товар номер 3 с id - 5":  " снят с резерва",
				"Товар номер 4 с id - 7":  " снят с резерва",
				"Товар номер 5 с id - 11": " снят с резерва",
				"Товар номер 6 с id - 13": " снят с резерва",
			}},
			responseOut: &product_service.CancelReservationProductResponse{},
		},
		{
			scenario: Success,
			method:   getAllProducts,
			request:  &product_service.GetAllProductsRequest{WarehouseId: 1},
			expectedResponse: &product_service.GetAllProductsResponse{Product: []*product_service.Product{{
				Id:          1,
				Name:        "test name",
				Size:        "test size",
				Quantity:    "test quantity",
				InReserve:   "0",
				WarehouseId: 1,
			}}},
			response: &product_service.GetAllProductsResponse{Product: []*product_service.Product{{
				Id:          1,
				Name:        "test name",
				Size:        "test size",
				Quantity:    "test quantity",
				InReserve:   "0",
				WarehouseId: 1,
			}}},
			responseOut: &product_service.GetAllProductsResponse{},
		},
		{
			scenario:         Error,
			method:           reservationProduct,
			request:          &product_service.ReservationProductRequest{WarehouseId: 0},
			ErrorStatus:      404,
			expectedMsgError: "rpc error: code = Code(404) desc = Record not found",
			MsgError:         "Record not found",
		},
		{
			scenario:         Error,
			method:           cancelReservationProduct,
			request:          &product_service.CancelReservationProductRequest{WarehouseId: 0},
			ErrorStatus:      404,
			expectedMsgError: "rpc error: code = Code(404) desc = Record not found",
			MsgError:         "Record not found",
		},
		{
			scenario:         Error,
			method:           getAllProducts,
			request:          &product_service.GetAllProductsRequest{WarehouseId: 0},
			ErrorStatus:      404,
			expectedMsgError: "rpc error: code = Code(404) desc = Record not found",
			MsgError:         "Record not found",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			_, dialer := createServiceServerMock(tc)(t)

			err := grpcmock.InvokeUnary(context.Background(),
				tc.method, tc.request, tc.responseOut,
				grpcmock.WithInsecure(),
				grpcmock.WithContextDialer(dialer),
			)

			if tc.MsgError != "" {
				t.Log(tc.expectedMsgError)
				t.Log(err)

				assert.EqualError(t, err, tc.expectedMsgError)
				return
			}

			require.NoError(t, err)

			t.Log(tc.expectedResponse)
			t.Log(tc.responseOut)

			assert.Equal(t, tc.expectedResponse, tc.responseOut)
		})
	}
}
