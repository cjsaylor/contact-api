# Contact API

Simple web service that takes a post request to a contact endpoint and sends an email to a configured target address with contact information. This service is particularly useful for sites deployed as static HTML.

## Usage

The following endpoint is exposed:

```
POST /contact?sender=&subject=&body=
```

Headers:
* `Content-Type`: `application/x-www-urlencoded-form`


## Example implementation

Here is an example html form (minimal) and JS script that utilizes the endpoint:

```html
<form id="contactform">
	<input type="email" name="sender" required>
	<input type="text" name="subject" required>
	<textarea name="body" required></textarea>
	<!-- optional recaptcha spam filter -->
	<div class="g-recaptcha" data-sitekey="YOUR-KEY-HERE"></div>
	<button type="submit">Submit</button>
</form>
```

---

```js
const form = document.querySelector('#contactform');
form.addEventListener('submit', function(e) {
	e.preventDefault();
	const formData = new FormData(form);
	fetch('https://yourdomain.com/contact', {
		method: 'POST',
		body: new URLSearchParams(formData),
	})
	.then(res => {
		if (!res.ok) {
			throw new Error(res.blob());
		}
		console.log('Message sent!');
	})
	.catch(err => {
		console.error('Message failed to send :(');
	});
});
```

## Supported services

Email Services:

* [Mailgun](https://www.mailgun.com)

Spam Prevention Services:

* [Recaptcha](https://www.google.com/recaptcha/)

## Installation

### Via Docker

```bash
$ docker pull cjsaylor/contact-api
$ docker run --port 80:8080 --env <config keys here> cjsaylor/contact-api
```

### Via Golang build

```bash
$ go build -o web ./cmd/web/main.go
$ ./web
```

## Configuration

All configuration is done by environment variables.

| Env variable | Description
| --- | ---
| `PORT` | Port to expose the service (default: `8080`)
| `CE_CORS_DOMAIN` | Domain to allow ajax requests through from the browser. See the [MDN documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) for more details. (default `http://localhost:1313`)
| `CE_TARGET_EMAIL` | Email address of whom will recieve the messages.
| `RECAPTCHA_SECRET_KEY` | Recaptcha private API key, this is for spam protection. (optional)
| `MG_DOMAIN` | Mailgun registerd domain
| `MG_API_KEY` | Mailgun private api key
| `CE_TEST_MODE` | By default test mode is engaged to prevent accidental email sending, you must specify `CE_TEST_MODE=false` in production. While engaged, it will only log the message posted.

## Happy Users of this service

* [Blue Trails Tech](https://www.bluetrailstech.com)

## License

[GNU GPL v3.0](license.md)