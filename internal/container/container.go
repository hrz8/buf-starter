package container

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	greeterv1 "github.com/hrz8/altalune/gen/greeter/v1"

	api_key_domain "github.com/hrz8/altalune/internal/domain/api_key"
	chatbot_domain "github.com/hrz8/altalune/internal/domain/chatbot"
	chatbot_node_domain "github.com/hrz8/altalune/internal/domain/chatbot_node"
	employee_domain "github.com/hrz8/altalune/internal/domain/employee"
	greeter_domain "github.com/hrz8/altalune/internal/domain/greeter"
	iam_mapper_domain "github.com/hrz8/altalune/internal/domain/iam_mapper"
	oauth_auth_domain "github.com/hrz8/altalune/internal/domain/oauth_auth"
	oauth_client_domain "github.com/hrz8/altalune/internal/domain/oauth_client"
	oauth_provider_domain "github.com/hrz8/altalune/internal/domain/oauth_provider"
	permission_domain "github.com/hrz8/altalune/internal/domain/permission"
	project_domain "github.com/hrz8/altalune/internal/domain/project"
	role_domain "github.com/hrz8/altalune/internal/domain/role"
	user_domain "github.com/hrz8/altalune/internal/domain/user"
	"github.com/hrz8/altalune/internal/postgres"
	"github.com/hrz8/altalune/internal/session"
	"github.com/hrz8/altalune/internal/shared/jwt"
	"github.com/hrz8/altalune/internal/shared/notification"
	"github.com/hrz8/altalune/internal/shared/notification/email"
	"github.com/hrz8/altalune/logger"

	migration_domain "github.com/hrz8/altalune/internal/domain/migration"
)

// Container manages dependency injection with private fields
type Container struct {
	// Configuration and logger
	config altalune.Config
	logger altalune.Logger

	// Database connection and manager
	db postgres.DB

	// Migrations
	migrationRepo    migration_domain.Migrator
	migrationService *migration_domain.Service

	// Example Repositories
	greeterRepo  greeter_domain.Repositor
	employeeRepo employee_domain.Repositor

	// API Key Repository
	apiKeyRepo api_key_domain.Repositor

	// Chatbot Repository
	chatbotRepo chatbot_domain.Repositor

	// Chatbot Node Repository
	chatbotNodeRepo chatbot_node_domain.Repositor

	// IAM Repositories
	userRepo          user_domain.Repository
	roleRepo          role_domain.Repository
	permissionRepo    permission_domain.Repository
	iamMapperRepo     iam_mapper_domain.Repository
	oauthProviderRepo oauth_provider_domain.Repository
	oauthClientRepo   oauth_client_domain.Repositor
	oauthAuthRepo     oauth_auth_domain.Repositor

	// OTP and Verification Repositories
	otpRepo              oauth_auth_domain.OTPRepositor
	otpUserRepo          oauth_auth_domain.UserLookupRepositor
	verificationUserRepo oauth_auth_domain.UserEmailVerificationRepositor
	verificationRepo     oauth_auth_domain.EmailVerificationRepositor

	// Repositories
	projectRepo project_domain.Repositor

	// Shared Providers (available across the app)
	notificationService *notification.NotificationService

	// Example Services
	greeterService  greeterv1.GreeterServiceServer
	employeeService altalunev1.EmployeeServiceServer

	// Domain Services
	projectService       altalunev1.ProjectServiceServer
	apiKeyService        altalunev1.ApiKeyServiceServer
	chatbotService       altalunev1.ChatbotServiceServer
	chatbotNodeService   altalunev1.ChatbotNodeServiceServer
	userService          altalunev1.UserServiceServer
	roleService          altalunev1.RoleServiceServer
	permissionService    altalunev1.PermissionServiceServer
	iamMapperService     altalunev1.IAMMapperServiceServer
	oauthProviderService altalunev1.OAuthProviderServiceServer
	oauthClientService   altalunev1.OAuthClientServiceServer

	// Auth Server Components (conditionally initialized)
	jwtSigner                *jwt.Signer
	sessionStore             *session.Store
	oauthAuthService         *oauth_auth_domain.Service
	otpService               *oauth_auth_domain.OTPService
	emailVerificationService *oauth_auth_domain.EmailVerificationService
}

