<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="https://user-images.githubusercontent.com/32605850/97817997-df110f80-1cf3-11eb-9fae-2db765d09563.png" alt="Project logo"></a>
</p>


<h3 align="center">GOP</h3>

<div align="center">

  [![Status](https://img.shields.io/badge/status-active-success.svg)]() 
  [![GitHub Issues](https://img.shields.io/github/issues/joshcarp/gop)](https://github.com/joshcarp/gop/issues)
  [![GitHub Pull Requests](https://img.shields.io/github/issues-pr/joshcarp/gop)](https://github.com/joshcarp/gop/pulls)
  [![License](https://img.shields.io/badge/license-apache2-blue.svg)](/LICENSE)

</div>

---


## 📝 Table of Contents
- [About](#about)
- [Getting Started](#getting_started)
- [Deployment](#deployment)
- [Usage](#usage)
- [Built Using](#built_using)
- [Authors](#authors)
- [Acknowledgments](#acknowledgement)

## 🧐 About <a name = "about"></a>
GOP: "Git Object Proxy" is a generic library used to implement moduling systems in programming languages. 
GOP defines two main interfaces that can be used to retrieve and cache resources from different sources. 
See [revision2.md](/design/revision2.md) to see design doc. 

## 🏁 Getting Started <a name = "getting_started"></a>

See [deployment](#deployment) for notes on how to deploy the project on a live system.

### Prerequisites
What things you need to install the software and how to install them.
- Go 1.13: currently google cloud functions only support upto the go 1.13 runtime

## 🔧 Running the tests <a name = "tests"></a>

`go test ./...`

## 🎈 Usage <a name="usage"></a>
`go run ./cmd/servegop`
- This will run a gop server on `localhost:8080`


## 🚀 Deployment <a name = "deployment"></a>

- See .github/workflows/cloud-function-deploy.yml

## ⛏️ Built Using <a name = "built_using"></a>
- [Google Cloud Functions](https://cloud.google.com/functions/) - Deployment
- [Google Cloud Storage](https://cloud.google.com/storage/) - Asset caching
- [Golang](https://golang.org/) - Server 

## ✍️ Authors <a name = "authors"></a>
- [@joshcarp](https://github.com/joshcarp)

## 🎉 Acknowledgements <a name = "acknowledgement"></a>
- Go Modules: Athens Project: https://github.com/gomods/athens 