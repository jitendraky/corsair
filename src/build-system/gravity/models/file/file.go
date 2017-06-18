package file

import (
	"net/http"
	"strconv"
	"time"

	"../../hateoas"

	"github.com/gin-gonic/gin"
)

// Bucket is the name of the bucket storing all the files
const (
	Bucket = "files"
	Type   = "file"
)

// File is the main struct
type File struct {
	ID                 int       `json:"id"`
	Name       string    `storm:"index" json:"product_title"`
	Description string    `json:"product_description"`
	CreatedAt          time.Time `storm:"index"`
	ModifiedAt          time.Time `storm:"index"`
	AccessedAt          time.Time `storm:"index"`
}

// Validate validates that all the required files are not empty.
func (file File) Validate() hateoas.Errors {
	var errors hateoas.Errors
	if file.Name == "" {
		errors = append(errors, hateoas.Error{
			Status: http.StatusBadRequest,
			Title:  "name field is required",
		})
	}
	return errors
}

// Data contains the Type of the request and the Attributes
type Data struct {
	Type       string `json:"type,omitempty"`
	Attributes *File `json:"attributes,omitempty"`
	Links      *Links `json:"links,omitempty"`
}

// Links represent a list of links
type Links map[string]string

// Wrapper is the HATEOAS wrapper
type Wrapper struct {
	Data   *Data           `json:"data,omitempty"`
	Errors *hateoas.Errors `json:"errors,omitempty"`
}

// MultiWrapper is a wrapper that can accept multiple Data
type MultiWrapper struct {
	Data   *[]Data         `json:"data,omitempty"`
	Errors *hateoas.Errors `json:"errors,omitempty"`
}

// Post is the handler to POST a new File 
func Post(c *gin.Context) {
	var err error
	var json = Wrapper{}
	if err = c.BindJSON(&json); err == nil {
		errors := json.Data.Attributes.Validate()
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, Wrapper{Errors: &errors})
			return
		}
		var file *File
		file = json.Data.Attributes
		if err = file.Save(); err != nil {
			json.Data = nil
			json.Errors = &hateoas.Errors{hateoas.Error{Status: http.StatusInternalServerError, Title: "could not save file"}}
			c.JSON(http.StatusInternalServerError, json)
		} else {
			json.Data.Links = &Links{"self": c.Request.URL.RequestURI() + strconv.Itoa(json.Data.Attributes.ID)}
			c.JSON(http.StatusCreated, json)
		}
	} else {
		json.Data = nil
		json.Errors = &hateoas.Errors{hateoas.Error{Status: http.StatusInternalServerError, Title: "Bad json format"}}
		c.JSON(http.StatusBadRequest, json)
	}
}

func List(c *gin.Context) {
	var json = MultiWrapper{}
	var datas = []Data{}
	files, err := All()
	if err != nil {
		json.Errors = &hateoas.Errors{hateoas.Error{Status: http.StatusInternalServerError, Title: "could not retrieve files"}}
		c.JSON(http.StatusInternalServerError, json)
		return
	}
	for index := range files {
		datas = append(datas, Data{Type: Type, Attributes: &files[index]})
	}
	json.Data = &datas
	c.JSON(http.StatusOK, json)
}

func Get(c *gin.Context) {
	var err error
	var file File
	var json = Wrapper{}

	id, _ := strconv.Atoi(c.Param("id"))
	if file, err = file.Get(id); err != nil {
		json.Errors = &hateoas.Errors{hateoas.Error{Status: http.StatusNotFound, Title: "id could not be found"}}
		c.JSON(http.StatusNotFound, json)
		return
	}
	json.Data = &Data{Type: Type, Attributes: &file}
	c.JSON(http.StatusOK, json)
}

func Patch(c *gin.Context) {
	var err error
	var file File
	var json = Wrapper{}
	id, _ := strconv.Atoi(c.Param("id"))
	if file, err = file.Get(id); err != nil {
		json.Errors = &hateoas.Errors{hateoas.Error{Status: http.StatusNotFound, Title: "id could not be found"}}
		c.JSON(http.StatusNotFound, json)
		return
	}
	json.Data = &Data{Type: Type, Attributes: &file}
	if err = c.BindJSON(&json); err == nil {
		if err = json.Data.Attributes.Save(); err != nil {
			json.Data = nil
			json.Errors = &hateoas.Errors{hateoas.Error{Status: http.StatusInternalServerError, Title: "could not save file"}}
			c.JSON(http.StatusInternalServerError, json)
		} else {
			c.JSON(http.StatusCreated, json)
		}
	} else {
		json.Data = nil
		json.Errors = &hateoas.Errors{hateoas.Error{Status: http.StatusInternalServerError, Title: "Bad json format"}}
		c.JSON(http.StatusBadRequest, json)
	}
}

func Delete(c *gin.Context) {
	var err error
	var file File
	var json = Wrapper{}
	id, _ := strconv.Atoi(c.Param("id"))
	if file, err = file.Get(id); err != nil {
		json.Errors = &hateoas.Errors{hateoas.Error{Status: http.StatusNotFound, Title: "id could not be found"}}
		c.JSON(http.StatusNotFound, json)
		return
	}
	if err = file.Delete(); err != nil {
		json.Errors = &hateoas.Errors{hateoas.Error{Status: http.StatusInternalServerError, Title: "couldn't delete resource"}}
		c.JSON(http.StatusInternalServerError, json)
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}
