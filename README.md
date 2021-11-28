# linkinjs - A Go based program to find links from list of Js files quickly

# Installation 

go get -u github.com/rc4ne/linkinjs

# Efficient Usage with other tools like Gau

cat list_of_subdomains.txt | gau -b woff,png,jpeg,jpg -o urls.txt

cat urls.txt | grep "\\.js" > js_files.txt

linkinjs -n 50 -dl js_files.txt -o js_links.txt

# Sample Usage

![linkinjs_0](https://user-images.githubusercontent.com/83397936/143763281-3f0b68f8-3869-4ef1-a821-c2f5dd2d959d.JPG)

![linkinjs_1](https://user-images.githubusercontent.com/83397936/143763475-314f0f7e-ce00-4419-96b5-9ec496333fd3.JPG)


# Some Points to note

1. Using -m will match for the basedomain in the found links. Example: For http://xyz.tld/file.js, all the links matching keyword "xyz" will be considered. Default value-false.

![linkinjs_2](https://user-images.githubusercontent.com/83397936/143763460-9063eedc-4376-4d95-9eb4-091805971d1f.JPG)

![linkinjs_3](https://user-images.githubusercontent.com/83397936/143763465-ca4c91ee-b6d7-490e-aa30-bf5b5df3be74.JPG)

2. Flag -n is for concurrency, no of goroutines to use at once.

# Improvement

Lots of scope for improvement. I made this as side project during academics when trying Golang for first time. Most of inspiration from https://github.com/0xsha/GoLinkFinder. 

-> Concurrency can be implemented in better way

-> Implementing SecretFinder