// CreateContainer creates a new dependency injection container with proper error handling
func CreateContainer(ctx context.Context, cfg altalune.Config) (*Container, error) {
	container := &Container{
		config: cfg,
		logger: logger.New(cfg.GetServerLogLevel()),
	}

	// Initialize components in dependency order:
	// 1. Database connection
	// 2. Repositories (data access layer)
	// 3. Providers (shared infrastructure services like notification)
	// 4. Services (domain business logic)
	// 5. Auth components (auth-specific services)
	if err := container.initDatabase(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	if err := container.initRepositories(); err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}
	if err := container.initProviders(); err != nil {
		return nil, fmt.Errorf("failed to initialize providers: %w", err)
	}
	if err := container.initServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}
	return container, nil
}

// Private initialization methods
func (c *Container) initDatabase(ctx context.Context) error {
	conn := postgres.MustConnect(postgres.ConnectionOptions{
		URL:            c.config.GetDatabaseURL(),
		MaxConnections: c.config.GetDatabaseMaxConnections(),
		MaxIdleTime:    c.config.GetDatabaseMaxIdleTime(),
		ConnectTimeout: c.config.GetDatabaseConnectTimeout(),
	})
	if err := conn.TestConnection(ctx); err != nil {
		return fmt.Errorf("database connection test failed: %w", err)
	}
	c.db = conn
	return nil
}

func (c *Container) initRepositories() error {
	c.migrationRepo = migration_domain.NewAltaluneMigrationRepo(c.db)
	c.greeterRepo = greeter_domain.NewRepo()
	c.employeeRepo = employee_domain.NewRepo(c.db)
	c.projectRepo = project_domain.NewRepo(c.db)
	c.apiKeyRepo = api_key_domain.NewRepo(c.db)
	c.chatbotRepo = chatbot_domain.NewRepo(c.db)
	c.chatbotNodeRepo = chatbot_node_domain.NewRepo(c.db)
	c.userRepo = user_domain.NewRepo(c.db)
	c.roleRepo = role_domain.NewRepo(c.db)
	c.permissionRepo = permission_domain.NewRepo(c.db)
	c.iamMapperRepo = iam_mapper_domain.NewRepo(c.db)
	c.oauthProviderRepo = oauth_provider_domain.NewRepo(c.db, c.config.GetIAMEncryptionKey())
	c.oauthClientRepo = oauth_client_domain.NewRepo(c.db)
	c.oauthAuthRepo = oauth_auth_domain.NewRepo(c.db)

	// OTP and Verification repositories
	c.otpRepo = oauth_auth_domain.NewOTPRepo(c.db)
	userRepo := oauth_auth_domain.NewUserRepo(c.db)
	c.otpUserRepo = userRepo          // UserLookupRepositor for OTP service
	c.verificationUserRepo = userRepo // UserEmailVerificationRepositor for verification service
	c.verificationRepo = oauth_auth_domain.NewEmailVerificationRepo(c.db)
	return nil
}

func (c *Container) initProviders() error {
	emailProvider := c.config.GetNotificationEmailProvider()
	if emailProvider != "" {
		var emailSender email.EmailSender

		fromEmail := c.config.GetNotificationEmailFromEmail()
		fromName := c.config.GetNotificationEmailFromName()

		switch emailProvider {
		case "resend":
			apiKey := c.config.GetNotificationResendAPIKey()
			if apiKey != "" {
				emailSender = email.NewResendEmailSender(apiKey, fromEmail, fromName)
			}
		case "ses":
			region := c.config.GetNotificationSESRegion()
			if region != "" {
				emailSender = email.NewSESEmailSender(region, fromEmail)
			}
		}

		if emailSender != nil {
			baseURL := c.config.GetNotificationAuthBaseURL()
			if baseURL == "" {
				return fmt.Errorf("notification.authBaseURL is required when email provider is configured")
			}
			notificationSvc, err := notification.NewNotificationService(emailSender, baseURL)
			if err != nil {
				return fmt.Errorf("create notification service: %w", err)
			}
			c.notificationService = notificationSvc
		}
	}

	return nil
}

