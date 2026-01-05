package container

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	greeterv1 "github.com/hrz8/altalune/gen/greeter/v1"

	api_key_domain "github.com/hrz8/altalune/internal/domain/api_key"
	employee_domain "github.com/hrz8/altalune/internal/domain/employee"
	greeter_domain "github.com/hrz8/altalune/internal/domain/greeter"
	iam_mapper_domain "github.com/hrz8/altalune/internal/domain/iam_mapper"
	permission_domain "github.com/hrz8/altalune/internal/domain/permission"
	project_domain "github.com/hrz8/altalune/internal/domain/project"
	role_domain "github.com/hrz8/altalune/internal/domain/role"
	user_domain "github.com/hrz8/altalune/internal/domain/user"
	"github.com/hrz8/altalune/internal/postgres"
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

	// IAM Repositories
	userRepo       user_domain.Repository
	roleRepo       role_domain.Repository
	permissionRepo permission_domain.Repository
	iamMapperRepo  iam_mapper_domain.Repository

	// Example Services
	greeterService  greeterv1.GreeterServiceServer
	employeeService altalunev1.EmployeeServiceServer

	// Repositories
	projectRepo project_domain.Repositor

	// Services
	projectService    altalunev1.ProjectServiceServer
	apiKeyService     altalunev1.ApiKeyServiceServer
	userService       altalunev1.UserServiceServer
	roleService       altalunev1.RoleServiceServer
	permissionService altalunev1.PermissionServiceServer
	iamMapperService  altalunev1.IAMMapperServiceServer
}

// CreateContainer creates a new dependency injection container with proper error handling
func CreateContainer(ctx context.Context, cfg altalune.Config) (*Container, error) {
	container := &Container{
		config: cfg,
		logger: logger.New(cfg.GetServerLogLevel()),
	}

	// Initialize components in dependency order
	if err := container.initDatabase(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	if err := container.initRepositories(); err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
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
	c.userRepo = user_domain.NewRepo(c.db)
	c.roleRepo = role_domain.NewRepo(c.db)
	c.permissionRepo = permission_domain.NewRepo(c.db)
	c.iamMapperRepo = iam_mapper_domain.NewRepo(c.db)
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
	c.userService = user_domain.NewService(validator, c.logger, c.userRepo)
	c.roleService = role_domain.NewService(validator, c.logger, c.roleRepo)
	c.permissionService = permission_domain.NewService(validator, c.logger, c.permissionRepo)
	c.iamMapperService = iam_mapper_domain.NewService(validator, c.logger, c.iamMapperRepo, c.userRepo, c.roleRepo, c.permissionRepo, c.projectRepo)

	return nil
}
