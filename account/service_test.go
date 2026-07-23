package account

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

func TestAccountEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/account/details/wallet-balance":
			_, _ = w.Write([]byte(`{"data":{"balance_amount":"-291539.83"}}`))
		case "/v1/external/account/details/statement":
			if got := r.URL.Query().Encode(); got != "from=2017-08-12&page=5&per_page=2&to=2017-09-12" {
				t.Fatalf("unexpected query: %s", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"transaction_id":"","order_id":"","channel_order_id":"","awb_code":"","return_awb_code":null,"applied_weight":"","charged_weight":"","billed_weight":"","action":"","charge":"","description":"Wallet Balance","debit_amount":"","credit_amount":"","balance_amount":"0","balance_weight":0,"volumetric_weight":"","entered_weight":"","created_at":"","can_ship":true}]}`))
		case "/v1/external/billing/discrepancy":
			_, _ = w.Write([]byte(`{"status":200,"data":[],"upper_fold_text":"Upper text","lower_fild_text":"Lower text"}`))
		case "/v1/external/errors/20212061/check":
			_, _ = w.Write([]byte(`{"data":{"status":"3","message":"Error in reading file data!"}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	s := NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
	wallet, err := s.GetWalletBalance(context.Background())
	if err != nil || wallet.Data.BalanceAmount.String() != "-291539.83" {
		t.Fatalf("unexpected wallet response: %+v err=%v", wallet, err)
	}
	statement, err := s.GetStatement(context.Background(), &StatementParams{Page: 5, PerPage: 2, From: "2017-08-12", To: "2017-09-12"})
	if err != nil || len(statement.Data) != 1 {
		t.Fatalf("unexpected statement response: %+v err=%v", statement, err)
	}
	discrepancy, err := s.GetDiscrepancy(context.Background())
	if err != nil || discrepancy.Status != 200 {
		t.Fatalf("unexpected discrepancy response: %+v err=%v", discrepancy, err)
	}
	importResp, err := s.CheckImport(context.Background(), &ImportCheckRequest{ImportID: "20212061"})
	if err != nil || importResp.Data.Status != "3" {
		t.Fatalf("unexpected import response: %+v err=%v", importResp, err)
	}
}