func (c *Container) initServices() error {
	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("failed to create validator: %w", err)
	}
	c.migrationService = migration_domain.NewService(c.logger, c.migrationRepo)
	c.greeterService = greeter_domain.NewService(validator, c.logger, c.greeterRepo)
	c.employeeService = employee_domain.NewService(validator, c.logger, c.projectRepo, c.employeeRepo)
	c.projectService = project_domain.NewService(validator, c.logger, c.projectRepo)
	c.apiKeyService = api_key_domain.NewService(validator, c.logger, c.projectRepo, c.apiKeyRepo)
	c.chatbotService = chatbot_domain.NewService(validator, c.logger, c.projectRepo, c.chatbotRepo)
	c.chatbotNodeService = chatbot_node_domain.NewService(validator, c.logger, c.projectRepo, c.chatbotNodeRepo)
	c.roleService = role_domain.NewService(validator, c.logger, c.roleRepo)
	c.permissionService = permission_domain.NewService(validator, c.logger, c.permissionRepo)
	c.iamMapperService = iam_mapper_domain.NewService(validator, c.logger, c.iamMapperRepo, c.userRepo, c.roleRepo, c.permissionRepo, c.projectRepo)
	c.oauthProviderService = oauth_provider_domain.NewService(validator, c.logger, c.oauthProviderRepo)
	c.oauthClientService = oauth_client_domain.NewService(validator, c.logger, c.projectRepo, c.oauthClientRepo)

	if err := c.initAuthComponents(); err != nil {
		return fmt.Errorf("failed to initialize auth components: %w", err)
	}

	c.userService = user_domain.NewService(validator, c.logger, c.userRepo, c.roleRepo, c.iamMapperRepo, c.emailVerificationService)

	return nil
}

func (c *Container) initAuthComponents() error {
	// JWT Signer - only initialize if key paths are configured
	if c.config.GetJWTPrivateKeyPath() != "" && c.config.GetJWTPublicKeyPath() != "" {
		signer, err := jwt.NewSigner(
			c.config.GetJWTPrivateKeyPath(),
			c.config.GetJWTPublicKeyPath(),
			c.config.GetJWKSKid(),
			c.config.GetJWTIssuer(),
		)
		if err != nil {
			return fmt.Errorf("create jwt signer: %w", err)
		}
		c.jwtSigner = signer
	}

	// Session Store - only initialize if session secret is configured
	if c.config.GetSessionSecret() != "" {
		c.sessionStore = session.NewStore(c.config.GetSessionSecret(), false, 86400)
	}

	// OAuth Auth Service - only initialize if JWT signer is available
	if c.jwtSigner != nil {
		permissionProvider := oauth_auth_domain.NewPermissionService(c.iamMapperRepo)
		scopeHandlerRegistry := oauth_auth_domain.NewScopeHandlerRegistry()

		c.oauthAuthService = oauth_auth_domain.NewService(
			c.logger,
			c.oauthAuthRepo,
			c.otpUserRepo,
			c.jwtSigner,
			c.config,
			permissionProvider,
			scopeHandlerRegistry,
		)
	}

	// Initialize OTP and Email Verification Services if notification service is available
	if c.notificationService != nil {
		c.otpService = oauth_auth_domain.NewOTPService(
			c.otpRepo,
			c.otpUserRepo,
			c.notificationService,
			c.logger,
			c.config,
		)

		c.emailVerificationService = oauth_auth_domain.NewEmailVerificationService(
			c.verificationRepo,
			c.verificationUserRepo,
			c.notificationService,
			c.logger,
			c.config,
		)
	}

	return nil
}
