package matcher

import (
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

func idName(id *primitive.ObjectID) string {
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

// CacheTree caches a tree
func CacheTree(id *primitive.ObjectID, bytes []byte) error {
	cacheDirPath, err := cacheDir(true)
	if err != nil {
		return err
	}

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

	return ioutil.WriteFile(path.Join(cacheDirPath, idName(id)), bytes, 0644)
}

// ObtainCachedTree might return a cached tree
func ObtainCachedTree(id *primitive.ObjectID) (fs.FileInfo, *os.File, error) {
	cacheDirPath, err := cacheDir(false)
	if err != nil {
		return nil, nil, err
	}
	f, err := os.Open(path.Join(cacheDirPath, idName(id)))
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
