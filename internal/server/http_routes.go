package server

import (
	"encoding/json"
	"net/http"

	"connectrpc.com/connect"
	"github.com/hrz8/altalune/gen/altalune/v1/altalunev1connect"
	"github.com/hrz8/altalune/gen/greeter/v1/greeterv1connect"
	"github.com/hrz8/altalune/internal/auth"
	api_key_domain "github.com/hrz8/altalune/internal/domain/api_key"
	chatbot_domain "github.com/hrz8/altalune/internal/domain/chatbot"
	chatbot_node_domain "github.com/hrz8/altalune/internal/domain/chatbot_node"
	config_domain "github.com/hrz8/altalune/internal/domain/config"
	employee_domain "github.com/hrz8/altalune/internal/domain/employee"
	greeter_domain "github.com/hrz8/altalune/internal/domain/greeter"
	iam_mapper_domain "github.com/hrz8/altalune/internal/domain/iam_mapper"
	oauth_client_domain "github.com/hrz8/altalune/internal/domain/oauth_client"
	oauth_provider_domain "github.com/hrz8/altalune/internal/domain/oauth_provider"
	permission_domain "github.com/hrz8/altalune/internal/domain/permission"
	project_domain "github.com/hrz8/altalune/internal/domain/project"
	role_domain "github.com/hrz8/altalune/internal/domain/role"
	user_domain "github.com/hrz8/altalune/internal/domain/user"
)

func (s *Server) setupRoutes() *http.ServeMux {
	connectrpcMux := http.NewServeMux()

	// Setup auth interceptor if JWT validator is configured
	var handlerOptions []connect.HandlerOption
	if validator := s.c.GetJWTValidator(); validator != nil {
		authInterceptor := auth.NewAuthInterceptor(validator)
		handlerOptions = append(handlerOptions, connect.WithInterceptors(authInterceptor))
	}

	// Get authorizer for handlers that need authorization checks
	authorizer := s.c.GetAuthorizer()

	// Examples
	greeterHandler := greeter_domain.NewHandler(s.c.GetGreeterService())
	employeeHandler := employee_domain.NewHandler(s.c.GetEmployeeService(), authorizer)
	greeterPath, greeterConnectHandler := greeterv1connect.NewGreeterServiceHandler(greeterHandler, handlerOptions...)
	employeePath, employeeConnectHandler := altalunev1connect.NewEmployeeServiceHandler(employeeHandler, handlerOptions...)
	connectrpcMux.Handle(greeterPath, greeterConnectHandler)
	connectrpcMux.Handle(employeePath, employeeConnectHandler)

	// Domains
	projectHandler := project_domain.NewHandler(s.c.GetProjectService(), authorizer)
	projectPath, projectConnectHandler := altalunev1connect.NewProjectServiceHandler(projectHandler, handlerOptions...)
	connectrpcMux.Handle(projectPath, projectConnectHandler)

	apiKeyHandler := api_key_domain.NewHandler(s.c.GetApiKeyService(), authorizer)
	apiKeyPath, apiKeyConnectHandler := altalunev1connect.NewApiKeyServiceHandler(apiKeyHandler, handlerOptions...)
	connectrpcMux.Handle(apiKeyPath, apiKeyConnectHandler)

	chatbotHandler := chatbot_domain.NewHandler(s.c.GetChatbotService(), authorizer)
	chatbotPath, chatbotConnectHandler := altalunev1connect.NewChatbotServiceHandler(chatbotHandler, handlerOptions...)
	connectrpcMux.Handle(chatbotPath, chatbotConnectHandler)

	chatbotNodeHandler := chatbot_node_domain.NewHandler(s.c.GetChatbotNodeService(), authorizer)
	chatbotNodePath, chatbotNodeConnectHandler := altalunev1connect.NewChatbotNodeServiceHandler(chatbotNodeHandler, handlerOptions...)
	connectrpcMux.Handle(chatbotNodePath, chatbotNodeConnectHandler)

	// IAM Domains
	userHandler := user_domain.NewHandler(s.c.GetUserService(), authorizer)
	userPath, userConnectHandler := altalunev1connect.NewUserServiceHandler(userHandler, handlerOptions...)
	connectrpcMux.Handle(userPath, userConnectHandler)

	roleHandler := role_domain.NewHandler(s.c.GetRoleService(), authorizer)
	rolePath, roleConnectHandler := altalunev1connect.NewRoleServiceHandler(roleHandler, handlerOptions...)
	connectrpcMux.Handle(rolePath, roleConnectHandler)

	permissionHandler := permission_domain.NewHandler(s.c.GetPermissionService(), authorizer)
	permissionPath, permissionConnectHandler := altalunev1connect.NewPermissionServiceHandler(permissionHandler, handlerOptions...)
	connectrpcMux.Handle(permissionPath, permissionConnectHandler)

	iamMapperHandler := iam_mapper_domain.NewHandler(s.c.GetIAMMapperService(), authorizer)
	iamMapperPath, iamMapperConnectHandler := altalunev1connect.NewIAMMapperServiceHandler(iamMapperHandler, handlerOptions...)
	connectrpcMux.Handle(iamMapperPath, iamMapperConnectHandler)

	oauthProviderHandler := oauth_provider_domain.NewHandler(s.c.GetOAuthProviderService(), authorizer)
	oauthProviderPath, oauthProviderConnectHandler := altalunev1connect.NewOAuthProviderServiceHandler(oauthProviderHandler, handlerOptions...)
	connectrpcMux.Handle(oauthProviderPath, oauthProviderConnectHandler)

	oauthClientHandler := oauth_client_domain.NewHandler(s.c.GetOAuthClientService(), authorizer)
	oauthClientPath, oauthClientConnectHandler := altalunev1connect.NewOAuthClientServiceHandler(oauthClientHandler, handlerOptions...)
	connectrpcMux.Handle(oauthClientPath, oauthClientConnectHandler)

	// Public Config (no auth required - register without auth interceptor)
	configHandler := config_domain.NewHandler(s.cfg)
	configPath, configConnectHandler := altalunev1connect.NewConfigServiceHandler(configHandler)
	connectrpcMux.Handle(configPath, configConnectHandler)

	// main server mux
	mux := http.NewServeMux()

	// OAuth BFF endpoints (keeps tokens secure on backend via httpOnly cookies)
	s.registerBFFRoutes(mux)

	// Connect-RPC API routes
	mux.Handle("/api/", http.StripPrefix("/api", connectrpcMux))

	// Health check endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		health := map[string]any{
			"status": "ok",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(health); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	// Static file serving for SPA frontend
	s.registerStaticRoutes(mux)

	return mux
}
