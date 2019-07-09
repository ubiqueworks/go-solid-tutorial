package main

import (
	"context"
	"log"
	"os"

	"github.com/ubiqueworks/go-solid-tutorial/pb"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "create new user",
			Action: createUser,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "addr",
					Value: "0.0.0.0:9000",
					Usage: "Server RPC address",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "Name of the user to create",
				},
				cli.StringFlag{
					Name:  "password",
					Usage: "Password of the user to create",
				},
			},
		},
		{
			Name:   "delete",
			Usage:  "delete existing user",
			Action: deleteUser,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "addr",
					Value: "0.0.0.0:9000",
					Usage: "Server RPC address",
				},
				cli.StringFlag{
					Name:  "id",
					Usage: "ID of the user to delete",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func createUser(c *cli.Context) error {
	conn, err := grpc.Dial(c.String("addr"), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewUserAdminClient(conn)
	reply, err := client.CreateUser(context.Background(), &pb.CreateUserRequest{
		Name: c.String("name"),
		Password: c.String("password"),
	})

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	log.Printf(`User "%s" created with ID: %s`, reply.User.Name, reply.User.Id)
	return nil
}

func deleteUser(c *cli.Context) error {
	conn, err := grpc.Dial(c.String("addr"), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userID := c.String("id")

	client := pb.NewUserAdminClient(conn)
	_, err = client.DeleteUser(context.Background(), &pb.DeleteUserRequest{
		Id: userID,
	})
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	log.Printf(`User with ID "%s" deleted`, userID)
	return nil
}
