package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/dimiro1/banner/autoload"
)

// Regex to compare lines to find API keys
// This was fucking JSON value which was not even valid. Made it into array so that it will be easier to process
var _regexes = [34]string{`AIza[0-9A-Za-z-_]{35}`, `AAAA[A-Za-z0-9_-]{7}:[A-Za-z0-9_-]{140}`, `6L[0-9A-Za-z-_]{38}|^6[0-9a-zA-Z_-]{39}$`, `ya29\\.[0-9A-Za-z\\-_]+`, `A[SK]IA[0-9A-Z]{16}`, `amzn\\.mws\\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`, `s3\\.amazonaws.com[/]+|[a-zA-Z0-9_-]*\\.s3\\.amazonaws.com`, `[a-zA-Z0-9-\\.\\_]+\\.s3\\.amazonaws\\.com`, `|s3://[a-zA-Z0-9-\\.\\_]+`, `|s3-[a-zA-Z0-9-\\.\\_\\/]+`, `|s3.amazonaws.com/[a-zA-Z0-9-\\.\\_]+`, `|s3.console.aws.amazon.com/s3/buckets/[a-zA-Z0-9-\\.\\_]+)`, `EAACEdEose0cBA[0-9A-Za-z]+`, `basic [a-zA-Z0-9=:_\\+\\/-]{5,100}`, `bearer [a-zA-Z0-9_\\-\\.=:_\\+\\/]{5,100}`, `api[key|_key|\\s+]+[a-zA-Z0-9_\\-]{5,100}`, `key-[0-9a-zA-Z]{32}`, `SK[0-9a-fA-F]{32}`, `AC[a-zA-Z0-9_\\-]{32}`, `AP[a-zA-Z0-9_\\-]{32}`, `access_token\\$production\\$[0-9a-z]{16}\\$[0-9a-f]{32}`, `sq0csp-[ 0-9A-Za-z\\-_]{43}|sq0[a-z]{3}-[0-9A-Za-z\\-_]{22,43}`, `sqOatp-[0-9A-Za-z\\-_]{22}|EAAA[a-zA-Z0-9]{60}`, `sk_live_[0-9a-zA-Z]{24}`, `rk_live_[0-9a-zA-Z]{24}`, `[a-zA-Z0-9_-]*:[a-zA-Z0-9_\\-]+@github\\.com*`, `-----BEGIN RSA PRIVATE KEY-----`, `-----BEGIN DSA PRIVATE KEY-----`, `-----BEGIN EC PRIVATE KEY-----`, `-----BEGIN PGP PRIVATE KEY BLOCK-----`, `ey[A-Za-z0-9-_=]+\\.[A-Za-z0-9-_=]+\\.?[A-Za-z0-9-_.+/=]*$`, `"api_token":"(xox[a-zA-Z]-[a-zA-Z0-9-]+)"`, `([-]+BEGIN [^\\s]+ PRIVATE KEY[-]+[\\s]*[^-]*[-]+END [^\\s]+ PRIVATE KEY[-]+)`, `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`}
var _keys = [34]string{"google_api", "firebase", "google_captcha", "google_oauth", "amazon_aws_access_key_id", "amazon_mws_auth_token", "amazon_aws_url", "amazon_aws_url2", "amazon_aws_url3", "amazon_aws_url4", "amazon_aws_url5", "amazon_aws_url6", "facebook_access_token", "authorization_basic", "authorization_bearer", "authorization_api", "mailgun_api_key", "twilio_api_key", "twilio_account_sid", "twilio_app_sid", "paypal_braintree_access_token", "square_oauth_secret", "square_access_token", "stripe_standard_api", "stripe_restricted_api", "github_access_token", "rsa_private_key", "ssh_dsa_private_key", "ssh_dc_private_key", "pgp_private_block", "json_web_token", "slack_token", "SSH_privKey", "Heroku API KEY"}

// Main function
func main() {
	// fmt.Println("\nHello World!!!")
	fmt.Printf("\n")
	DIR := os.Args[1]
	paths := getdir(DIR) // Get files names/paths and store in an array
	for _, file := range paths {
		if !checkignore(file) { // Directory check. Only reads paths if not a directory
			if !checkifdir(file) {
				readFile(file) // Returns the file content. Should we return this? Need to work on that
			}
		}
	}
}

// Recursivley go thourgh directories and find all the files and return an array of files/paths
func getdir(path string) []string {
	var arr []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		arr = append(arr, path)
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return arr
}

// Read file and store the data in data probably the main function
func readFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Println("[-]Error!!!!")
	}
	// fmt.Println(_regexes[1])
	// fmt.Println(_keys[1])
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text() // the content is read line by line
		fmt.Println(line)
		break
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// This function checks if the path is a directory or not
func checkifdir(file string) bool {
	fi, err := os.Stat(file)
	if err != nil {
		log.Println(err)
		return true
	}
	mode := fi.Mode()
	return mode.IsDir() //return true if file is directory and false if not
}

func checkignore(file string) bool {
	ignorefile := ".goignore"
	ignore, err := os.Open(ignorefile)
	if err != nil {
		log.Println(err)
		return true
	}
	defer ignore.Close()
	scanner := bufio.NewScanner(ignore)
	for scanner.Scan() {
		line := scanner.Text()
		if file == line {
			return true
		}
	}
	return false
}
