package cwl

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto"
)

// SecondaryFile represents an element of "secondaryFiles".
type SecondaryFile struct {
	Entry string
}

// NewList constructs list of "SecondaryFile".
func (s SecondaryFile) NewList(i interface{}) []SecondaryFile {
	var dest []SecondaryFile
	switch x := i.(type) {
	case []interface{}:
		for _, v := range x {
			dest = append(dest, SecondaryFile{Entry: v.(string)})
		}
	}
	return dest
}

// ToFiles resolves the SecondaryFile to a concrete file or files with locations
// set. The supplied vm must have 'self' set to the details of the parent file.
// dir is the main configured output directory. The file must exist on disk.
// Setting the copy boolean to true will copy secondary files that aren't in the
// same dir as parentPath to that dir and alter the location accordingly.
func (s SecondaryFile) ToFiles(dir, parentPath string, vm *otto.Otto, copy bool) ([]interface{}, error) {
	str, _, obj, err := evaluateExpression(s.Entry, vm)
	if err != nil {
		return nil, err
	}

	var files []interface{}
	if obj != nil {
		files, err = ottoObjToFiles(obj, dir)
		if err != nil {
			return nil, err
		}
	} else if str != "" {
		if str == s.Entry {
			// wasn't an expression, it's a suffix for the
			// main file's basename
			basename := filepath.Base(parentPath)
			if strings.HasPrefix(str, "^") {
				for {
					str = strings.TrimPrefix(str, "^")
					basename = strings.TrimSuffix(basename, filepath.Ext(basename))

					if !strings.HasPrefix(str, "^") {
						break
					}
				}
			}
			str = basename + str
		}

		sFile := make(map[interface{}]interface{})
		sFile[fieldClass] = typeFile
		sFile[fieldLocation] = str
		files = append(files, sFile)
	}

	// check each file exists and fill in size and checksum
	for _, f := range files {
		file := f.(map[interface{}]interface{})
		thisPath := filepath.Join(dir, file[fieldLocation].(string))

		f, err := os.Stat(thisPath)
		if err != nil {
			return nil, err
		}

		if f.IsDir() {
			file[fieldClass] = typeDirectory
			listing, errd := dirToListing(thisPath)
			if errd != nil {
				return nil, errd
			}
			file[fieldListing] = listing
		} else {
			size, sha, errf := fileSizeAndSha(thisPath)
			if errf != nil {
				return nil, errf
			}
			file[fieldSize] = size
			file[fieldChecksum] = sha
		}

		parentDir := filepath.Dir(parentPath)
		if !strings.HasPrefix(thisPath, parentDir) && copy {
			// copy the file to our output dir *** doesn't support copying dirs...
			newLocation := filepath.Base(thisPath)
			newPath := filepath.Join(parentDir, newLocation)
			err = copyFile(thisPath, newPath)
			if err != nil {
				return nil, err
			}
			file[fieldLocation] = newLocation
		}
	}

	return files, nil
}

// dirToListing reads the contents of the given directory and returns a slice of
// file or dirs.
func dirToListing(path string) ([]interface{}, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = dir.Close()
	}()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var listing []interface{}
	for _, file := range files {
		abs := filepath.Join(path, file.Name())
		if file.IsDir() {
			recurse, errr := dirToListing(abs)
			if errr != nil {
				return nil, errr
			}

			sDir := make(map[interface{}]interface{})
			sDir[fieldClass] = typeDirectory
			sDir[fieldLocation] = file.Name()
			sDir[fieldListing] = recurse
			listing = append(listing, sDir)
		} else {
			size, sha, errf := fileSizeAndSha(abs)
			if errf != nil {
				return nil, errf
			}

			sFile := make(map[interface{}]interface{})
			sFile[fieldClass] = typeFile
			sFile[fieldLocation] = file.Name()
			sFile[fieldBasename] = file.Name()
			sFile[fieldSize] = size
			sFile[fieldChecksum] = sha
			listing = append(listing, sFile)
		}
	}
	return listing, err
}

// SecondaryFiles is a slice of SecondaryFile.
type SecondaryFiles []SecondaryFile

// ToFiles calls ToFiles() on each constituent SecondaryFile and returns the
// combined set of files. It also sets up self in the vm for you.
func (sfs SecondaryFiles) ToFiles(dir, parentPath string, vm *otto.Otto, copy bool) ([]interface{}, error) {
	rel, err := filepath.Rel(dir, parentPath)
	if err != nil {
		return nil, err
	}

	parentFile := map[string]interface{}{
		fieldClass:    typeFile,
		fieldLocation: rel,
	}

	err = vm.Set("self", fileToSelf(parentPath, parentFile))
	if err != nil {
		return nil, err
	}

	var thisSlice []interface{}
	for _, sf := range sfs {
		theseFiles, err := sf.ToFiles(dir, parentPath, vm, copy)
		if err != nil {
			return nil, err
		}
		thisSlice = append(thisSlice, theseFiles...)
	}

	return thisSlice, nil
}
