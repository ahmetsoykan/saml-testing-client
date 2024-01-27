# saml-testing-web-client

Simple web application for testing saml setup.

Each service provider(SP) must have an self-signed X.509 key pair established. You can generate your own with something like this:

    openssl req -x509 -newkey rsa:2048 -keyout myservice.key -out myservice.cert -days 365 -nodes -subj "/CN=myservice.example.com"

## The flow should look like this:

1. You browse to `localhost:8000/hello`
2. The middleware redirects you to IDP Login.
3. Prompts you for a username and password and returns you an HTML document which contains an HTML form setup to POST to `localhost:8000/saml/acs`. The form is automatically submitted if you have javascript enabled.
4. The local service validates the response, issues a session cookie, and redirects you to the original URL, `localhost:8000/hello`.
5. This time when `localhost:8000/hello` is requested there is a valid session and so the main content is served.

## IDP Setup:

When i'm testing saml, i've put these configuration to saml applications:

- Entity ID: test (you name it)
- ACS URLs: http://localhost:8000/saml/acs (This is the URL that IDP sends the user information after successful login)
- Attribute mappings:
  - email: Email Address

## Env Vars and Running

### Env Vars

- APP_IDPMETADATAURL (specific for saml application, this app gets xml metadata from this url)
- APP_ENTITYID (specific for saml application)

### Running

`go run main.go`

## Result

You should be seeing the 'Welcome, $your_email!' messages on your browser.
