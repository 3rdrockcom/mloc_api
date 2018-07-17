package cmd

import (
	"github.com/epointpayment/mloc_api_go/app/config"
	"github.com/epointpayment/mloc_api_go/app/controllers"
	"github.com/epointpayment/mloc_api_go/app/database"
	"github.com/epointpayment/mloc_api_go/app/log"
	"github.com/epointpayment/mloc_api_go/app/router"
	"github.com/epointpayment/mloc_api_go/app/services"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

// serveCmd executes webserver
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start server",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// initiate logging
		log.Start()
		defer log.Stop()

		// load config
		cfg, err := config.New()
		if err != nil {
			log.Fatal(err)
		}
		log.SetMode(cfg.Environment)

		// Create new connection handler for database
		db := database.NewDatabase(cfg.DB.Driver, cfg.DB.DSN)
		err = db.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Setup services
		err = services.New(db)
		if err != nil {
			log.Fatal(err)
		}

		// Setup router and run
		c := controllers.NewControllers(db)
		r := router.NewRouter(c)
		err = r.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}
