package testdata

import ( // want `adding "github.com/pkg/errors" import`
	"fmt"
	"log"
	"os"
	"time"
)

var ErrOpenFile = fmt.Errorf("error opening file")
var ErrReadFile = fmt.Errorf("error reading file")
var ErrFriday = fmt.Errorf("error - friday's are for fun, not files!")

func readFileBasicExample(name string) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, ErrOpenFile // want `unwrapped error found 'return nil, ErrOpenFile'`
	}

	var body []byte
	if _, err = file.Read(body); err != nil {
		return nil, ErrReadFile // want `unwrapped error found 'return nil, ErrReadFile'`
	}

	if time.Now().UTC().Weekday() == time.Friday {
		return nil, ErrFriday
	}

	return body, nil
}

func mainBasicExample() {
	body, err := readFileBasicExample(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(body)
}
