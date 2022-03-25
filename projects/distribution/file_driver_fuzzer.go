package filesystem

import (
	"context"
	"os"
	storagedriver "github.com/distribution/distribution/v3/registry/storage/driver"
	fuzz "github.com/AdaLogics/go-fuzz-headers"
)

func FuzzFilesystemDriver(data []byte) int {
	f := fuzz.NewConsumer(data)
	err := os.Mkdir("fuzz-dir", 0755)
	if err != nil {
		return 0
	}
	defer os.RemoveAll("fuzz-dir")
	params := map[string]interface{}{
		"maxthreads": 1,
		"rootdirectory": "fuzz-dir",
	}
	driver, err := FromParameters(params)
	if err != nil {
		return 0
	}

	noOfOps, err := f.GetInt()
	if err != nil {
		return 0
	}
	for i:=0;i<noOfOps%10;i++ {
		opType, err := f.GetInt()
		if err != nil {
			return 0
		}
		switch opType%10 {
		case 0:
			err := f.CreateFiles("fuzz-dir")
			if err != nil {
				return 0
			}
		case 1:
			path, err := f.GetString()
			if err != nil {
				return 0
			}
			_, _ = driver.GetContent(context.Background(), path)
		case 2:
			subPath, err := f.GetString()
			if err != nil {
				return 0
			}
			contents, err := f.GetBytes()
			if err != nil {
				return 0
			}
			_ = driver.PutContent(context.Background(), subPath, contents)
		case 3:
			path, err := f.GetString()
			if err != nil {
				return 0
			}
			offset, err := f.GetInt()
			if err != nil {
				return 0
			}
			reader, err := driver.Reader(context.Background(), path, int64(offset))
			if err == nil {
				defer reader.Close()
			}
		case 4:
			subPath, err := f.GetString()
			if err != nil {
				return 0
			}
			_, _ = driver.Stat(context.Background(), subPath)
		case 5:
			subPath, err := f.GetString()
			if err != nil {
				return 0
			}
			_, _ = driver.List(context.Background(), subPath)
		case 6:
			sourcePath, err := f.GetString()
			if err != nil {
				return 0
			}
			destPath, err := f.GetString()
			if err != nil {
				return 0
			}
			_ = driver.Move(context.Background(), sourcePath, destPath)
		case 7:
			subPath, err := f.GetString()
			if err != nil {
				return 0
			}
			_ = driver.Delete(context.Background(), subPath)
		case 8:
			path, err := f.GetString()
			if err != nil {
				return 0
			}
			err = driver.Walk(context.Background(), path, func(fileInfo storagedriver.FileInfo) error {
				return nil
			})
		case 9:
			subPath, err := f.GetString()
			if err != nil {
				return 0
			}
			append, err := f.GetBool()
			if err != nil {
				return 0
			}
			fw, err := driver.Writer(context.Background(), subPath, append)
			if err != nil {
				return 0
			}
			defer fw.Close()
			p, err := f.GetBytes()
			if err != nil {
				return 0
			}
			_, err = fw.Write(p)
			if err != nil {
				return 0
			}
			_ = fw.Commit()
		}
	}
	return 1
}