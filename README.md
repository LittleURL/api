# LittleURL API

This is the API that powers LittleURL backend functionality. While you are free to host it yourself, please keep in mind a few caveats;

1. This API is very heavily tied into AWS, to the point at which it can't realistically run anywhere else.

2. Running the API yourself is technically unsupported, while everything is open-source and you are free to deploy it yourself,
   I will be unlikely to offer technical support in doing so (except from fixing bugs, of course).

## SSM Parameters

There are some values that are required to be set in SSM, these are generally deployed via the [infrastructure stack](https://github.com/deltabyte/littleurl-infrstructure).
The reason for this is either due to needing abnormally high permissions, or being a globally unique resource that it used in multiple codebases.

| name                             | description                                               |
| -------------------------------- | --------------------------------------------------------- |
| `/littleurl/cloudflare-zone`     | ZoneID of the Cloudflare zone for deploying DNS records   |
| `/littleurl/api-certificate-arn` | ARN of the ACM certificate to be used for the API Gateway |

## Tracing

All of the lambda functions have support for [Lumigo Tracing](https://lumigo.io/) built in, it will automatically
be enabled/disabled based on the presence of the terraform variable `lumigo_token`.

## Email templates

The email templates are written using [MJML](https://mjml.io), as such they require compilation using a Node.js CLI.
Because the generated HTML is required for embedding into some of the golang binaries, the MJML compilation is automatically
handled whenever you run the deploy command in the Makefile.

For testing purposes, and so you don't have to re-generate the templates constantly, you will need to run each build step
manually.

```shell
# normally build everything
make build

# build templates once
make build-templates
# build functions as often as needed
make build-functions
```
