# LittleURL API

This is the API that powers LittleURL backend functionality. While you are free to host it yourself, please keep in mind a few caveats;

1. This API is very heavily tied into AWS, to the point at which it can't realistically run anywhere else.

2. Running the API yourself is technically unsupported, while everything is open-source and you are free to deploy it yourself,
   I will be unlikely to offer technical support in doing so (except from fixing bugs, of course).

## Tracing

All of the lambda functions have support for [Lumigo Tracing](https://lumigo.io/) built in, it will automatically
be enabled/disabled based on the presence of the terraform variable `lumigo_token`.
