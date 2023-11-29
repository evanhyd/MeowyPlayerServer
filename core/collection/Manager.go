package collection

import (
	"errors"
	"io"
	"os"
	"sync"
)

var instance manager

type manager struct {
	mux sync.Map //id -> *sync.RWMutex, guard file read/write
}

func Initialize() error {
	return os.MkdirAll(CollectionPath(), 0777)
}

func getMux(id string) *sync.RWMutex {
	rwMux, _ := instance.mux.LoadOrStore(id, &sync.RWMutex{})
	return rwMux.(*sync.RWMutex)
}

/*
Must be called before using user provided filename to avoid malicious injection.
*/
func IsValidFileName(fileName string) bool {
	return fileNameValidator.MatchString(fileName)
}

func List() ([]CollectionInfo, error) {
	entries, err := os.ReadDir(CollectionPath())
	if err != nil {
		return nil, errors.New(err.Error() + " - collection path doesn't exist")
	}

	infos := make([]CollectionInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, errors.New(err.Error() + " - failed to read entry info: " + entry.Name())
		}
		infos = append(infos, CollectionInfo{info.Name(), info.ModTime(), info.Size()})
	}
	return infos, nil
}

func Update(src io.Reader, id string) error {
	mux := getMux(id)
	mux.Lock()
	defer mux.Unlock()

	file, err := os.OpenFile(CollectionFile(id), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		return errors.New(err.Error() + " - failed to update the collection: " + id)
	}
	defer file.Close()

	_, err = io.Copy(file, src)
	return err
}

func Fetch(dst io.Writer, id string) error {
	mux := getMux(id)
	mux.RLock()
	defer mux.RUnlock()

	file, err := os.Open(CollectionFile(id))
	if err != nil {
		return errors.New(err.Error() + " - failed to fetch the collection: " + id)
	}
	defer file.Close()

	_, err = io.Copy(dst, file)
	return err
}
