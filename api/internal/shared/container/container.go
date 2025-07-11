package container

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/okta"

	authApp "github.com/BennyEisner/test-results/internal/auth/application"
	authDB "github.com/BennyEisner/test-results/internal/auth/infrastructure/database"
	authHTTP "github.com/BennyEisner/test-results/internal/auth/infrastructure/http"
	authMiddleware "github.com/BennyEisner/test-results/internal/auth/infrastructure/middleware"
	buildApp "github.com/BennyEisner/test-results/internal/build/application"
	buildDB "github.com/BennyEisner/test-results/internal/build/infrastructure/database"
	buildHTTP "github.com/BennyEisner/test-results/internal/build/infrastructure/http"
	buildExecApp "github.com/BennyEisner/test-results/internal/build_test_case_execution/application"
	buildExecDB "github.com/BennyEisner/test-results/internal/build_test_case_execution/infrastructure/database"
	buildExecHTTP "github.com/BennyEisner/test-results/internal/build_test_case_execution/infrastructure/http"
	failureApp "github.com/BennyEisner/test-results/internal/failure/application"
	failureDB "github.com/BennyEisner/test-results/internal/failure/infrastructure/database"
	failureHTTP "github.com/BennyEisner/test-results/internal/failure/infrastructure/http"
	projectApp "github.com/BennyEisner/test-results/internal/project/application"
	projectDB "github.com/BennyEisner/test-results/internal/project/infrastructure/database"
	projectHTTP "github.com/BennyEisner/test-results/internal/project/infrastructure/http"
	"github.com/BennyEisner/test-results/internal/shared/middleware"
	testCaseApp "github.com/BennyEisner/test-results/internal/test_case/application"
	testCaseDB "github.com/BennyEisner/test-results/internal/test_case/infrastructure/database"
	testCaseHTTP "github.com/BennyEisner/test-results/internal/test_case/infrastructure/http"
	testSuiteApp "github.com/BennyEisner/test-results/internal/test_suite/application"
	testSuiteDB "github.com/BennyEisner/test-results/internal/test_suite/infrastructure/database"
	testSuiteHTTP "github.com/BennyEisner/test-results/internal/test_suite/infrastructure/http"
	userApp "github.com/BennyEisner/test-results/internal/user/application"
	userDB "github.com/BennyEisner/test-results/internal/user/infrastructure/database"
	userHTTP "github.com/BennyEisner/test-results/internal/user/infrastructure/http"
	userConfigApp "github.com/BennyEisner/test-results/internal/user_config/application"
	userConfigDB "github.com/BennyEisner/test-results/internal/user_config/infrastructure"
	userConfigHTTP "github.com/BennyEisner/test-results/internal/user_config/infrastructure/http"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter creates and configures the HTTP router with all handlers
