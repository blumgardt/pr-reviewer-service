package app

import (
	"fmt"
	"log"
	http2 "net/http"

	"github.com/blumgardt/pr-reviewer-service.git/internal/config"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/handlers/pull_requests"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/handlers/teams"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/handlers/users"
	"github.com/blumgardt/pr-reviewer-service.git/internal/repository/postgres"
	"github.com/blumgardt/pr-reviewer-service.git/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	config       *config.Config
	logger       *log.Logger
	db           *pgxpool.Pool
	Router       *http.Router
	UsersHandler *users.UsersHandler
	TeamHandler  *teams.TeamHandler
	PRHandler    *pull_requests.PullRequestHandler
}

func NewApp(config *config.Config, logger *log.Logger, pool *pgxpool.Pool) *App {
	// Repositories
	teamRepo := postgres.NewTeamRepository(pool)
	usersRepo := postgres.NewUserRepository(pool)
	prRepo := postgres.NewPullRequestRepository(pool)

	// Services
	teamService := service.NewTeamService(teamRepo)
	usersService := service.NewUserService(usersRepo)
	prService := service.NewPullRequestService(prRepo, usersRepo, teamRepo)

	// Handlers
	teamHandler := teams.NewTeamHandler(teamService)
	usersHandler := users.NewUsersHandler(usersService)
	prHandler := pull_requests.NewPullRequestHandler(prService)

	app := &App{
		config:       config,
		logger:       logger,
		Router:       http.NewRouter(),
		UsersHandler: usersHandler,
		TeamHandler:  teamHandler,
		PRHandler:    prHandler,
	}

	app.configureRouter()

	return app
}

func (a *App) Start() error {
	addr := fmt.Sprintf(":%d", a.config.HTTP.Port)
	a.logger.Printf("starting http server on %s", addr)
	return http2.ListenAndServe(addr, a.Router.Handler())
}

func (a *App) configureRouter() {
	// Teams
	a.Router.HandleFunc("/team/add", a.TeamHandler.Add)
	a.Router.HandleFunc("/team/get", a.TeamHandler.Get)

	// Users
	a.Router.HandleFunc("/users/setIsActive", a.UsersHandler.SetIsActive)
	a.Router.HandleFunc("/users/getReview", a.UsersHandler.GetReview)

	// Pull Requests
	a.Router.HandleFunc("/pullRequest/create", a.PRHandler.Create)
	a.Router.HandleFunc("/pullRequest/merge", a.PRHandler.Merge)
	a.Router.HandleFunc("/pullRequest/reassign", a.PRHandler.ReAssign)
}
