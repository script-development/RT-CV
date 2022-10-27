package matcher

import (
	"encoding/gob"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/apex/log"
	"github.com/script-development/RT-CV/helpers/random"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*

TODO: Implement caching

*/

func treeIDFileName(id *primitive.ObjectID) string {
	if id == nil {
		return "root.json"
	}
	return id.Hex() + ".json"
}

func cacheDir(createDirIfNotExsisting bool) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	cacheDir = path.Join(cacheDir, "rtcv/tree")

	if !createDirIfNotExsisting {
		return cacheDir, nil
	}

	err = os.MkdirAll(cacheDir, 0755)
	return cacheDir, err
}

// NukeCache well the name explains itself it clears the complete cache
func NukeCache() error {
	cacheDir, err := cacheDir(false)
	if err != nil {
		return nil
	}
	_, err = os.Stat(cacheDir)
	if err != nil {
		return nil
	}
	return os.RemoveAll(cacheDir)
}

// checkCacheDirSize caches the
func checkCacheDirSize(cacheDirPath string) error {
	entries, err := os.ReadDir(cacheDirPath)
	if err != nil {
		return err
	}
	if len(entries) > 50 {
		// Remove 4 random entries
		for i := 0; i < 4; i++ {
			entry := random.SliceIndex(entries)
			name := path.Join(cacheDirPath, entry.Name())
			err = os.RemoveAll(name)
			if err != nil {
				log.WithError(err).WithField("name", name).Warn("unable to remove cache entry")
				continue
			}
		}
	}

	return nil
}

// CacheTree caches a tree
func CacheTree(id *primitive.ObjectID, bytes []byte) error {
	cacheDirPath, err := cacheDir(true)
	if err != nil {
		return err
	}

	err = checkCacheDirSize(cacheDirPath)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(cacheDirPath, treeIDFileName(id)), bytes, 0644)
}

// ObtainCachedTree might return a cached tree
func ObtainCachedTree(id *primitive.ObjectID) (fs.FileInfo, *os.File, error) {
	cacheDirPath, err := cacheDir(false)
	if err != nil {
		return nil, nil, err
	}
	f, err := os.Open(path.Join(cacheDirPath, treeIDFileName(id)))
	if err != nil {
		return nil, nil, err
	}

	inf, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, nil, err
	}

	return inf, f, nil
}

// cacheSearch safes a list of optimized data for searching
func cacheSearch(data []FuzzySearchCacheEntry) error {
	cacheDirPath, err := cacheDir(true)
	if err != nil {
		return err
	}

	err = checkCacheDirSize(cacheDirPath)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path.Join(cacheDirPath, "search.gob"), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	for _, entry := range data {
		err = encoder.Encode(entry)
		if err != nil {
			return err
		}
	}

	return nil
}

// searchCache reads a optimized searching list
func searchCache(found func(FuzzySearchCacheEntry) (stop bool)) error {
	cacheDirPath, err := cacheDir(true)
	if err != nil {
		return err
	}

	f, err := os.Open(path.Join(cacheDirPath, "search.gob"))
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)
	first := true
	for {
		resp := FuzzySearchCacheEntry{}
		err = decoder.Decode(&resp)
		if err != nil {
			if !first && err.Error() == "gob: unknown type id or corrupted data" {
				// This error is thrown when we reached the end of the file / no bytes left over to read
				return nil
			}
			return err
		}
		if found(resp) {
			return nil
		}
		first = false
	}
}
