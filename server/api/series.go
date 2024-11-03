package api

import (
	"github.com/AniComix/query"
	"github.com/AniComix/server/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CreateSeries(c *gin.Context) {
	_, err := getCurrentUser(c)
	if err != nil {
		badRequest(c, "invalid user")
		return
	}
	series := &models.Series{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Seasons:     0,
	}
	err = query.Series.Create(series)
	if err != nil {
		internalServerError(c)
		return
	}
	c.JSON(200, gin.H{"message": series.ID})
}

func DeleteSeries(c *gin.Context) {
	_, err := getCurrentUser(c)
	if err != nil {
		badRequest(c, "invalid user")
		return
	}
	id, err := strconv.Atoi(c.PostForm("series_id"))
	if err != nil {
		badRequest(c, "invalid series id")
		return
	}
	s, err := query.Series.Where(query.Series.ID.Eq(uint(id))).First()
	if err != nil {
		badRequest(c, "series not found")
		return
	}
	_, err = query.Series.Delete(s)
	if err != nil {
		internalServerError(c)
	}
	c.JSON(200, gin.H{"message": "series deleted"})
}

func SearchSeries(c *gin.Context) {
	_, err := getCurrentUser(c)
	if err != nil {
		badRequest(c, "invalid user")
		return
	}
	title := c.PostForm("title")
	series, err := query.Series.Where(query.Series.Title.Like(title)).Find()
	if err != nil {
		internalServerError(c)
		return
	}
	c.JSON(200, gin.H{"message": series})
}
