package handler

import "time"

type MangaSwag struct {
	Name          string        `json:"name"`
	Img           string        `json:"img"`
	ImgHeader     string        `json:"imgHeader" db:"imgHeader"`
	Describe      string        `json:"describe"`
	Genres        []string      `json:"genres" db:"genres"`
	Author        string        `json:"author"`
	Country       string        `json:"country"`
	Published     int           `json:"published"`
	AverageRating float64       `json:"averageRating" db:"averageRating"`
	RatingCount   int           `json:"ratingCount" db:"ratingCount"`
	Status        string        `json:"status"`
	Popularity    int           `json:"popularity"`
	Id            int           `json:"id"`
	Chapters      []ChapterSwag `json:"chapters"`
}

type ChapterSwag struct {
	Chapter   int       `json:"chapter"`
	Img       []string  `json:"genres" db:"img"`
	Name      string    `json:"name"`
	AnimeName string    `json:"animeName" db:"animeName"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
}

type UserSwag struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Favorite  []string  `json:"favorite"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
}
