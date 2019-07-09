package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ubiqueworks/go-interface-usage/pb"
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
		{
			Name:   "list",
			Usage:  "list all users",
			Action: listUsers,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "addr",
					Value: "0.0.0.0:9000",
					Usage: "Server RPC address",
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
	})

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	log.Printf(`User "%s" created with ID: %s`, reply.User.Name, reply.User.Id)
	return nil
}

func deleteUser(c *cli.Context) error {
	fmt.Println(c.String("addr"))
	fmt.Println(c.String("id"))

	// conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer conn.Close()
	//
	// cli := pb.NewUserAdminClient(conn)

	return nil
}

func listUsers(c *cli.Context) error {
	fmt.Println(c.String("addr"))

	return nil
}
