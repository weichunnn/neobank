package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/weichunnn/neobank/api"
	db "github.com/weichunnn/neobank/db/sqlc"
	"github.com/weichunnn/neobank/gapi"
	"github.com/weichunnn/neobank/pb"
	"github.com/weichunnn/neobank/util"
	"github.com/weichunnn/neobank/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/weichunnn/neobank/docs/statik"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot load configuration")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot connect to db")
	}

	// run db migration
	runDBMigration(config.MigrationURL, config.DBSource)
	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot create a migrate instance")
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).
			Msg("failed to run migrate u")
	}

	log.Info().Msg("db migrated successfully")
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot create server")

	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterNeoBankServer(grpcServer, server) // 2nd param needs to implement the function introduced in gRPC
	// provides information about publicly-accessible gRPC services on a server,
	// and assists clients at runtime to construct RPC requests and responses without precompiled service information
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot create listener")
	}

	log.Info().Msgf("starting gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot start grpc")
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot create server")
	}

	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOptions)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterNeoBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot register handler")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// server swagger docs
	// debug: https://stackoverflow.com/a/27946132
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot create statik file system")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	// strip to make the route relative to root dir
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot create l")
	}

	log.Printf("starting HTTP gateway server at %s", listener.Addr().String())
	httpHandler := gapi.HttpLogger(mux)
	err = http.Serve(listener, httpHandler)
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot start HTTP gateway")
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot start task processor")
	}
}

// run http server
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot create server")

	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot start server")
	}
}
