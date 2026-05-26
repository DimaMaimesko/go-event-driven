package entities

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type AppendToTrackerPayload struct {
	TicketID      string `json:"ticket_id"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`
}
