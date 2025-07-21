package client

import (
	"context"
	"errors"
	"time"

	"github.com/vodolaz095/go-investAPI/investapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TinkoffClient struct {
	token             string
	conn              *grpc.ClientConn
	instrumentsClient investapi.InstrumentsServiceClient
	marketClient      investapi.MarketDataServiceClient
	operationsClient  investapi.OperationsServiceClient
	usersClient       investapi.UsersServiceClient
}

func NewTinkoffClient(token string) *TinkoffClient {
	conn, _ := grpc.Dial("invest-public-api.tinkoff.ru:443",
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))

	return &TinkoffClient{
		token:             token,
		conn:              conn,
		instrumentsClient: investapi.NewInstrumentsServiceClient(conn),
		marketClient:      investapi.NewMarketDataServiceClient(conn),
		operationsClient:  investapi.NewOperationsServiceClient(conn),
		usersClient:       investapi.NewUsersServiceClient(conn),
	}
}

func (c *TinkoffClient) getAuthContext() context.Context {
	return metadata.AppendToOutgoingContext(context.Background(), "Authorization", "Bearer "+c.token)
}

func (c *TinkoffClient) GetOperations() ([]*investapi.Operation, error) {
	ctx := c.getAuthContext()
	accResp, err := c.usersClient.GetAccounts(ctx, &investapi.GetAccountsRequest{})
	if err != nil || len(accResp.Accounts) == 0 {
		return nil, errors.New("не удалось получить аккаунт")
	}

	resp, err := c.operationsClient.GetOperations(ctx, &investapi.OperationsRequest{
		AccountId: accResp.Accounts[1].Id,
		From:      timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
		To:        timestamppb.New(time.Now()),
	})
	if err != nil {
		return nil, err
	}
	return resp.Operations, nil
}

func (c *TinkoffClient) GetFigiPrice(figi string) (string, float64, error) {
	ctx := c.getAuthContext()

	instr, err := c.instrumentsClient.GetInstrumentBy(ctx, &investapi.InstrumentRequest{
		IdType: investapi.InstrumentIdType_INSTRUMENT_ID_TYPE_FIGI,
		Id:     figi,
	})
	if err != nil {
		return "", 0, err
	}

	prices, err := c.marketClient.GetLastPrices(ctx, &investapi.GetLastPricesRequest{Figi: []string{figi}})
	if err != nil || len(prices.LastPrices) == 0 {
		return "", 0, errors.New("цена не найдена")
	}

	p := prices.LastPrices[0].Price
	return instr.Instrument.Name, float64(p.Units) + float64(p.Nano)/1e9, nil
}
