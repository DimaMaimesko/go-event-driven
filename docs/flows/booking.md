mermaid sequenceDiagram
autonumber
participant Client
participant API
participant EB as Event Bus
participant H1 as BookPlace Handler
participant H2 as IssueReceipt Handler
participant H3 as PrintTicket Handler
participant H4 as StoreTickets Handler
participant Ext as External Services

    Client->>API: Book tickets
    API->>EB: Publish BookingMade
    EB->>H1: BookingMade
    H1->>Ext: Reserve seats
    H1-->>EB: ACK

    API->>EB: Publish TicketBookingConfirmed
    EB->>H2: TicketBookingConfirmed
    H2->>Ext: Issue receipt
    EB->>H3: TicketBookingConfirmed
    H3->>Ext: Generate ticket file
    EB->>H4: TicketBookingConfirmed
    H4->>Ext: Persist/store ticket data