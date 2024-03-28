package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func NewMangaHandler(db *sqlx.DB, rdb *redis.Client) *MangaHandler {
	return &MangaHandler{db: db, rdb: rdb}
}

type MangaHandler struct {
	db  *sqlx.DB
	rdb *redis.Client
}

type Manga struct {
	Name          string         `json:"name"`
	Img           string         `json:"img"`
	ImgHeader     string         `json:"imgHeader" db:"imgHeader"`
	Describe      string         `json:"describe"`
	Genres        pq.StringArray `json:"genres" db:"genres"`
	Author        string         `json:"author"`
	Country       string         `json:"country"`
	Published     int            `json:"published"`
	AverageRating float64        `json:"averageRating" db:"averageRating"`
	RatingCount   int            `json:"ratingCount" db:"ratingCount"`
	Status        string         `json:"status"`
	Popularity    int            `json:"popularity"`
	Id            int            `json:"id"`
	Chapters      []Chapter      `json:"chapters"`
}

type Chapter struct {
	Chapter   int            `json:"chapter"`
	Img       pq.StringArray `json:"img" db:"img"`
	Name      string         `json:"name"`
	AnimeName string         `json:"animeName" db:"animeName"`
	CreatedAt time.Time      `json:"createdAt" db:"createdAt"`
}

// @Summary Get all mangas
// @Description Retrieve a list of all mangas
// @Tags Manga
// @ID get-all-mangas
// @Accept  json
// @Produce  json
// @Success 200 {array} MangaSwag
// @Router /mangas [get]
func (m *MangaHandler) Mangas(w http.ResponseWriter, r *http.Request) {

	query := `SELECT * FROM "Anime"`
	var mangas []Manga
	err := m.db.Select(&mangas, query)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mangas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Get a manga by name
// @Description Retrieve a manga by its name
// @Tags Manga
// @ID get-manga-by-name
// @Accept  json
// @Produce  json
// @Param  name query string true "Name of the Manga"
// @Success 200 {object} MangaSwag
// @Router /manga [get]
func (m *MangaHandler) Manga(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	name := r.URL.Query().Get("name")

	val, err := m.rdb.Get(ctx, name).Result()
	if err == redis.Nil {

		query := `SELECT * FROM "Anime" WHERE name=$1`
		chaptersQuery := `SELECT * FROM "Chapter" WHERE "animeName" =$1`
		var manga Manga
		err := m.db.Get(&manga, query, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var chapters []Chapter
		err = m.db.Select(&chapters, chaptersQuery, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		manga.Chapters = chapters

		mangaJSON, err := json.Marshal(manga)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = m.rdb.Set(ctx, name, mangaJSON, time.Minute).Err()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(manga); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		manga := Manga{}
		err := json.Unmarshal([]byte(val), &manga)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(manga); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// @Summary Get a chapter
// @Description Find Manga Chapter
// @Tags Manga
// @ID get-chapter
// @Accept  json
// @Produce  json
// @Param  name path string true "Name of the Manga"
// @Param  chapter path string true "Chapter of the Manga"
// @Success 200 {object} ChapterSwag
// @Router /manga/{name}/{chapter} [get]
func (m *MangaHandler) Chapter(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	chapt := r.PathValue("chapter")

	var chapter Chapter

	query := `SELECT * FROM "Chapter" WHERE "animeName" =$1 AND chapter=$2`

	err := m.db.Get(&chapter, query, name, chapt)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chapter); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Get popular mangas
// @Description Retrieve a list of popular mangas
// @Tags Manga
// @ID get-popular-manga
// @Accept  json
// @Produce  json
// @Success 200 {array} MangaSwag
// @Router /popular [get]
func (m *MangaHandler) Popular(w http.ResponseWriter, r *http.Request) {

	query := `SELECT * FROM "Anime" ORDER BY "ratingCount" DESC LIMIT 14 `
	var animes []Manga
	err := m.db.Select(&animes, query)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(animes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (m *MangaHandler) Search(w http.ResponseWriter, r *http.Request) {

	query := `SELECT * FROM "Anime"`
	var animes []Manga
	err := m.db.Select(&animes, query)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(animes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type FilterParams struct {
	Name       string   `schema:"name"`
	Genres     []string `schema:"genres"`
	Status     string   `schema:"status"`
	Country    string   `schema:"country"`
	OrderField string   `schema:"orderField"`
	OrderSort  string   `schema:"orderSort"`
	Page       int      `schema:"page"`
	PerPage    int      `schema:"perPage"`
}

// @Summary Get a chapter
// @Description Find Manga Chapter
// @Tags Manga
// @ID Filter-anime
// @Accept  json
// @Produce  json
// @Param  name query string false "Name of the Manga"
// @Param  genres query []string false "Chapter of the Manga"
// @Param  status query string false "Name of the Manga"
// @Param  country query string false "Chapter of the Manga"
// @Param  orderField query string false "field of the Manga"
// @Param  orderSort query string false "sort of the Manga"
// @Param  page query int false "page not 0"
// @Param  perPage query int false "perPage"
// @Success 200 {array} MangaSwag
// @Router /filter [get]
func (m *MangaHandler) Filter(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	name := params.Get("name")
	genres := params["genres[]"]
	status := params.Get("status")
	country := params.Get("country")
	orderField := params.Get("orderField")
	orderSort := params.Get("orderSort")

	page, err := strconv.Atoi(params.Get("page"))
	if err != nil {
		log.Println("not have page")
	}
	perPage, err := strconv.Atoi(params.Get("perPage"))
	if err != nil {
		log.Println("not have perPage")
	}

	var mangas []Manga
	query := `SELECT * FROM "Anime"`
	args := []interface{}{}
	i := 1

	if name != "" || status != "" || country != "" || (len(genres) > 0 && genres[0] != "") {
		query += " WHERE"
	}

	if name != "" {
		name = "%" + name + "%"
		query += fmt.Sprintf(` "name" ILIKE $%d AND`, i)
		args = append(args, name)
		i++
	}
	if status != "" {
		query += fmt.Sprintf(` "status" = $%d AND`, i)
		args = append(args, status)
		i++
	}
	if country != "" {
		query += fmt.Sprintf(` "country" = $%d AND`, i)
		args = append(args, country)
		i++
	}
	if len(genres) > 0 && genres[0] != "" {
		for _, genre := range genres {
			query += fmt.Sprintf(` "genres" @> ARRAY[$%d] AND`, i)
			args = append(args, genre)
			i++
		}
	}

	query = strings.TrimSuffix(query, "AND")
	if orderField != "" && orderSort != "" {
		query += fmt.Sprintf(` ORDER BY "%s" %s`, orderField, orderSort)
	}
	if page > 0 && perPage > 0 {
		query += fmt.Sprintf(` LIMIT %d OFFSET %d`, perPage, (page-1)*perPage)
	}

	log.Println("q", query)
	err = m.db.Select(&mangas, query, args...)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mangas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
