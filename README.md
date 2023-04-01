## GitHub installation tokens, made easy

Fetcher encapsulates the logic to retrieve a GitHub App installation's token.

It Generates the request jwt using:

- iat (issued at): One hour ago by the time the claim is generated. This [is the recomended setting by GitHub.](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-a-json-web-token-jwt-for-a-github-app#about-json-web-tokens-jwts).
- ext (expires at): 2 hours ahead from the claim generation. This is a hardset value that will be able to be passed on 0.2.
- iss (Issuer): the Application ID. GitHub ultimately uses this information to find the corresponding public key to validate the jwt signature.

The private key is *only* used to sign the jwt payload and it is *not* sent to GitHub's api. Jwt.io has a good [introductory page](https://jwt.io/introduction) about signing.

![Build and release](https://img.shields.io/github/v/release/luizfnunesmarques/github-token-fetcher.svg)
![Tests](https://github.com/luizfnunesmarques/token-fetcher/actions/workflows/lint-and-test.yml/badge.svg)
---
## use cases
- Local development helping to discover the capabilities of the GitHub API.
- Part of an automation workflow when the token to authenticate "as" a GitHub installation is required.

---
## usage & running

### Typical local example:

`./github-token-fetcher -applicationID="1" -installationID="1" -privateKeyFilePath="/path/to/file"`

or

`docker run -t luizfnunesmarques/token-fetcher . -applicationID "<ID>" -installationID "<ID>" -privateKeyFilePath "<path-to-file>"`

ps: Note that for the docker run the private key must be available to the running container

### Part of a worfklow in a container-based environment:
- Image can be tied to a step and have the environments accesible to its running container.

### The program accepts three arguments:
- applicationID: The actual application ID as shown to the owner of the app.
- installationID: The ID of the installation on which the token will be on behalf of.
- privateKeyFilePath: Path to find the private key file. Double check if the file and its DIR has appropriate permissions for read.



## :warning: Security Advisory
This program neither stores nor logs the content of the variable (specially the private key content) but it doesn't protect against a rogue agent with enough elevated access to read memory pages and therefore the contents of the program process' stack. Typically the agent would need `root`: in most cases where the host is safe, this program is safe to run.
