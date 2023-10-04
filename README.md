# DocsPdf-go 
<!-- Idea -->
## Main Idea
- Pass the url of your page
- Pass the base path of your docs : (/docs or /documents)
- Get all the paths in form of three that are under the docs
- Make nested folders for each pdf like file system routing

## Run
```
go run main.go "https://vuejs.org/guide/introduction.html" "guide"
```

```
go run main.go "https://nextjs.org/docs" "docs"
```



## End 
https://vuejs.org/guide/introduction.html -> assets/vue.js/guide/introduction.html.pdf

