flowchart LR
U[Customer / Admin] --> API[Tickets Service API]
API --> APP[Tickets Service]

    APP --> B[(Message Broker)]
    APP --> DB[(PostgreSQL)]
    APP --> EXT1[Payments Service]
    APP --> EXT2[Receipts Service]
    APP --> EXT3[Spreadsheet Service]
    APP --> EXT4[File Service]
    APP --> EXT5[External Ticketing Provider]