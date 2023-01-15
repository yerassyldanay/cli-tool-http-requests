# CLI tool for making concurrent HTTP requests

This is a command-line tool for making concurrent HTTP requests to a list of URLs provided by the user. The tool returns the responses in decreasing order of the length of the response body.

## Features

- Concurrent requests: The tool can make multiple requests at the same time, which improves the performance of the tool.
- URL validation: The tool validates the URLs before making the requests to ensure that only valid URLs are processed.
- Sorted responses: The tool returns the responses in decreasing order of the length of the response body.
- Logging: The tool uses the zap package for logging and provides detailed logs for debugging.

## Getting started

1. Clone the repo by running `git clone github.com/yerassyldanay/cli-tool-http-requests`
2. Build the tool by running `make <OS type: windows,linux,darwin>`
3. Run the tool by providing a list of URLs through the `--urls` flag and any additional flags you want to use. Example: 

```code
./bin/linux/ydcli --urls=https://example.com,https://example2.com --verbose --error
```

You can use following example to make requests (45 valid, 5 invalid urls):
```code
--urls=https://google.com,https://youtube.com,https://facebook.com,https://twitter.com,https://instagram.com,https://linkedin.com,https://pinterest.com,https://reddit.com,https://apple.com,https://amazon.com,https://ebay.com,https://netflix.com,https://spotify.com,https://gmail.com,https://skype.com,https://whatsapp.com,https://tiktok.com,https://zoom.us,https://github.com,https://dropbox.com,https://soundcloud.com,https://wordpress.com,https://slack.com,https://telegram.com,https://trello.com,https://buffer.com,https://hootsuite.com,https://hubspot.com,https://surveymonkey.com,https://wix.com,https://shopify.com,https://squarespace.com,https://canva.com,https://issuu.com,https://weebly.com,https://wufoo.com,https://typeform.com,https://jotform.com,https://zendesk.com,https://salesforce.com,https://zoho.com,https://freshbooks.com,https://invalid.com,https://fake.com,https://notreal.com,https://nonexistent.com,https://doesnotexist.com
```

## Downloads

You can download the pre-built binary files in
[releases](https://github.com/yerassyldanay/cli-tool-http-requests/releases)

## Flags 

```code
  -error
        print error for each url at the end (if there is any)
  -urls string
        a list of URLs separated by a comma
  -verbose
        print more logs to monitor the process
```

## Dependencies

- Golang version 1.18 or higher

## Limitations

- The tool is not suitable for handling a large number of URLs (billions). In such cases, external sorting should be used to sort the responses.
- There should be a limiter for goroutines

Note: I do not believe that the user can bring a problem by providing a large number of urls
