package main

import (
  "os"
  "github.com/codegangsta/cli"
)

func main() {
    app := cli.NewApp()
    app.Name = "alex"
    app.Usage = "Alexandria CMDB CLI"
    app.Version = "1.0.0"
  
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "u, url",
            Value: "http://localhost:3000",
            Usage: "Specify the API base URL",
            EnvVar: "ALEX_API_URL",
        },
        cli.BoolFlag{
            Name:  "verbose",
            Usage: "Show more output",
        },
    }
    
    // Commands
    app.Commands = []cli.Command{
    {
        Name: "build",
        Flags: []cli.Flag{
            cli.BoolFlag{
                Name:  "no-cache",
                Usage: "Do not use cache when building the image.",
            },
        },
        Usage:  "Build or rebuild services",
        Action: CmdBuild,
    }}

    app.Action = func(c *cli.Context) {
        println("Hello friend!")
    }

    app.Run(os.Args)
}

func CmdBuild(c *cli.Context) {
    println("Yep!")
}