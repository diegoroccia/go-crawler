package main

import (
	"channels/crawler"
	"channels/fetcher"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	r := gin.Default()

	// Define a struct to represent a job
	type Job struct {
		ID     string
		Status string
		Result []string
	}

	jobStatus := make(map[string]*Job)

	r.GET("/crawl", func(c *gin.Context) {

		url := c.Query("url")

		jobID := uuid.New().String()

		job := &Job{
			ID:     jobID,
			Status: "Pending",
		}

		go func() {

			fetcher := fetcher.NewFetcher()

			crawler := crawler.Crawler{
				Url:     "https://" + url + "/",
				Fetcher: fetcher,
			}

			job.Status = "Running"
			jobStatus[jobID] = job

			urls := crawler.Crawl()

			job.Status = "Finished"
			job.Result = urls

			jobStatus[jobID] = job
		}()

		c.JSON(200, gin.H{
			"job_id": jobID,
		})

	})

	r.GET("/job/:jobID", func(c *gin.Context) {
		jobID := c.Param("jobID")

		// Check if the job ID exists in the map
		if job, ok := jobStatus[jobID]; ok {
			c.JSON(200, gin.H{
				"job_id": job.ID,
				"status": job.Status,
				"result": job.Result,
			})
		} else {
			c.JSON(404, gin.H{
				"error": "Job not found",
			})
		}
	})

	r.Run()

}
