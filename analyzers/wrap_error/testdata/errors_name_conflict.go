package testdata

import ( // want `adding "github.com/pkg/errors" import`
	"fmt"
	"log"
	"os"
	"time"

	"./errors"
)

func readFileErrorsNameConflict(name string) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, errors.ErrOpenFile // want `unwrapped error found 'return nil, errors.ErrOpenFile'`
	}

	var body []byte
	if _, err = file.Read(body); err != nil {
		return nil, errors.ErrReadFile // want `unwrapped error found 'return nil, errors.ErrReadFile'`
	}

	if time.Now().UTC().Weekday() == time.Friday {
		return nil, errors.ErrFriday
	}

	return body, nil
}

func mainErrorsNameConflict() {
	body, err := readFileErrorsNameConflict(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(body)
}
