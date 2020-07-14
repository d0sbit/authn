# authn
authn is an authentication library written in Go

## Getting Started

### Login

## Things to Know

- Each feature stands alone and can be wired separately.  Quick and simple to start, easy to add what you need as you build out your app.  Examples for each thing so you can just copy and paste and wire in the specifics for your app.
- Every handler does exactly one thing and so does not look at the path, only the HTTP method (or the `method` URL parameter can be used to override it)
- (explain inputs are JSON or URL or POST params)
- We don't serve pages, this is just an API.  Your application already will have a UI so trying to include that here would only get in your page.  TODO: add UI examples, probably using Vugu.
