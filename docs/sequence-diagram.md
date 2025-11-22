# Flight Search Sequence Diagram

## Main Flight Search Flow

```
    title Flight Search Flow

    participant Client
    participant Controller
    participant Middleware
    participant Usecase
    participant Service
    participant Cache
    participant Provider1 as Garuda API
    participant Provider2 as Lion Air API
    participant Provider3 as AirAsia API
    participant Provider4 as Batik Air API

    Client->>Controller: POST /api/flights/search
    Controller->>Middleware: Rate Limit Check
    Middleware-->>Controller: OK
    
    Controller->>Controller: Validate Request
    Controller->>Usecase: SearchFlights(request)
    
    Usecase->>Cache: Check Cache
    Cache-->>Usecase: Cache Miss
    
    Usecase->>Service: FetchFlights(criteria)
    
    par Concurrent Provider Calls
        Service->>Provider1: Search Flights
        Service->>Provider2: Search Flights  
        Service->>Provider3: Search Flights
        Service->>Provider4: Search Flights
    end
    
    Provider1-->>Service: Flight Data (Format A)
    Provider2-->>Service: Flight Data (Format B)
    Provider3-->>Service: Flight Data (Format C)
    Provider4-->>Service: Flight Data (Format D)
    
    Service->>Service: Normalize Data Formats
    Service-->>Usecase: Aggregated Flight Data
    
    Usecase->>Usecase: Apply Filters
    Usecase->>Usecase: Sort Results
    Usecase->>Usecase: Calculate Best Value
    
    Usecase->>Cache: Store Results
    Usecase-->>Controller: Filtered & Sorted Results
    
    Controller-->>Client: JSON Response
```

## Error Handling Flow

```
    title Error Handling Flow

    participant Client
    participant Controller
    participant Service
    participant Provider1
    participant Provider2
    participant Provider3

    Client->>Controller: POST /api/flights/search
    Controller->>Service: FetchFlights()
    
    par Provider Calls with Retry
        Service->>Provider1: Request
        Provider1-->>Service: Success
        
        Service->>Provider2: Request
        Provider2-->>Service: Timeout
        Service->>Provider2: Retry (1/3)
        Provider2-->>Service: Error
        Service->>Provider2: Retry (2/3)
        Provider2-->>Service: Success
        
        Service->>Provider3: Request
        Provider3-->>Service: Error
        Service->>Provider3: Retry (1/3)
        Provider3-->>Service: Error
        Service->>Provider3: Retry (2/3)
        Provider3-->>Service: Error
        Service->>Provider3: Retry (3/3)
        Provider3-->>Service: Final Error
    end
    
    Service-->>Controller: Partial Results + Warning
    Controller-->>Client: 200 OK with Warning
```

## Cache Hit Flow

```
    title Cache Hit Flow

    participant Client
    participant Controller
    participant Usecase
    participant Cache

    Client->>Controller: POST /api/flights/search
    Controller->>Usecase: SearchFlights(request)
    
    Usecase->>Cache: Get(cacheKey)
    Cache-->>Usecase: Cached Results
    
    Usecase->>Usecase: Apply Dynamic Filters
    Usecase->>Usecase: Sort Results
    
    Usecase-->>Controller: Filtered Results
    Controller-->>Client: JSON Response (Fast)
```