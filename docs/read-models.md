mermaid flowchart LR
E[Domain Events] --> P1[Projection: Refund Sheet]
E --> P2[Projection: Tracker/Analytics]
E --> P3[Projection: Ticket Storage Read Model]
P1 --> RM1[(Read DB / Sheet)]
P2 --> RM2[(Analytics Store)]
P3 --> RM3[(Ticket Read DB)]