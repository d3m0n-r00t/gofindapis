package main

// All imports
import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"

	_ "github.com/dimiro1/banner/autoload"
)

// Regex to compare lines to find API keys
// This was JSON value which was not even valid. Made it into array so that it will be easier to process
var _regexes = [33]string{`AIza[0-9A-Za-z-_]{35}`, `AAAA[A-Za-z0-9_-]{7}:[A-Za-z0-9_-]{140}`, `6L[0-9A-Za-z-_]{38}|^6[0-9a-zA-Z_-]{39}$`, `ya29\\.[0-9A-Za-z\\-_]+`, `A[SK]IA[0-9A-Z]{16}`, `amzn\\.mws\\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`, `s3\\.amazonaws.com[/]+|[a-zA-Z0-9_-]*\\.s3\\.amazonaws.com`, `[a-zA-Z0-9-\\.\\_]+\\.s3\\.amazonaws\\.com`, `|s3://[a-zA-Z0-9-\\.\\_]+`, `|s3-[a-zA-Z0-9-\\.\\_\\/]+`, `|s3.amazonaws.com/[a-zA-Z0-9-\\.\\_]+`, `EAACEdEose0cBA[0-9A-Za-z]+`, `basic [a-zA-Z0-9=:_\\+\\/-]{5,100}`, `bearer [a-zA-Z0-9_\\-\\.=:_\\+\\/]{5,100}`, `api[key|_key|\\s+]+[a-zA-Z0-9_\\-]{5,100}`, `key-[0-9a-zA-Z]{32}`, `SK[0-9a-fA-F]{32}`, `AC[a-zA-Z0-9_\\-]{32}`, `AP[a-zA-Z0-9_\\-]{32}`, `access_token\\$production\\$[0-9a-z]{16}\\$[0-9a-f]{32}`, `sq0csp-[ 0-9A-Za-z\\-_]{43}|sq0[a-z]{3}-[0-9A-Za-z\\-_]{22,43}`, `sqOatp-[0-9A-Za-z\\-_]{22}|EAAA[a-zA-Z0-9]{60}`, `sk_live_[0-9a-zA-Z]{24}`, `rk_live_[0-9a-zA-Z]{24}`, `[a-zA-Z0-9_-]*:[a-zA-Z0-9_\\-]+@github\\.com*`, `-----BEGIN RSA PRIVATE KEY-----`, `-----BEGIN DSA PRIVATE KEY-----`, `-----BEGIN EC PRIVATE KEY-----`, `-----BEGIN PGP PRIVATE KEY BLOCK-----`, `ey[A-Za-z0-9-_=]+\\.[A-Za-z0-9-_=]+\\.?[A-Za-z0-9-_.+/=]*$`, `"api_token":"(xox[a-zA-Z]-[a-zA-Z0-9-]+)"`, `([-]+BEGIN [^\\s]+ PRIVATE KEY[-]+[\\s]*[^-]*[-]+END [^\\s]+ PRIVATE KEY[-]+)`, `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`}
var _keys = [33]string{"google_api", "firebase", "google_captcha", "google_oauth", "amazon_aws_access_key_id", "amazon_mws_auth_token", "amazon_aws_url", "amazon_aws_url2", "amazon_aws_url3", "amazon_aws_url4", "amazon_aws_url5", "facebook_access_token", "authorization_basic", "authorization_bearer", "authorization_api", "mailgun_api_key", "twilio_api_key", "twilio_account_sid", "twilio_app_sid", "paypal_braintree_access_token", "square_oauth_secret", "square_access_token", "stripe_standard_api", "stripe_restricted_api", "github_access_token", "rsa_private_key", "ssh_dsa_private_key", "ssh_dc_private_key", "pgp_private_block", "json_web_token", "slack_token", "SSH_privKey", "Heroku API KEY"}

// var wg sync.WaitGroup

var dirarr []string

// Main function
func main() {
	// fmt.Println("\nHello World!!!") Well here is where I started Go programming
	fmt.Printf("\n\n")
	DIR := os.Args[1]
	paths := getdir(DIR)     // Get files names/paths and store in an array
	c := make(chan []string) // Create a channel for concurrency
	ignorefile := ".goignore"
	ignore, err := os.Open(ignorefile)
	checkerror(err)
	defer ignore.Close()
	scanner := bufio.NewScanner(ignore)
	for scanner.Scan() {
		line := scanner.Text()
		if checkifdir(line) {
			filepath.WalkDir(line, walk)
		} else {
			line = DIR + "/" + line
			dirarr = append(dirarr, line)
		}
	}
	for _, file := range paths {
		if !checkignore(dirarr, file) { // Ignore check. If you don't want to scan any file in code base just add it in .goignore. Remember to add file name
			if !checkifdir(file) { // Directory check. Only reads paths if not a directory
				// wg.Add(1)
				go doMagic(c, file) // Returns the file content. Should we return this? Need to work on that
			}
		}
	}
	<-c // Wait for goroutine to finish scanning all files before iterating over the channels
	for item := range c {
		if len(item) > 0 {
			os.Exit(1) //Exit with non-zero code just to add the script as pre-commit hook
		}
	}
	fmt.Println("[+] Done No keys found!!!")
	// wg.Wait()
}

// Ignore files callback to return full path entries on walkdir
func walk(file string, d fs.DirEntry, e error) error {
	if e != nil {
		return e
	}
	if !d.IsDir() {
		dirarr = append(dirarr, file)
	}
	return nil
}

// Recursivley go thourgh directories and find all the files and return an array of files/paths
func getdir(path string) []string {
	var arr []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		arr = append(arr, path) // Make an array of paths from
		return nil
	})
	checkerror(err)
	return arr
}

// Read file and store the data in data probably the main function
func doMagic(c chan<- []string, file string) {
	// defer wg.Done()
	f, err := os.Open(file)
	var new_arr []string
	checkerror(err)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		line := scanner.Text() // the content is read line by line
		for i := range _regexes {
			r, err := regexp.Compile(_regexes[i])
			checkerror(err)
			found := r.Find([]byte(line)) // Find keys with regex
			if len(found) > 0 {
				fmt.Println("[+] Found ", _keys[i], " ", string(found), "in file", file) // Final Result!! Prints result if any key found
				new_arr = append(new_arr, string(found))
			}
		}
	}
	c <- new_arr
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// This function checks if the path is a directory or not
func checkifdir(file string) bool {
	fi, err := os.Stat(file)
	checkerror(err)
	mode := fi.Mode()
	return mode.IsDir() //return true if file is directory and false if not
}

// If there are some files we need to ignore and not scan add it to .goignore file
func checkignore(arr []string, file string) bool {
	for line := range arr {
		if file == arr[line] {
			return true
		}
	}
	return false // Return false if not in ignore file
}

// A common error function
func checkerror(e error) {
	if e != nil {
		log.Println("[-]Error", e)
	}
}
