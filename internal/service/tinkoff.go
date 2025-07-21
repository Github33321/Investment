package service

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/vodolaz095/go-investAPI/investapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TinkoffClient struct {
	conn        *grpc.ClientConn
	ctx         context.Context
	accountID   string
	operations  investapi.OperationsServiceClient
	instruments investapi.InstrumentsServiceClient
	prices      investapi.MarketDataServiceClient
}

func NewTinkoffClient() *TinkoffClient {
	_ = godotenv.Load()
	token := os.Getenv("TINKOFF_TOKEN")
	if token == "" {
		log.Fatal("TINKOFF_TOKEN не найден в .env")
	}

	conn, err := grpc.Dial("invest-public-api.tinkoff.ru:443",
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", "Bearer "+token)

	usersClient := investapi.NewUsersServiceClient(conn)
	accountsResp, err := usersClient.GetAccounts(ctx, &investapi.GetAccountsRequest{})
	if err != nil || len(accountsResp.Accounts) == 0 {
		log.Fatalf("Ошибка получения аккаунтов: %v", err)
	}

	return &TinkoffClient{
		conn:        conn,
		ctx:         ctx,
		accountID:   accountsResp.Accounts[1].Id,
		operations:  investapi.NewOperationsServiceClient(conn),
		instruments: investapi.NewInstrumentsServiceClient(conn),
		prices:      investapi.NewMarketDataServiceClient(conn),
	}
}

type Operation struct {
	ID           string  `json:"id"`
	Currency     string  `json:"currency"`
	FloatPayment float64 `json:"float_payment"`
	Date         string  `json:"date"`
	Type         string  `json:"type"`
	Operation    string  `json:"operation_type"`
	Figi         string  `json:"figi,omitempty"`
	Quantity     float64 `json:"quantity"`
	Price        float64 `json:"price"`
	IsCanceled   bool    `json:"is_canceled"`
}

func (c *TinkoffClient) GetOperations() ([]Operation, error) {
	req := &investapi.OperationsRequest{
		AccountId: c.accountID,
		From:      timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
		To:        timestamppb.New(time.Now()),
	}
	resp, err := c.operations.GetOperations(c.ctx, req)
	if err != nil {
		return nil, err
	}

	var out []Operation
	for _, op := range resp.Operations {
		payment := float64(op.Payment.Units) + float64(op.Payment.Nano)/1e9
		quantity := float64(op.Quantity)
		price := float64(op.Price.Units) + float64(op.Price.Nano)/1e9

		out = append(out, Operation{
			ID:           op.Id,
			Currency:     op.Currency,
			FloatPayment: payment,
			Date:         op.Date.AsTime().Format("02/01/2006"),
			Type:         op.Type,
			Operation:    op.OperationType.String(),
			Figi:         op.Figi,
			Quantity:     quantity,
			Price:        price,
			IsCanceled:   op.State == investapi.OperationState_OPERATION_STATE_CANCELED,
		})
	}
	return out, nil
}

type StockData struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (c *TinkoffClient) GetFigiPrice(figi string) (StockData, error) {
	instrResp, err := c.instruments.GetInstrumentBy(c.ctx, &investapi.InstrumentRequest{
		IdType: investapi.InstrumentIdType_INSTRUMENT_ID_TYPE_FIGI,
		Id:     figi,
	})
	if err != nil {
		return StockData{}, err
	}

	priceResp, err := c.prices.GetLastPrices(c.ctx, &investapi.GetLastPricesRequest{
		Figi: []string{figi},
	})
	if err != nil || len(priceResp.LastPrices) == 0 {
		return StockData{}, err
	}

	p := priceResp.LastPrices[0].Price
	price := float64(p.Units) + float64(p.Nano)/1e9

	return StockData{
		Name:  instrResp.Instrument.Name,
		Price: price,
	}, nil
}
