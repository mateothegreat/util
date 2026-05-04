package context

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/localsession"
)

// ExampleBasicSessionUsage demonstrates the basic usage of Session interface with both SessionCtx and SessionMap implementations.
//
// This example shows how to create sessions, store and retrieve values, and understand the differences
// between SessionCtx (isolated per goroutine) and SessionMap (shared across goroutines).
func ExampleBasicSessionUsage() {
	fmt.Println("=== Basic Session Usage Example ===")

	// Example 1: Using SessionCtx (context-based, isolated storage)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "userID", "12345")
	ctx = context.WithValue(ctx, "requestID", "req-001")

	sessionCtx := localsession.NewSessionCtx(ctx)
	fmt.Printf("SessionCtx - UserID: %v\n", sessionCtx.Get("userID"))
	fmt.Printf("SessionCtx - RequestID: %v\n", sessionCtx.Get("requestID"))

	// Add new value to session
	newSessionCtx := sessionCtx.WithValue("role", "admin")
	fmt.Printf("SessionCtx - Role: %v\n", newSessionCtx.Get("role"))

	// Example 2: Using SessionMap (map-based, shared storage)
	sessionMap := localsession.NewSessionMap(map[interface{}]interface{}{
		"apiKey":    "secret-key-123",
		"rateLimit": 1000,
	})
	fmt.Printf("\nSessionMap - API Key: %v\n", sessionMap.Get("apiKey"))
	fmt.Printf("SessionMap - Rate Limit: %v\n", sessionMap.Get("rateLimit"))

	// Add new value to session map
	newSessionMap := sessionMap.WithValue("region", "us-west-2")
	fmt.Printf("SessionMap - Region: %v\n", newSessionMap.Get("region"))
}

// ExampleSessionManager demonstrates how to use SessionManager to manage sessions globally.
//
// This example shows how to create a session manager, bind sessions to the current goroutine,
// and retrieve session data without explicitly passing context.
func ExampleSessionManager() {
	fmt.Println("\n=== Session Manager Example ===")

	// Create a new session manager with custom options
	manager := localsession.NewSessionManager(localsession.ManagerOptions{
		ShardNumber:                   10,
		EnableImplicitlyTransmitAsync: false, // We'll enable this in another example
		GCInterval:                    time.Hour,
	})

	// Initialize a session with some user context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "tenantID", "tenant-123")
	ctx = context.WithValue(ctx, "userEmail", "user@example.com")

	session := localsession.NewSessionCtx(ctx)
	session = session.WithValue("permissions", []string{"read", "write"})

	// Bind the session to current goroutine
	manager.BindSession(session)

	// Now we can access the session from anywhere in the same goroutine
	currentSession, ok := manager.CurSession()
	if !ok {
		panic("Failed to get current session")
	}

	fmt.Printf("Current Tenant ID: %v\n", currentSession.Get("tenantID"))
	fmt.Printf("Current User Email: %v\n", currentSession.Get("userEmail"))
	fmt.Printf("Current Permissions: %v\n", currentSession.Get("permissions"))

	// Clean up
	manager.UnbindSession()
}

// ExampleExplicitAsyncTransmission demonstrates how to explicitly transmit sessions to child goroutines.
//
// This example shows the recommended way to pass session context to goroutines using Go() and GoSession() methods.
func ExampleExplicitAsyncTransmission() {
	fmt.Println("\n=== Explicit Async Transmission Example ===")

	// Initialize the default manager
	localsession.InitDefaultManager(localsession.DefaultManagerOptions())

	// Create initial session
	session := localsession.NewSessionMap(map[interface{}]interface{}{
		"traceID":     "trace-123",
		"spanID":      "span-001",
		"serviceName": "api-gateway",
	})

	// Bind to current goroutine
	localsession.BindSession(session)

	var wg sync.WaitGroup
	wg.Add(2)

	// Example 1: Using Go() to inherit parent session
	localsession.Go(func() {
		defer wg.Done()

		currentSession, _ := localsession.CurSession()
		fmt.Printf("Child goroutine 1 - TraceID: %v\n", currentSession.Get("traceID"))
		fmt.Printf("Child goroutine 1 - ServiceName: %v\n", currentSession.Get("serviceName"))

		// Create a new session with additional context
		newSession := currentSession.WithValue("spanID", "span-002")
		newSession = newSession.WithValue("handlerName", "userHandler")

		// Launch another goroutine with the new session
		wg.Add(1)
		localsession.GoSession(newSession, func() {
			defer wg.Done()

			currentSession, _ := localsession.CurSession()
			fmt.Printf("Nested goroutine - TraceID: %v\n", currentSession.Get("traceID"))
			fmt.Printf("Nested goroutine - SpanID: %v\n", currentSession.Get("spanID"))
			fmt.Printf("Nested goroutine - HandlerName: %v\n", currentSession.Get("handlerName"))
		})
	})

	// Example 2: Using Go() for parallel processing
	localsession.Go(func() {
		defer wg.Done()

		currentSession, _ := localsession.CurSession()
		fmt.Printf("Child goroutine 2 - TraceID: %v\n", currentSession.Get("traceID"))

		// Simulate some work
		time.Sleep(10 * time.Millisecond)
		fmt.Println("Child goroutine 2 completed work")
	})

	wg.Wait()
	localsession.UnbindSession()
}

