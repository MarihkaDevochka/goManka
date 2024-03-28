package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/chimas/GoProject/config"
	"github.com/chimas/GoProject/db"
	_ "github.com/chimas/GoProject/docs"
	"github.com/chimas/GoProject/handler"
	"github.com/chimas/GoProject/middleware"
	"github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

//		@title			Manka Api
//		@version		1.0
//		@description	Manga search
//	 @BasePath	/
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		// log.Fatal("Error loading .env file")
	}

	router := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4000", "http://localhost:3000"},
	})

	db, err := db.DBConnection()
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
		return
	}
	defer db.Close()

	// opt, err := redis.ParseURL(config.LoadEnv().REDIS_URL)
	// if err != nil {
	// 	panic(err)
	// }
	// rdb := redis.NewClient(opt)
	opt, err := redis.ParseURL(config.LoadEnv().REDIS_URL)
	if err != nil {
		panic(err)
	}

	opt.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	rdb := redis.NewClient(opt)

	handlerM := handler.NewMangaHandler(db, rdb)
	handlerU := handler.NewUserHandler(db, rdb)
	router.HandleFunc("GET /yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})
	router.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	router.HandleFunc("GET /mangas", handlerM.Mangas)
	router.HandleFunc("GET /manga", handlerM.Manga)
	router.HandleFunc("GET /manga/{name}/{chapter}", handlerM.Chapter)
	router.HandleFunc("GET /popular", handlerM.Popular)
	router.HandleFunc("GET /filter", handlerM.Filter)
	router.HandleFunc("GET /user/{email}", handlerU.GetUser)
	router.HandleFunc("POST /user/create", handlerU.CreateUserIfNotExists)
	router.HandleFunc("POST /user/favorite/{name}/{email}", handlerU.ToggleFavorite)
	router.HandleFunc("GET /user/favorite/one", handlerU.IsUserFavorite)
	router.HandleFunc("GET /user/favorite/list", handlerU.UserFavList)
	router.HandleFunc("DELETE /user/delete", handlerU.DeleteUser)

	// router.HandleFunc("DELETE /user",handler)

	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "4000"
	}
	server := http.Server{
		Addr:    ":" + PORT,
		Handler: middleware.Logging(c.Handler(router)),
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
