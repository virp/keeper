package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"keeper/cmd/client/command"
	"keeper/cmd/client/help"
	pb "keeper/gen/service"
	"keeper/internal/services"
	"keeper/pkg/logger"
)

var build = "develop"

type config struct {
	conf.Version
	Address string `conf:"default:localhost:3200,help:Server address"`
	Args    conf.Args
}

func main() {
	log, err := logger.New("KEEPER-CLIENT", "history.log")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("loop", "ERROR", err)
		//goland:noinspection GoUnhandledErrorResult
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// =========================================================================
	// Configuration

	cfg := config{
		Version: conf.Version{
			Build: build,
			Desc:  "GophKeeper Client",
		},
	}

	const prefix = "KEEPER"
	h, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(h)
			return nil
		}
		return fmt.Errorf("parse config: %w", err)
	}

	// =========================================================================
	// Check command
	if len(cfg.Args) == 2 && cfg.Args.Num(1) == "help" {
		switch cfg.Args.Num(0) {
		case "register":
			help.RegisterUsage()
		case "login":
			help.LoginUsage()
		case "ls":
			help.LsUsage()
		case "get":
			help.GetUsage()
		case "add":
			help.AddUsage()
		case "edit":
			help.EditUsage()
		case "del":
			help.DelUsage()
		default:
			help.Usage()
		}
		return nil
	}

	conn, err := grpc.Dial(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer conn.Close()
	grpcClient := pb.NewKeeperServiceClient(conn)
	client := services.NewClientService(grpcClient)
	cmd := command.NewCommand(log, client)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch cfg.Args.Num(0) {
	case "register":
		return cmd.Register(ctx)
	case "login":
		return cmd.Login(ctx)
	case "ls":
		return cmd.List(ctx)
	case "get":
		if len(cfg.Args) < 2 {
			help.GetUsage()
			return nil
		}
		return cmd.Get(ctx, cfg.Args.Num(1))
	case "add":
		if len(cfg.Args) < 3 {
			help.AddUsage()
			return nil
		}
		return cmd.Add(ctx, cfg.Args.Num(1), cfg.Args.Num(2))
	case "edit":
		if len(cfg.Args) < 2 {
			help.EditUsage()
			return nil
		}
		return cmd.Edit(ctx, cfg.Args.Num(1))
	case "del":
		if len(cfg.Args) < 2 {
			help.DelUsage()
			return nil
		}
		return cmd.Del(ctx, cfg.Args.Num(1))
	default:
		help.Usage()
		return nil
	}
}
