package main

// How to write to spreadsheet:  https://stackoverflow.com/questions/39691100/golang-google-sheets-api-v4-write-update-example
import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"math/rand"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
        "gopkg.in/yaml.v2"
)

func main() {
        ctx := context.Background()

        type T struct {
                Credentials string
                Spreadsheetid string
        }

        t := T{}

        configdata, err := ioutil.ReadFile("config.yaml")
        if err != nil {
                log.Fatalf("Unable to read config.yaml file: %v", err)
        }
	//s := string(configdata)
	//fmt.Printf(s)
        err = yaml.Unmarshal(configdata, &t)
        if err != nil {
                log.Fatalf("error: %v", err)
        }

	//reader := bufio.NewReader(os.Stdin)
	//text, _ := reader.ReadString('\n')
	//fmt.Println(text)
	//text, _ = reader.ReadString('\n')
	//fmt.Println(text)
	//text, _ = reader.ReadString('\n')
	//fmt.Println(text)
	//text, _ = reader.ReadString('\n')
	//fmt.Println(text)

	homePath, err := expandHome(t.Credentials)
	absPath, err := filepath.Abs(homePath)
        if err != nil {
                log.Fatalf("error: %v", err)
        }
	b, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	//config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	config, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := config.Client(oauth2.NoContext)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	spreadsheetId := t.Spreadsheetid
	if spreadsheetId == "" {
		log.Fatalf("Missing 'spreadsheetid' from config file")
	}

	//os.Exit(3)

	// The A1 notation of a range to search for a logical table of data.
	// Values will be appended after the last row of the table.
	range2 := "A1"

	// How the input data should be interpreted.
	valueInputOption := "RAW"

	// How the input data should be inserted.
	insertDataOption := "OVERWRITE"

	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	var vr sheets.ValueRange
	myval := []interface{}{
		time.Now().UTC().Format(time.RFC3339),
		r.Float32(),
		r.Float32(),
	}
	vr.Values = append(vr.Values, myval)

	resp, err := srv.Spreadsheets.Values.Append(spreadsheetId, range2, &vr).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", resp)

	//  resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	//  if err != nil {
	//    log.Fatalf("Unable to retrieve data from sheet. %v", err)
	//  }
	//
	//  if len(resp.Values) > 0 {
	//    fmt.Println("Name, Major:")
	//    for _, row := range resp.Values {
	//      // Print columns A and E, which correspond to indices 0 and 4.
	//      fmt.Printf("%#v\n", row)
	//      fmt.Printf("%s, %s\n", row[0], row[1])
	//    }
	//  } else {
	//    fmt.Print("No data found.")
	//  }

}

func expandHome(path string) (string, error) {
    if len(path) == 0 || path[0] != '~' {
        return path, nil
    }

    usr, err := user.Current()
    if err != nil {
        return "", err
    }
    return filepath.Join(usr.HomeDir, path[1:]), nil
}