func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// --- Setup Goth providers ---
	setupGothProviders()

	// --- Standard Kubernetes health endpoints ---
	// @Summary Liveness probe
	// @Description Kubernetes liveness probe endpoint
	// @Tags health
	// @Produce plain
	// @Success 200 {string} string "ok"
	// @Router /livez [get]
	mux.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			slog.Error("failed to write response", "error", err)
		}
	})

	// @Summary Readiness probe
	// @Description Kubernetes readiness probe endpoint with database connectivity check
	// @Tags health
	// @Produce plain
	// @Success 200 {string} string "ok"
	// @Failure 503 {string} string "db not ready"
	// @Router /readyz [get]
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			if _, err := w.Write([]byte("db not ready")); err != nil {
				slog.Error("failed to write response", "error", err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			slog.Error("failed to write response", "error", err)
		}
	})

	// @Summary Health check
	// @Description Comprehensive health check endpoint with database connectivity check
	// @Tags health
	// @Produce plain
	// @Success 200 {string} string "ok"
	// @Failure 503 {string} string "db not healthy"
	// @Router /healthz [get]
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			if _, err := w.Write([]byte("db not healthy")); err != nil {
				slog.Error("failed to write response", "error", err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			slog.Error("failed to write response", "error", err)
		}
	})

	// --- API Documentation ---
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Wire up repositories
	authRepo := authDB.NewSQLAuthRepository(db)
	projectRepo := projectDB.NewSQLProjectRepository(db)
	buildRepo := buildDB.NewSQLBuildRepository(db)
	buildExecRepo := buildExecDB.NewSQLBuildTestCaseExecutionRepository(db)
	failureRepo := failureDB.NewSQLFailureRepository(db)
	userRepo := userDB.NewSQLUserRepository(db)
	testSuiteRepo := testSuiteDB.NewSQLTestSuiteRepository(db)
	testCaseRepo := testCaseDB.NewSQLTestCaseRepository(db)
	userConfigRepo := userConfigDB.NewSQLUserConfigRepository(db)

	// Wire up services
	authService := authApp.NewAuthService(authRepo)
	projectService := projectApp.NewProjectService(projectRepo)
	buildService := buildApp.NewBuildService(buildRepo)
	buildExecService := buildExecApp.NewBuildTestCaseExecutionService(buildExecRepo)
	failureService := failureApp.NewFailureService(failureRepo)
	userService := userApp.NewUserService(userRepo)
	testSuiteService := testSuiteApp.NewTestSuiteService(testSuiteRepo)
	testCaseService := testCaseApp.NewTestCaseService(testCaseRepo)
	userConfigService := userConfigApp.NewUserConfigService(userConfigRepo)

	// Wire up HTTP handlers
	authHandler := authHTTP.NewAuthHandler(authService)
	projectHandler := projectHTTP.NewProjectHandler(projectService)
	buildHandler := buildHTTP.NewBuildHandler(buildService)
	buildExecHandler := buildExecHTTP.NewBuildTestCaseExecutionHandler(buildExecService)
	failureHandler := failureHTTP.NewFailureHandler(failureService)
	userHandler := userHTTP.NewUserHandler(userService)
	testSuiteHandler := testSuiteHTTP.NewTestSuiteHandler(testSuiteService)
	testCaseHandler := testCaseHTTP.NewTestCaseHandler(testCaseService)
	userConfigHandler := userConfigHTTP.NewUserConfigHandler(userConfigService)

	// Wire up middleware
	authMiddleware := authMiddleware.NewAuthMiddleware(authService)

	// --- API subrouter ---
	apiMux := http.NewServeMux()
	registerRoutes(apiMux, projectHandler, buildHandler,
		buildExecHandler, failureHandler, userHandler, testSuiteHandler, testCaseHandler, userConfigHandler)
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	// --- Auth routes (not under /api prefix) ---
	authMux := http.NewServeMux()
	registerAuthRoutes(authMux, authHandler, authMiddleware)
	mux.Handle("/auth/", http.StripPrefix("/auth", authMux))

	// Apply middleware
	logger := slog.Default()
	return middleware.Cors(middleware.LoggingMiddleware(logger)(mux))
}

// setupGothProviders configures OAuth2 providers for authentication
func setupGothProviders() {
	// Configure session store for Goth
	key := os.Getenv("SESSION_SECRET")
	if key == "" {
		key = "dev-secret-key-change-in-production"
	}

	maxAge := 86400 * 30 // 30 days
	isProd := os.Getenv("ENVIRONMENT") == "production"

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	// Register OAuth2 providers
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Development provider (GitHub)
	if githubClientID := os.Getenv("GITHUB_CLIENT_ID"); githubClientID != "" {
		githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
		goth.UseProviders(
			github.New(githubClientID, githubClientSecret, baseURL+"/auth/github/callback"),
		)
	}

	// Production provider (Okta)
	if oktaClientID := os.Getenv("OKTA_CLIENT_ID"); oktaClientID != "" {
		oktaClientSecret := os.Getenv("OKTA_CLIENT_SECRET")
		oktaDomain := os.Getenv("OKTA_DOMAIN")
		if oktaDomain != "" {
			goth.UseProviders(
				okta.New(oktaClientID, oktaClientSecret, baseURL+"/auth/okta/callback", oktaDomain),
			)
		}
	}
}

