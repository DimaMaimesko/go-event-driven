sequenceDiagram
autonumber
participant Client
participant API as Tickets API
participant CB as Command Bus
participant CH as RefundTicket Handler
participant Receipts
participant Payments
participant EB as Event Bus
participant Proj as Refund Projection(s)

    Client->>API: POST /tickets/{id}/refund
    API->>CB: Publish RefundTicket(ticket_id, idempotency_key)
    CB->>CH: Deliver command
    CH->>Receipts: VoidReceipt(ticket_id, idempotency_key)
    Receipts-->>CH: OK
    CH->>Payments: RefundPayment(ticket_id, idempotency_key)
    Payments-->>CH: OK
    CH->>EB: Publish TicketRefunded(ticket_id)
    EB->>Proj: Deliver TicketRefunded
    Proj-->>EB: ACK
    API-->>Client: 202 Accepted