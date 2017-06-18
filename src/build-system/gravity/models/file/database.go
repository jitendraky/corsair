package file

import (
	"fmt"

	"../../database"
)

func (file *File) Save() error {
	if !database.Storage.Opened {
		return fmt.Errorf("Database Error: DB must be opened before deleting.")
	}
	return database.Storage.DB.Save(file)
}

func (file File) Delete() error {
	if !database.Storage.Opened {
		return fmt.Errorf("Database Error: DB must be opened before deleting.")
	}
	return database.Storage.DB.Remove(&file)
}

func (file File) Get(key int) (File, error) {
	if !database.Storage.Opened {
		return file, fmt.Errorf("Database Error: DB must be opened before deleting.")
	}
	err := database.Storage.DB.One("ID", key, &file)
	return file, err
}

// All returns all the files
func All() ([]File, error) {
	var err error
	var files []File
	if !database.Storage.Opened {
		return files, fmt.Errorf("Database must be opened first.")
	}
	database.Storage.DB.All(&files)
	return files, err
}
