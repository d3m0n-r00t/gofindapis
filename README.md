## GO FIND APIS

```
   _____  ____    ______ _____ _   _ _____             _____ _____  _____ 
  / ____|/ __ \  |  ____|_   _| \ | |  __ \      /\   |  __ \_   _|/ ____|
 | |  __| |  | | | |__    | | |  \| | |  | |    /  \  | |__) || | | (___  
 | | |_ | |  | | |  __|   | | | . ` | |  | |   / /\ \ |  ___/ | |  \___ \ 
 | |__| | |__| | | |     _| |_| |\  | |__| |  / ____ \| |    _| |_ ____) |
  \_____|\____/  |_|    |_____|_| \_|_____/  /_/    \_\_|   |_____|_____/ 
                                                                          
                                                               @demonroot     
```

This tool is build to find APIs. Just go recursively go around files and read them and scan for API keys with regex.
This tool is still in the developement phase. There is posibility of false positive. Please do contribute.

### Installation
 - To install ```gofindapis```...

```bash
curl https://raw.githubusercontent.com/d3m0n-r00t/gofindapis/master/install.sh > ~/install-gofindapis.sh && chmod +x ~/install-gofindapis.sh
~/install-gofindapis.sh
gofindapis <path_to_scan>
```

 - To add ```gofindapis``` as pre-commit hook.

 ```bash
 cd git-project
 curl https://raw.githubusercontent.com/d3m0n-r00t/gofindapis/master/pre-commit > .git/hooks/pre-commit
 ```

### The .goignore file
Any files can be avoided from being scanned by ```gofindapis``` by the use of a ```.goignore``` file. Just add the file in the git directory 
and the files in the file wont be scanned. For example to avoid scanning ```.git``` folder and its files add the following in ```.goignore```.
```
.git
```


#### TO DO
- Find a way to add this tool as git pre-commit hook check. --- **DONE**
- Optimization. --- **DONE**
- Build the script and make it a binary so that we can run it from anywhere. --- **DONE**
- Find a general way to give root path of the git project in ```pre-commit``` script in ```.git/hooks/pre-commit```. --- **DONE**
- Find more regex and reduce false positives.