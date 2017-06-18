package main

import (
  "fmt"
  "time"
  "io"
  "net/http"
  "html/template"
  "regexp"

  "./webui/layout"

  "github.com/gin-gonic/gin"
  "github.com/aviddiviner/gin-limit"
)

func compileHost(address, port string) (string) {
  return (address + ":" + port)
}

func (app *Application)InitiateWebUI() {
  //gin.SetMode(gin.ReleaseMode)
  r := gin.Default()
  r.Use(limit.MaxAllowed(20))

  html := template.Must(template.New("search_results").Parse(layout.DefaultLayout()))
  r.SetHTMLTemplate(html)

  //r.Use(rateLimit, gin.Recovery())
  searchWhiteList := "\\.(py|js|coffee|go|yaml|scss|css|html|c|cpp|m|h|java|txt|ini|yml|tmpl|tmp)$"
  searchBlackList := "/(node_modules|build|coverage)/"
  var searchableFiles [](*FileData)
  fmt.Println("Caching file contents...")
  var filePathIncludeRegexp, filePathExcludeRegexp *regexp.Regexp
  filePathIncludeRegexp, _ = regexp.Compile(searchWhiteList)
  filePathExcludeRegexp, _ = regexp.Compile(searchBlackList)
  searchableFiles, _ = GetFilepathsInDir(".", filePathIncludeRegexp, filePathExcludeRegexp)
  cacheAllFiles(searchableFiles)

  r.GET("/search", func(c *gin.Context) {
    query := c.Query("q")
    searchOptions := SearchOptions{FilePathInclude: "."}
    foundFiles, _ := SearchInDir(searchableFiles, "*", query, &searchOptions, 10, 10)
    c.HTML(200, "search_results", gin.H{
      "Query": query,
      "FoundCount": len(foundFiles),
      "FoundFiles": foundFiles,
    })
  })

  r.GET("/api/v1/search", func(c *gin.Context) {
    startTime := time.Now()
    // TODO: Always set to project root path
    query := c.Query("q")
    searchOptions := SearchOptions{FilePathInclude: ""}
    searchResults, _ := SearchInDir(searchableFiles, "*", query, &searchOptions, 10, 10)
    searchTime := time.Since(startTime)
    totalTime := time.Since(startTime)
    fmt.Printf("Found %v files (%v latency, %v total)\n", len(searchResults), searchTime, totalTime)
    c.JSON(http.StatusOK, searchResults)
    _ = cacheAllFiles(searchableFiles)
  })

  r.GET("/events", func (c *gin.Context) {
    c.Stream(func(w io.Writer) bool {
      time.Sleep(time.Second)
      c.SSEvent("ping", time.Now().String())
      return true
    })
  })

  // REST API
  //currentAPI := r.Group("/api/v" + config.State.APIVersion)
  //fileEndpoint := currentAPI.Group("/file")
  //{
  //  fileEndpoint.GET("/", file.List)
  //  fileEndpoint.GET("/:id", file.Get)
  //  fileEndpoint.GET("/:id/delete", file.Delete)
  //  fileEndpoint.POST("/", file.Post)
  //  fileEndpoint.POST("/:id/patch", file.Patch)
  //}
  host := compileHost(app.Config.WebUI.Host, app.Config.WebUI.Port)
  r.Run(host) // config.State.Port)
  fmt.Println("Serving WebUI on ", host)
}
