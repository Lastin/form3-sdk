##  Form3 SDK
#### Instructions

- `make test` starts stack of docker containers and runs the tests on all packages in this repository
- `start-stack` will start the stack required for testing

#### Structure
##### form3_sdk
Repository contains root package `form3_sdk` which provides a layer of abstracts
fetching data from the API - such as buffering the responses, or flexible way of filtering (see [buildFilter](util.go))

Root package allows for easy configurability:
- overwriting of the default HTTP client, for stubbing the requests/responses
- easy switching between multiple API hosts - when running more than one API version
- easy switching of API keys - for example for testing permission boundaries

##### Accounts
Accounts package abstracts interactions with Accounts API and provides range of functions to
create, query and delete resources.

This package utilises root package to make interactions with the API, thus is unaware
of hosts/api key and other configurations defined by the SDK client itself.

