package lobster

import "fmt"
import "net/http"

type PaymentInterface interface {
	Payment(w http.ResponseWriter, r *http.Request, db *Database, frameParams FrameParams, userId int, username string, amount float64)
}

var paymentInterfaces map[string]PaymentInterface = make(map[string]PaymentInterface)

func paymentMethodList() []string {
	var methods []string
	for method, _ := range paymentInterfaces {
		methods = append(methods, method)
	}
	return methods
}

func paymentHandle(method string, w http.ResponseWriter, r *http.Request, db *Database, frameParams FrameParams, userId int, username string, amount float64) {
	if amount < cfg.Billing.DepositMinimum || amount > cfg.Billing.DepositMaximum {
		redirectMessage(w, r, "/panel/billing", fmt.Sprintf("Error: amount must be between $%.2f and $%.2f.", cfg.Billing.DepositMinimum, cfg.Billing.DepositMaximum))
		return
	}

	payInterface, ok := paymentInterfaces[method]
	if ok {
		payInterface.Payment(w, r, db, frameParams, userId, username, amount)
	} else {
		redirectMessage(w, r, "/panel/billing", "Error: invalid payment method specified.")
	}
}