// registerRoutes registers all HTTP routes
func registerRoutes(mux *http.ServeMux,
	projectHandler *projectHTTP.ProjectHandler,
	buildHandler *buildHTTP.BuildHandler,
	buildExecHandler *buildExecHTTP.BuildTestCaseExecutionHandler,
	failureHandler *failureHTTP.FailureHandler,
	userHandler *userHTTP.UserHandler,
	testSuiteHandler *testSuiteHTTP.TestSuiteHandler,
	testCaseHandler *testCaseHTTP.TestCaseHandler,
	userConfigHandler *userConfigHTTP.UserConfigHandler) {

	// Project routes
	mux.HandleFunc("GET /projects", projectHandler.GetAllProjects)
	mux.HandleFunc("GET /projects/{id}", projectHandler.GetProjectByID)
	mux.HandleFunc("POST /projects", projectHandler.CreateProject)
	mux.HandleFunc("PUT /projects/{id}", projectHandler.UpdateProject)
	mux.HandleFunc("DELETE /projects/{id}", projectHandler.DeleteProject)

	// Build routes
	mux.HandleFunc("GET /builds", buildHandler.GetBuilds)
	mux.HandleFunc("GET /builds/{id}", buildHandler.GetBuildByID)
	mux.HandleFunc("POST /builds", buildHandler.CreateBuild)
	mux.HandleFunc("PUT /builds/{id}", buildHandler.UpdateBuild)
	mux.HandleFunc("DELETE /builds/{id}", buildHandler.DeleteBuild)

	// Build Test Case Execution routes
	mux.HandleFunc("GET /builds/{buildID}/executions", buildExecHandler.GetExecutionsByBuildID)
	mux.HandleFunc("GET /executions/{id}", buildExecHandler.GetExecutionByID)
	mux.HandleFunc("POST /builds/{buildID}/executions", buildExecHandler.CreateExecution)
	mux.HandleFunc("PUT /executions/{id}", buildExecHandler.UpdateExecution)
	mux.HandleFunc("DELETE /executions/{id}", buildExecHandler.DeleteExecution)

	// Failure routes
	mux.HandleFunc("GET /executions/{executionID}/failure", failureHandler.GetFailureByExecution)
	mux.HandleFunc("GET /failures/{id}", failureHandler.GetFailureByID)
	mux.HandleFunc("POST /executions/{executionID}/failures", failureHandler.CreateFailure)
	mux.HandleFunc("PUT /failures/{id}", failureHandler.UpdateFailure)
	mux.HandleFunc("DELETE /failures/{id}", failureHandler.DeleteFailure)

	// User routes
	mux.HandleFunc("GET /user/{id}", userHandler.GetUserByID)
	mux.HandleFunc("GET /user/username/{username}", userHandler.GetUserByUsername)
	mux.HandleFunc("GET /user/email/{email}", userHandler.GetUserByEmail)
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("PUT /user/{id}", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /user/{id}", userHandler.DeleteUser)

	// Test suite routes
	mux.HandleFunc("GET /test-suites", testSuiteHandler.GetTestSuites)
	mux.HandleFunc("POST /test-suites", testSuiteHandler.CreateTestSuite)
	mux.HandleFunc("PUT /test-suites", testSuiteHandler.UpdateTestSuite)
	mux.HandleFunc("DELETE /test-suites", testSuiteHandler.DeleteTestSuite)

	// Test case routes
	mux.HandleFunc("GET /test-cases", testCaseHandler.GetTestCases)
	mux.HandleFunc("POST /test-cases", testCaseHandler.CreateTestCase)
	mux.HandleFunc("PUT /test-cases", testCaseHandler.UpdateTestCase)
	mux.HandleFunc("DELETE /test-cases", testCaseHandler.DeleteTestCase)

	// User config routes
	mux.HandleFunc("GET /users/{userID}/configs", userConfigHandler.GetUserConfigs)
	mux.HandleFunc("POST /users/{userID}/configs", userConfigHandler.SaveUserConfig)
}

// registerAuthRoutes registers authentication HTTP routes
func registerAuthRoutes(mux *http.ServeMux, authHandler *authHTTP.AuthHandler, authMiddleware *authMiddleware.AuthMiddleware) {
	// OAuth2 authentication routes
	mux.HandleFunc("GET /{provider}", authHandler.BeginOAuth2Auth)
	mux.HandleFunc("GET /{provider}/callback", authHandler.OAuth2Callback)

	// Session management routes
	mux.HandleFunc("POST /logout", authHandler.Logout)

	// User management routes (protected)
	mux.Handle("GET /me", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.GetCurrentUser)))

	// API key management routes (protected)
	mux.Handle("GET /api-keys", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.ListAPIKeys)))
	mux.Handle("POST /api-keys", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.CreateAPIKey)))
	mux.Handle("DELETE /api-keys/{id}", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.DeleteAPIKey)))
}