// ExampleImplicitAsyncTransmission demonstrates implicit session transmission across goroutines.
//
// When EnableImplicitlyTransmitAsync is true, every goroutine automatically inherits its parent's session.
func ExampleImplicitAsyncTransmission() {
	fmt.Println("\n=== Implicit Async Transmission Example ===")

	// Reset manager with implicit transmission enabled
	localsession.ResetDefaultManager(localsession.ManagerOptions{
		ShardNumber:                   10,
		EnableImplicitlyTransmitAsync: true, // Enable implicit transmission
		GCInterval:                    time.Hour,
	})

	// Create a session with request context
	session := localsession.NewSessionMap(map[interface{}]interface{}{
		"requestID":   "req-456",
		"authToken":   "bearer-xyz",
		"startTime":   time.Now(),
		"environment": "production",
	})

	localsession.BindSession(session)

	var wg sync.WaitGroup
	wg.Add(1)

	// Regular go routine - session is automatically inherited
	go func() {
		defer wg.Done()

		currentSession, ok := localsession.CurSession()
		if !ok {
			panic("Failed to get inherited session")
		}

		fmt.Printf("Implicitly inherited - RequestID: %v\n", currentSession.Get("requestID"))
		fmt.Printf("Implicitly inherited - Environment: %v\n", currentSession.Get("environment"))

		// Nested goroutine also inherits automatically
		wg.Add(1)
		go func() {
			defer wg.Done()

			currentSession, _ := localsession.CurSession()
			fmt.Printf("Nested implicit - RequestID: %v\n", currentSession.Get("requestID"))

			// Add some processing metadata
			newSession := currentSession.WithValue("processingStage", "validation")
			localsession.BindSession(newSession)

			fmt.Printf("Nested implicit - Processing Stage: %v\n", newSession.Get("processingStage"))
		}()
	}()

	wg.Wait()
	localsession.UnbindSession()
}

// ExampleRealWorldHTTPMiddleware demonstrates a real-world use case for session management in HTTP middleware.
//
// This example shows how to use localsession to propagate request context through middleware chains
// and handlers without explicitly passing context parameters.
func ExampleRealWorldHTTPMiddleware() {
	fmt.Println("\n=== Real-World HTTP Middleware Example ===")

	// Initialize default manager
	localsession.InitDefaultManager(localsession.DefaultManagerOptions())

	// Simulate HTTP request processing
	simulateHTTPRequest := func(path string, userID string) {
		// Authentication middleware
		authMiddleware := func() {
			session := localsession.NewSessionMap(map[interface{}]interface{}{
				"userID":        userID,
				"authenticated": true,
				"roles":         []string{"user", "admin"},
			})
			localsession.BindSession(session)
			fmt.Printf("[Auth Middleware] User %s authenticated\n", userID)
		}

		// Logging middleware
		loggingMiddleware := func() {
			currentSession, _ := localsession.CurSession()
			requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())

			newSession := currentSession.WithValue("requestID", requestID)
			newSession = newSession.WithValue("path", path)
			newSession = newSession.WithValue("timestamp", time.Now())

			localsession.BindSession(newSession)
			fmt.Printf("[Logging Middleware] Request ID: %s, Path: %s\n", requestID, path)
		}

		// Rate limiting middleware
		rateLimitMiddleware := func() {
			currentSession, _ := localsession.CurSession()
			userID := currentSession.Get("userID")

			// Simulate rate limit check
			newSession := currentSession.WithValue("rateLimit", 100)
			newSession = newSession.WithValue("remainingRequests", 95)

			localsession.BindSession(newSession)
			fmt.Printf("[RateLimit Middleware] User %v has 95 requests remaining\n", userID)
		}

		// Business logic handler
		businessHandler := func() {
			currentSession, _ := localsession.CurSession()

			// Access all accumulated context without explicit passing
			userID := currentSession.Get("userID")
			requestID := currentSession.Get("requestID")
			roles := currentSession.Get("roles")
			path := currentSession.Get("path")

			fmt.Printf("[Business Handler] Processing request %v for user %v\n", requestID, userID)
			fmt.Printf("[Business Handler] User roles: %v, Path: %v\n", roles, path)

			// Simulate async operation
			var wg sync.WaitGroup
			wg.Add(1)

			localsession.Go(func() {
				defer wg.Done()

				// Context is available in async operations too
				currentSession, _ := localsession.CurSession()
				fmt.Printf("[Async Operation] Processing for request %v\n", currentSession.Get("requestID"))
			})

			wg.Wait()
		}

		// Execute middleware chain
		authMiddleware()
		loggingMiddleware()
		rateLimitMiddleware()
		businessHandler()

		// Cleanup
		localsession.UnbindSession()
		fmt.Println("Request completed")
	}

	// Simulate multiple requests
	simulateHTTPRequest("/api/users", "user-123")
	fmt.Println()
	simulateHTTPRequest("/api/orders", "user-456")
}

// ExampleSessionLifecycle demonstrates session lifecycle management including validation and cleanup.
//
