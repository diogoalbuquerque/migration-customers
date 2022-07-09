package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/diogoalbuquerque/migration-customers/config"
	"github.com/diogoalbuquerque/migration-customers/internal"
	"github.com/diogoalbuquerque/migration-customers/pkg/secret"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func syncPeopleMigrationHandle(ctx context.Context) error {
	c, err := config.NewConfig()

	if err != nil {
		log.Println(err)
		return err
	}

	s, err := secret.New(ctx, *c).Load()

	if err != nil {
		log.Println(err)
		return err
	}

	return internal.Run(ctx, *c, *s)
}

func main() {
	lambda.Start(syncPeopleMigrationHandle)
}
